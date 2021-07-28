package cloud66

import (
	"fmt"
	"time"
)

// ApplicationVariant indicates a rails/rack application deployment variant
type ApplicationVariant struct {
	UID        string    `json:"uid"`
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

func (a *ApplicationVariant) TypeString() string {
	return fmt.Sprintf("%s/%s", a.Type, a.SubType)
}

func (a *ApplicationVariant) PercentageString() string {
	return fmt.Sprintf("%v%%", a.Percentage)
}

// GetApplicationVariants returns list of application variants
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

// CommitRolloutVariant locks in the selected rollout variant
func (c *Client) CommitRolloutVariant(stackUid string, rolloutVariant ApplicationVariant) (*AsyncResult, error) {
	requestBody := struct {
		Action string `json:"action"`
	}{
		Action: "commit",
	}
	req, err := c.NewRequest("PATCH", "/stacks/"+stackUid+"/application_variants/"+rolloutVariant.UID+".json", requestBody, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}

// UpdateCanaryRolloutPercentage updates the canary variant percentage
func (c *Client) UpdateCanaryRolloutPercentage(stackUid string, canaryVariant ApplicationVariant, canaryPercentage int) (*AsyncResult, error) {
	requestBody := struct {
		Action           string `json:"action"`
		CanaryPercentage int    `json:"canary_percentage"`
	}{
		Action:           "percentage",
		CanaryPercentage: canaryPercentage,
	}
	req, err := c.NewRequest("PATCH", "/stacks/"+stackUid+"/application_variants/"+canaryVariant.UID+".json", requestBody, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}

// DeletePreviewVariant removes the preview variant
func (c *Client) DeletePreviewVariant(stackUid string, previewVariant ApplicationVariant) (*AsyncResult, error) {
	req, err := c.NewRequest("DELETE", "/stacks/"+stackUid+"/application_variants/"+previewVariant.UID+".json", nil, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}
