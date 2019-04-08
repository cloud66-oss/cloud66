package cloud66

import (
	"time"
)

type Policy struct {
	Uid       string    `json:"uid"`
	Name      string    `json:"name"`
	Selector  string    `json:"selector"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at_iso"`
	UpdatedAt time.Time `json:"updated_at_iso"`
	Tags      []string  `json:"tags"`
}

func (p Policy) String() string {
	return p.Name
}

func (c *Client) AddPolicies(stackUid string, formationUid string, policies []*Policy, message string) ([]Policy, error) {
	params := struct {
		Message  string    `json:"message"`
		Policies []*Policy `json:"policies"`
	}{
		Message:  message,
		Policies: policies,
	}

	var policiesRes []Policy

	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/formations/"+formationUid+"/policies.json", params, nil)
	if err != nil {
		return nil, err
	}

	policiesRes = nil
	err = c.DoReq(req, &policiesRes, nil)
	if err != nil {
		return nil, err
	}

	return policiesRes, nil
}
