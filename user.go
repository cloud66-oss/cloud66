package cloud66

import (
	"time"
)

type Account struct {
	Owner      string    `json:"owner"`
	StackCount int       `json:"stack_count"`
	UsedClouds []string  `json:"used_clouds"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c *Client) AccountInfos() ([]Account, error) {
	req, err := c.NewRequest("GET", "/accounts.json", nil)
	if err != nil {
		return nil, err
	}
	var accountRes []Account
	return accountRes, c.DoReq(req, &accountRes)
}
