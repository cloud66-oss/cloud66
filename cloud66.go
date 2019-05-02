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
	"strconv"
	"strings"

	"github.com/khash/oauth/oauth"
	"github.com/pborman/uuid"
	"github.com/toqueteos/webbrowser"
)

const (
	baseURL = "https://app.cloud66.com"
)

type GenericResponse struct {
	Status  bool   `json:"ok"`
	Message string `json:"message"`
}

type Client struct {
	HTTP              *http.Client
	URL               string
	UserAgent         string
	AccountId         *int
	Debug             bool
	AdditionalHeaders http.Header

	defaultUserAgent string
	baseAPIURL       string
	defaultAPIURL    string
	authURL          string
	tokenURL         string
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
		apiURL = c.defaultAPIURL
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
	if c.AccountId != nil {
		req.Header.Set("X-Account", strconv.Itoa(*c.AccountId))
	}
	useragent := c.UserAgent
	if useragent == "" {
		useragent = c.defaultUserAgent
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

func (c *Client) Authorize(tokenDir, tokenFile, clientID, clientSecret, redirectURL, scope string) {
	err := os.MkdirAll(tokenDir, 0777)
	if err != nil {
		fmt.Printf("Failed to create directory for the token at %s\n", tokenDir)
	}
	cachefile := filepath.Join(tokenDir, tokenFile)

	config := &oauth.Config{
		ClientId:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scope:        scope,
		AuthURL:      c.authURL,
		TokenURL:     c.tokenURL,
		TokenCache:   oauth.CacheFile(cachefile),
	}
	transport := &oauth.Transport{Config: config}
	_, err = config.TokenCache.Token()

	// do we already have access?
	if err != nil {

		url := config.AuthCodeURL("")
		fmt.Printf("Openning %s\n", url)
		e := webbrowser.Open(url)
		if e != nil {
			fmt.Printf("Counldn't open the browser because %s\n", e.Error())
			fmt.Println("Please open the following URL in your browser and paste the access code here:")
			fmt.Println(url)
		} else {
			fmt.Println("Opening the browser so you can approve the client access")
		}

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

func GetClient(baseAPIURL string, tokenDir, tokenFile, version, agentPrefix, clientId, clientSecret, redirectURL, scope string) Client {
	c := Client{
		defaultAPIURL: baseAPIURL + "/api/3",
		authURL:       baseAPIURL + "/oauth/authorize",
		tokenURL:      baseAPIURL + "/oauth/token",
	}

	cachefile := filepath.Join(tokenDir, tokenFile)
	c.defaultUserAgent = agentPrefix + "/" + version + " (" + runtime.GOOS + "; " + runtime.GOARCH + ")"

	config := &oauth.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scope:        scope,
		AuthURL:      c.authURL,
		TokenURL:     c.tokenURL,
		TokenCache:   oauth.CacheFile(cachefile),
	}

	transport := &oauth.Transport{Config: config}
	token, _ := config.TokenCache.Token()
	transport.Token = token
	c.HTTP = transport.Client()

	return c
}
