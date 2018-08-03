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
