package cloud66

import (
	"fmt"
	"time"
)

type GatewayKey struct {
	Id        int       `json:"id"`
	Ttl       string    `json:"ttl"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Client) ListGatewayKey(accountId int) (*GatewayKey, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("/accounts/%d/gateway_key.json", accountId), nil, nil)
	if err != nil {
		return nil, err
	}

	var gatewayKey *GatewayKey
	return gatewayKey, c.DoReq(req, &gatewayKey, nil)
}

func (c *Client) AddGatewayKey(accountId int, keyContent string, ttl int) error {
	params := struct {
		Content string `json:"content"`
		Ttl     int    `json:"ttl"`
	}{
		Content: keyContent,
		Ttl:     ttl,
	}

	req, err := c.NewRequest("POST", fmt.Sprintf("/accounts/%d/gateway_key.json", accountId), params, nil)
	if err != nil {
		return err
	}

	err = c.DoReq(req, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveGatewayKey(accountId int) error {
	req, err := c.NewRequest("DELETE", fmt.Sprintf("/accounts/%d/gateway_key.json", accountId), nil, nil)
	if err != nil {
		return err
	}

	err = c.DoReq(req, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
