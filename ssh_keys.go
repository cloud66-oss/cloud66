package cloud66

import (
	"fmt"
)

type AccessPublicKey struct {
	Key string `json:"public_key"`
}

func (c *Client) GetAccessPublicKey() (*AccessPublicKey, error) {
	fmt.Printf("")

	req, err := c.NewRequest("GET", "/ssh_keys/access.json", nil, nil)
	if err != nil {
		return nil, err
	}

	var key *AccessPublicKey
	return key, c.DoReq(req, &key, nil)
}
