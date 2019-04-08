package cloud66

import "time"

type HelmRelease struct {
	Name          string    `json:"name"`
	Version       string    `json:"version"`
	RepositoryURL string    `json:"repository"`
	Values        string    `json:"values"`
	Body          string    `json:"body"`
	CreatedAt     time.Time `json:"created_at_iso"`
	UpdatedAt     time.Time `json:"updated_at_iso"`
}

func (p HelmRelease) String() string {
	return p.Name
}

func (c *Client) AddHelmReleases(stackUid string, formationUid string, releases []*HelmRelease, message string) ([]HelmRelease, error) {
	params := struct {
		Message      string         `json:"message"`
		HelmReleases []*HelmRelease `json:"helm_releases"`
	}{
		Message:      message,
		HelmReleases: releases,
	}

	var releasesRes []HelmRelease

	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/formations/"+formationUid+"/helm_releases.json", params, nil)
	if err != nil {
		return nil, err
	}

	releasesRes = nil
	err = c.DoReq(req, &releasesRes, nil)
	if err != nil {
		return nil, err
	}

	return releasesRes, nil
}
