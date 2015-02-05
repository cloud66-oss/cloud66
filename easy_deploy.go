package cloud66

import (
	"fmt"
)

type EasyDeployMaintainer struct {
	Name    *string `json:"name"`
	Email   string  `json:"email"`
	Company *string `json:"company"`
	Offical *bool   `json:"official"`
	Trusted *bool   `json:"trusted"`
}

type EasyDeploy struct {
	Name        string               `json:"name"`
	DisplayName *string              `json:"display_name"`
	Version     string               `json:"version"`
	Uid         string               `json:"uid"`
	CreatedAt   string               `json:"created_at"`
	Logo        *string              `json:"logo"`
	Maintainer  EasyDeployMaintainer `json:"maintainer"`
}

func (c *Client) EasyDeployList() ([]string, error) {
	req, err := c.NewRequest("GET", "/easy_deploys.json", nil)
	if err != nil {
		return nil, err
	}
	var easyDeploy []string
	return easyDeploy, c.DoReq(req, &easyDeploy)
}

func (c *Client) EasyDeployInfo(name string) (*EasyDeploy, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("/easy_deploys/%s.json", name), nil)
	if err != nil {
		return nil, err
	}

	var easyDeploy *EasyDeploy
	return easyDeploy, c.DoReq(req, &easyDeploy)
}
