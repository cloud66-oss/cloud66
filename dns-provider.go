package cloud66

import (
	"fmt"
	"strconv"
	"time"
)

type DnsProvider struct {
	Uuid        string    `json:"uuid"`
	DisplayName string    `json:"display_name"`
	Type        string    `json:"type"`
	Key         string    `json:"key"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Client) ListDnsProviders(stackUID string) ([]DnsProvider, error) {
	queryStrings := make(map[string]string)
	queryStrings["page"] = "1"
	var p Pagination
	var result []DnsProvider
	var pageResult []DnsProvider
	for {
		req, err := c.NewRequest("GET", fmt.Sprintf("/stacks/%s/dns_providers.json", stackUID), nil, queryStrings)
		if err != nil {
			return nil, err
		}
		
		pageResult = nil
		err = c.DoReq(req, &pageResult, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, pageResult...)
		if p.Current < p.Next {
			queryStrings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}
	}
	return result, nil
}
