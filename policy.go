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
