package cloud66

import (
	"time"
)

type FormationFilter struct {
	Name                 string    `json:"name"`
	UID                  string    `json:"uid"`
	StencilFilter        string    `json:"stencil_filter"`
	HelmFilter           string    `json:"helm_filter"`
	TransformationFilter string    `json:"transformation_filter"`
	PolicyFilter         string    `json:"policy_filter"`
	WorkflowFilter       string    `json:"workflow_filter"`
	CreatedAt            time.Time `json:"created_at_iso"`
	UpdatedAt            time.Time `json:"updated_at_iso"`
}
