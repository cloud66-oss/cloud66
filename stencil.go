package cloud66

import "time"

type Stencil struct {
	Uid              string    `json:"uid"`
	Filename         string    `json:"filename"`
	TemplateFilename string    `json:"template_filename"`
	ContextID        string    `json:"context_id"`
	Status           int       `json:"status"`
	Tags             []string  `json:"tags"`
	Inline           bool      `json:"inline"`
	GitfilePath      string    `json:"gitfile_path"`
	Body             string    `json:"body"`
	BtrRepo          string    `json:"btr_repo"`
	BtrBranch        string    `json:"branch"`
	Sequence         int       `json:"sequence"`
	CreatedAt        time.Time `json:"created_at_iso"`
	UpdatedAt        time.Time `json:"updated_at_iso"`
}

func (s Stencil) String() string {
	return s.Filename
}

func (c *Client) AddStencils(stackUid string, formationUid string, baseTemplateUid string, stencils []*Stencil, message string) ([]Stencil, error) {
	params := struct {
		Message      string     `json:"message"`
		BaseTemplate string     `json:"btr_uuid"`
		Stencils     []*Stencil `json:"stencils"`
	}{
		Message:      message,
		BaseTemplate: baseTemplateUid,
		Stencils:     stencils,
	}

	var stencilRes []Stencil

	if len(stencils) > 0 {
		req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/formations/"+formationUid+"/stencils.json", params, nil)
		if err != nil {
			return nil, err
		}

		stencilRes = nil
		err = c.DoReq(req, &stencilRes, nil)
		if err != nil {
			return nil, err
		}
	}
	return stencilRes, nil
}
