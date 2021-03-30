package cloud66

import (
	"strconv"
	"time"
)

type ServiceYaml struct {
	Uid       string    `json:"uid"`
	Body      string    `json:"body"`
	Comments  string    `json:"comments"`
	CreatedAt time.Time `json:"created_at_iso"`
	UpdatedAt time.Time `json:"updated_at_iso"`
}

func (c *Client) ServiceYamlList(stackUid string, include_body bool) ([]ServiceYaml, error) {
	queryStrings := make(map[string]string)
	queryStrings["page"] = "1"
	queryStrings["include_body"] = strconv.FormatBool(include_body)

	var p Pagination
	var result []ServiceYaml
	var serviceYamlRes []ServiceYaml

	for {
		req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/service_yaml.json", nil, queryStrings)
		if err != nil {
			return nil, err
		}

		serviceYamlRes = nil
		err = c.DoReq(req, &serviceYamlRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, serviceYamlRes...)
		if p.Current < p.Next {
			queryStrings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	return result, nil
}

func (c *Client) ServiceYamlInfo(stackUid string, id string) (*ServiceYaml, error) {
	var serviceYamlRes ServiceYaml

	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/service_yaml/"+id+".json", nil, nil)
	if err != nil {
		return nil, err
	}

	err = c.DoReq(req, &serviceYamlRes, nil)
	if err != nil {
		return nil, err
	}

	return &serviceYamlRes, nil
}

func (c *Client) CreateServiceYaml(stackUid, serviceYaml, comments string) (*ServiceYaml, error) {
	params := struct {
		ServiceYaml string `json:"service_yaml"`
		Comments    string `json:"comments"`
	}{
		ServiceYaml: serviceYaml,
		Comments:    comments,
	}

	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/service_yaml", params, nil)
	if err != nil {
		return nil, err
	}

	var serviceYamlRes *ServiceYaml
	return serviceYamlRes, c.DoReq(req, &serviceYamlRes, nil)
}
