package cloud66

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type capturedRequest struct {
	method string
	path   string
	body   map[string]interface{}
}

func newGatewayTestClient(t *testing.T, captured *capturedRequest) *Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured.method = r.Method
		captured.path = r.URL.Path
		raw, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("reading request body: %v", err)
		}
		if err := json.Unmarshal(raw, &captured.body); err != nil {
			t.Fatalf("request body is not json: %v (%q)", err, raw)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response":{"ok":true}}`))
	}))
	t.Cleanup(server.Close)
	return &Client{URL: server.URL, UserAgent: "cloud66-test"}
}

func TestUpdateGatewayAttributesSendsOnlyTheFieldsThatAreSet(t *testing.T) {
	var captured capturedRequest
	client := newGatewayTestClient(t, &captured)

	err := client.UpdateGatewayAttributes(139, 7, GatewayUpdate{Address: "203.0.113.10", PrivateIp: "10.0.1.5"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if captured.method != "PUT" || captured.path != "/accounts/139/gateways/7.json" {
		t.Fatalf("got %s %s, want PUT /accounts/139/gateways/7.json", captured.method, captured.path)
	}
	want := map[string]interface{}{"address": "203.0.113.10", "private_ip": "10.0.1.5"}
	if len(captured.body) != len(want) {
		t.Fatalf("got body %v, want exactly %v", captured.body, want)
	}
	for key, value := range want {
		if captured.body[key] != value {
			t.Errorf("body[%q] = %v, want %v", key, captured.body[key], value)
		}
	}
}

func TestUpdateGatewayAttributesNeverSendsTtlOrContent(t *testing.T) {
	var captured capturedRequest
	client := newGatewayTestClient(t, &captured)

	err := client.UpdateGatewayAttributes(139, 7, GatewayUpdate{Name: "bastion-2", Username: "ec2-user"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if captured.method != "PUT" {
		t.Fatalf("no PUT was sent, got method %q", captured.method)
	}
	// the server closes the gateway on any ttl below 2, so even a zero ttl here would lock the customer out
	for _, key := range []string{"ttl", "content"} {
		if _, present := captured.body[key]; present {
			t.Errorf("body has %q, which would change the gateway's key or open state: %v", key, captured.body)
		}
	}
}
