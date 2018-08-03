package cloud66

import (
	"time"
)

type StencilGroup struct {
	Uid       string    `json:"uid"`
	Name      string    `json:"name"`
	Rules     string    `json:"rules"`
	CreatedAt time.Time `json:"created_at_iso"`
	UpdatedAt time.Time `json:"updated_at_iso"`
	Tags      []string  `json:"tags"`
}

func (s StencilGroup) String() string {
	return s.Name
}
