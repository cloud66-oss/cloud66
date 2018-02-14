package cloud66

import (
	"strconv"
	"time"
)

type Formation struct {
	Uid          string    `json:"uid"`
	Name         string    `json:"name"`
	TemplateBase string    `json:"template_base"`
	CreatedAt    time.Time `json:"created_at_iso"`
	UpdatedAt    time.Time `json:"updated_at_iso"`
	Tags         []string  `json:"tags"`
}

func (c *Client) Formations(stackUid string) ([]Formation, error) {
	query_strings := make(map[string]string)
	query_strings["page"] = "1"

	var p Pagination
	var result []Formation
	var formationRes []Formation

	for {
		req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/formations.json", nil, query_strings)
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
