package cloud66

import "time"

type Onprem struct {
	Uid       string      `json:"uid"`
	Name      string      `json:"name"`
	Config    interface{} `json:"config"`
	CreatedAt time.Time   `json:"created_at_iso"`
	UpdatedAt time.Time   `json:"updated_at_iso"`
}
