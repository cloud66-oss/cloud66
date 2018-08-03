package cloud66

import (
	"fmt"
	"strconv"
	"time"
)

type Formation struct {
	Uid           string         `json:"uid"`
	Name          string         `json:"name"`
	Stencils      []Stencil      `json:"stencils"`
	StencilGroups []StencilGroup `json:"stencil_groups"`
	BaseTemplate  BaseTemplate   `json:"base_template"`
	Policies      []Policy       `json:"policies"`
	CreatedAt     time.Time      `json:"created_at_iso"`
	UpdatedAt     time.Time      `json:"updated_at_iso"`
	Tags          []string       `json:"tags"`
}

func (c *Client) Formations(stackUid string, fullContent bool) ([]Formation, error) {
	query_strings := make(map[string]string)
	query_strings["page"] = "1"
	if fullContent {
		query_strings["full_content"] = "1"
	}

	var p Pagination
	var result []Formation
	var formationRes []Formation

	for {
		req, err := c.NewRequest("GET", fmt.Sprintf("/stacks/%s/formations.json", stackUid), nil, query_strings)
		if err != nil {
			return nil, err
		}

		formationRes = nil
		err = c.DoReq(req, &formationRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, formationRes...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	return result, nil
}
