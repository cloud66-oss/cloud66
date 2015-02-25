package cloud66

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/goauth2/oauth"
)

const (
	baseURL                = "https://app.cloud66.com"
	productionClientId     = "d4631fd51633bef0c04c6f946428a61fb9089abf4c1e13c15e9742cafd84a91f"
	productionClientSecret = "e663473f7b991504eb561e208995de15550f499b6840299df588cebe981ba48e"
	scope                  = "public redeploy jobs users admin"
	redirectURL            = "urn:ietf:wg:oauth:2.0:oob"
)

var (
	defaultUserAgent string
	baseAPIURL       string
	defaultAPIURL    string
	clientId         = os.Getenv("CX_APP_ID")
	clientSecret     = os.Getenv("CX_APP_SECRET")
	authURL          string
	tokenURL         string
)

type GenericResponse struct {
	Status  bool   `json:"ok"`
	Message string `json:"message"`
}

type Client struct {
	HTTP              *http.Client
	URL               string
	UserAgent         string
	Debug             bool
	AdditionalHeaders http.Header
}

type Response struct {
	Response   json.RawMessage
	Count      int
	Pagination json.RawMessage
}

type Pagination struct {
	Previous int
	Next     int
	Current  int
}

type filterFunction func(item interface{}) bool

func init() {
	baseAPIURL = os.Getenv("CLOUD66_API_URL")
	if baseAPIURL == "" {
		baseAPIURL = baseURL
	}

	defaultAPIURL = baseAPIURL + "/api/3"
	authURL = baseAPIURL + "/oauth/authorize"
	tokenURL = baseAPIURL + "/oauth/token"
}

func (c *Client) Get(v interface{}, path string, query_strings map[string]string, p *Pagination) error {
	return c.APIReq(v, "GET", path, nil, query_strings, p)
}

func (c *Client) Patch(v interface{}, path string, body interface{}) error {
	return c.APIReq(v, "PATCH", path, body, nil, nil)
}

func (c *Client) Post(v interface{}, path string, body interface{}) error {
	return c.APIReq(v, "POST", path, body, nil, nil)
}

func (c *Client) Put(v interface{}, path string, body interface{}) error {
	return c.APIReq(v, "PUT", path, body, nil, nil)
}

func (c *Client) Delete(path string) error {
	return c.APIReq(nil, "DELETE", path, nil, nil, nil)
}

func (c *Client) NewRequest(method, path string, body interface{}, query_strings map[string]string) (*http.Request, error) {
	var ctype string
	var rbody io.Reader

	switch t := body.(type) {
	case nil:
	case string:
		rbody = bytes.NewBufferString(t)
	case io.Reader:
		rbody = t
	default:
		v := reflect.ValueOf(body)
		if !v.IsValid() {
			break
		}
		if v.Type().Kind() == reflect.Ptr {
			v = reflect.Indirect(v)
			if !v.IsValid() {
				break
			}
		}

		j, err := json.Marshal(body)
		if err != nil {
			log.Fatal(err)
		}
		rbody = bytes.NewReader(j)
		ctype = "application/json"
	}
	apiURL := strings.TrimRight(c.URL, "/")
	if apiURL == "" {
		apiURL = defaultAPIURL
	}

	var qs string
	if (query_strings != nil) && (len(query_strings) > 0) {
		for key, value := range query_strings {
			if qs == "" {
				qs = "?"
			} else {
				qs = qs + "&"
			}
			qs = qs + key + "=" + value
		}
	}

	last_url := strings.TrimRight(apiURL+path, "/")
	if qs != "" {
		last_url = last_url + qs
	}

	req, err := http.NewRequest(method, last_url, rbody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Request-Id", uuid.New())
	if os.Getenv("CXTOKEN") != "" {
		req.Header.Set("X-CxToken", os.Getenv("CXTOKEN"))
	}
	useragent := c.UserAgent
	if useragent == "" {
		useragent = defaultUserAgent
	}
	req.Header.Set("User-Agent", useragent)
	if ctype != "" {
		req.Header.Set("Content-Type", ctype)
	}
	for k, v := range c.AdditionalHeaders {
		req.Header[k] = v
	}
	return req, nil
}

func (c *Client) APIReq(v interface{}, meth, path string, body interface{}, query_strings map[string]string, p *Pagination) error {
	req, err := c.NewRequest(meth, path, body, query_strings)
	if err != nil {
		return err
	}
	return c.DoReq(req, v, p)
}

func (c *Client) DoReq(req *http.Request, v interface{}, p *Pagination) error {

	if c.Debug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			log.Println(err)
		} else {
			os.Stderr.Write(dump)
			os.Stderr.Write([]byte{'\n', '\n'})
		}
	}

	var check_pagination bool
	if (req.Method == "GET") && (p != nil) {
		check_pagination = true
	} else {
		check_pagination = false
	}

	httpClient := c.HTTP
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if c.Debug {
		dump, err := httputil.DumpResponse(res, true)
		if err != nil {
			log.Println(err)
		} else {
			os.Stderr.Write(dump)
			os.Stderr.Write([]byte{'\n'})
		}
	}
	if err = checkResp(res); err != nil {
		return err
	}

	// open the wrapper
	var r Response
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(r.Response)

	switch t := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(t, buffer)
	default:
		err = json.NewDecoder(buffer).Decode(v)
	}

	if (err == nil) && check_pagination {
		pagination := bytes.NewBuffer(r.Pagination)
		err = json.NewDecoder(pagination).Decode(p)
	}

	return err
}

