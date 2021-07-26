package cloud66

import (
	"time"
)

// ApplicationVariant indicates a rails/rack application deployment variant
type ApplicationVariant struct {
	Id         int       `json:"id"`
	Type       string    `json:"type"`
	SubType    string    `json:"sub_type"`
	GitRepo    string    `json:"git_repo"`
	GitRef     string    `json:"git_ref"`
	GitHash    string    `json:"git_hash"`
	Tag        string    `json:"tag"`
	Percentage int       `json:"percentage"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// GetApplicationVariants starts a session via API
func (c *Client) GetApplicationVariants(stackUid string) ([]ApplicationVariant, error) {
	var variants []ApplicationVariant

	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/application_variants.json", nil, nil)
	if err != nil {
		return nil, err
	}

	err = c.DoReq(req, &variants, nil)
	if err != nil {
		return nil, err
	}
	return variants, nil
}