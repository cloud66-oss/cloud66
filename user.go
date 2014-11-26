package cloud66

import (
	"time"
	"fmt"
)

type UnmanagedServer struct {
	Vendor string `json:"vendor"`
	Id     string `json:"id"`
}

type Account struct {
	Id               int               `json:"id"`
	Owner            string            `json:"owner"`
	StackCount       int               `json:"stack_count"`
	UsedClouds       []string          `json:"used_clouds"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	CurrentAccount   bool              `json:"current_account"`
	UnmanagedServers []UnmanagedServer `json:"unmanaged_servers"`
}

func (c *Client) AccountInfos() ([]Account, error) {
	req, err := c.NewRequest("GET", "/accounts.json", nil)
	if err != nil {
		return nil, err
	}
	var accountRes []Account
	return accountRes, c.DoReq(req, &accountRes)
}

func (c *Client) AccountInfo(accountId int, getUnmanaged bool) (*Account, error) {
	params := struct {
		Value  bool `json:"include_servers"`
	}{
		Value: getUnmanaged,
	}

	req, err := c.NewRequest("GET", fmt.Sprintf("/accounts/%d.json", accountId), params)
	if err != nil {
		return nil, err
	}

	var accountRes *Account
	return accountRes, c.DoReq(req, &accountRes)
}