type Error struct {
	error
	Id string
}

type errorResp struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
	Details     string `json:"details"`
}

func checkResp(res *http.Response) error {
	if res.StatusCode/100 != 2 { // 200, 201, 202, etc
		var e errorResp
		err := json.NewDecoder(res.Body).Decode(&e)
		if err != nil {
			return errors.New("Unexpected error: " + res.Status)
		}
		if e.Details != "" {
			return Error{error: errors.New(e.Details), Id: e.Error}
		} else {
			return Error{error: errors.New(e.Description), Id: e.Error}
		}

	}
	if msg := res.Header.Get("X-Cloud66-Warning"); msg != "" {
		fmt.Fprintln(os.Stderr, strings.TrimSpace(msg))
	}
	return nil
}

func Authorize(tokenDir, tokenFile string) {
	err := os.MkdirAll(tokenDir, 0777)
	if err != nil {
		fmt.Printf("Failed to create directory for the token at %s\n", tokenDir)
	}
	cachefile := filepath.Join(tokenDir, tokenFile)

	if clientId == "" {
		clientId = productionClientId
	}
	if clientSecret == "" {
		clientSecret = productionClientSecret
	}
	config := &oauth.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scope:        scope,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
		TokenCache:   oauth.CacheFile(cachefile),
	}
	transport := &oauth.Transport{Config: config}
	_, err = config.TokenCache.Token()

	// do we already have access?
	if err != nil {

		url := config.AuthCodeURL("")
		fmt.Println("Please open the following URL in your browser and paste the access code here:")
		fmt.Println(url)

		var s string
		fmt.Println("Authorization Code:")
		fmt.Scan(&s)

		_, err := transport.Exchange(s)
		if err != nil {
			log.Fatal("Exchange:", err)
		}

		fmt.Printf("Token is cached in %v\n", config.TokenCache)
		os.Exit(1)
	}
}

func GetClient(tokenDir, tokenFile, version string) Client {
	cachefile := filepath.Join(tokenDir, tokenFile)
	defaultUserAgent = "cx/" + version + " (" + runtime.GOOS + "; " + runtime.GOARCH + ")"

	config := &oauth.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scope:        scope,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
		TokenCache:   oauth.CacheFile(cachefile),
	}

	transport := &oauth.Transport{Config: config}
	token, _ := config.TokenCache.Token()
	transport.Token = token

	return Client{HTTP: transport.Client()}
}
