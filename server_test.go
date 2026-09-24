package cloud66

import (
	"encoding/json"
	"testing"
)

func TestServerDecodesBothIpv4Addresses(t *testing.T) {
	payload := `{"uid":"srv-1","address":"203.0.113.20","ext_ipv4":"203.0.113.20","int_ipv4":"10.0.1.20"}`

	var server Server
	if err := json.Unmarshal([]byte(payload), &server); err != nil {
		t.Fatalf("decoding server: %v", err)
	}

	if server.IntIpV4 != "10.0.1.20" {
		t.Errorf("IntIpV4 = %q, want %q", server.IntIpV4, "10.0.1.20")
	}
	if server.ExtIpV4 != "203.0.113.20" {
		t.Errorf("ExtIpV4 = %q, want %q", server.ExtIpV4, "203.0.113.20")
	}
}
