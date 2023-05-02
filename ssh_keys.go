package cloud66

import (
	"encoding/base64"
)

type AccessPublicKey struct {
	Key string `json:"public_key"`
}

// GetAccessPublicKey gets the access public key for the currently authenticated user
func (c *Client) GetAccessPublicKey() (*AccessPublicKey, error) {
	req, err := c.NewRequest("GET", "/ssh_keys/access.json", nil, nil)
	if err != nil {
		return nil, err
	}
	var key *AccessPublicKey
	return key, c.DoReq(req, &key, nil)
}

// SetAccessPublicKey sets the access public key for the currently authenticated user
func (c *Client) SetAccessPublicKey(publicKey string) error {
	params := struct {
		PublicKey string `json:"public_key"`
	}{
		PublicKey: base64.StdEncoding.EncodeToString([]byte(publicKey)),
	}
	req, err := c.NewRequest("POST", "/ssh_keys/access.json", params, nil)
	if err != nil {
		return err
	}
	return c.DoReq(req, nil, nil)
}
