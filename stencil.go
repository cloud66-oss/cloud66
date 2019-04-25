package cloud66

import (
	"time"
)

type Stencil struct {
	Uid              string    `json:"uid"`
	Filename         string    `json:"filename"`
	TemplateFilename string    `json:"template_filename"`
	ContextID        string    `json:"context_id"`
	Status           int       `json:"status"`
	CreatedAt        time.Time `json:"created_at_iso"`
	UpdatedAt        time.Time `json:"updated_at_iso"`
	Tags             []string  `json:"tags"`
	Inline           bool      `json:"inline"`
	GitfilePath      string    `json:"gitfile_path"`
	Body             string    `json:"body"`
	Template         string    `json:"template"`
	Sequence         int       `json:"sequence"`
}

func (s Stencil) String() string {
	return s.Filename
}

func (c *Client) AddStencils(stackUid string, formationUid string, stencils []*Stencil, message string) ([]Stencil, error) {
	params := struct {
		Message  string     `json:"message"`
		Stencils []*Stencil `json:"stencils"`
	}{
		Message:  message,
		Stencils: stencils,
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
