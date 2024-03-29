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
	UpdatedAt  time.Time `json:"updated_at_iso"`
	CreatedAt  time.Time `json:"created_at_iso"`
}

func (a *ApplicationVariant) TypeString() string {
	if a.SubType == "" {
		return a.Type
	}
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
func (c *Client) CommitRolloutVariant(stackUid string, rolloutVariant ApplicationVariant) error {
	requestBody := struct {
		Operation string `json:"operation"`
	}{
		Operation: "commit",
	}
	req, err := c.NewRequest("PATCH", "/stacks/"+stackUid+"/application_variants/"+rolloutVariant.UID+".json", requestBody, nil)
	if err != nil {
		return err
	}
	return c.DoReq(req, nil, nil)
}

// UpdateCanaryRolloutPercentage updates the canary variant percentage
func (c *Client) UpdateCanaryRolloutPercentage(stackUid string, canaryVariant ApplicationVariant, canaryPercentage int) error {
	requestBody := struct {
		Operation        string `json:"operation"`
		CanaryPercentage int    `json:"canary_percentage"`
	}{
		Operation:        "commit",
		CanaryPercentage: canaryPercentage,
	}
	req, err := c.NewRequest("PATCH", "/stacks/"+stackUid+"/application_variants/"+canaryVariant.UID+".json", requestBody, nil)
	if err != nil {
		return err
	}
	return c.DoReq(req, nil, nil)
}

// DeletePreviewVariant removes the preview variant
func (c *Client) DeletePreviewVariant(stackUid string, previewVariant ApplicationVariant) error {
	req, err := c.NewRequest("DELETE", "/stacks/"+stackUid+"/application_variants/"+previewVariant.UID+".json", nil, nil)
	if err != nil {
		return err
	}
	return c.DoReq(req, nil, nil)
}
