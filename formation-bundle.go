package cloud66

import (
	"io/ioutil"
	"path/filepath"
	"time"
)

type FormationBundle struct {
	Version       string                `json:"version"`
	Metadata      *Metadata             `json:"metadata"`
	Uid           string                `json:"uid"`
	Name          string                `json:"name"`
	Stencils      []*BundleStencil      `json:"stencils"`
	StencilGroups []*BundleStencilGroup `json:"stencil_groups"`
	BaseTemplate  *BundleBaseTemplate   `json:"base_template"`
	Policies      []*BundlePolicy       `json:"policies"`
	Tags          []string              `json:"tags"`
}

type BundleBaseTemplate struct {
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
}

type Metadata struct {
	App       string    `json:"app"`
	Timestamp time.Time `json:"timestamp"`
}

type BundleStencil struct {
	Uid              string   `json:"uid"`
	Filename         string   `json:"filename"`
	TemplateFilename string   `json:"template_filename"`
	ContextID        string   `json:"context_id"`
	Status           int      `json:"status"`
	Tags             []string `json:"tags"`
	Sequence         int      `json:"sequence"`
}

type BundleStencilGroup struct {
	Uid  string   `json:"uid"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type BundlePolicy struct {
	Uid      string   `json:"uid"`
	Name     string   `json:"name"`
	Selector string   `json:"selector"`
	Tags     []string `json:"tags"`
}

func CreateFormationBundle(formation Formation, app string) *FormationBundle {
	bundle := &FormationBundle{
		Version: "1",
		Metadata: &Metadata{
			App:       app,
			Timestamp: time.Now().UTC(),
		},
		Uid:  formation.Uid,
		Name: formation.Name,
		Tags: formation.Tags,
		BaseTemplate: &BundleBaseTemplate{
			Repo:   formation.BaseTemplate.GitRepo,
			Branch: formation.BaseTemplate.GitBranch,
		},
		Stencils:      createStencils(formation.Stencils),
		StencilGroups: createStencilGroups(formation.StencilGroups),
		Policies:      createPolicies(formation.Policies),
	}

	return bundle
}

func createStencils(stencils []Stencil) []*BundleStencil {
	result := make([]*BundleStencil, len(stencils))
	for idx, st := range stencils {
		result[idx] = &BundleStencil{
			Uid:              st.Uid,
			Filename:         st.Filename,
			ContextID:        st.ContextID,
			TemplateFilename: st.TemplateFilename,
			Status:           st.Status,
			Tags:             st.Tags,
			Sequence:         st.Sequence,
		}
	}

	return result
}

func createStencilGroups(stencilGroups []StencilGroup) []*BundleStencilGroup {
	result := make([]*BundleStencilGroup, len(stencilGroups))
	for idx, st := range stencilGroups {
		result[idx] = &BundleStencilGroup{
			Name: st.Name,
			Uid:  st.Uid,
			Tags: st.Tags,
		}
	}

	return result
}

func createPolicies(policies []Policy) []*BundlePolicy {
	result := make([]*BundlePolicy, len(policies))
	for idx, st := range policies {
		result[idx] = &BundlePolicy{
			Uid:      st.Uid,
			Name:     st.Name,
			Selector: st.Selector,
			Tags:     st.Tags,
		}
	}

	return result
}

func (b *BundleStencil) AsStencil(bundlePath string) (*Stencil, error) {
	ext := filepath.Ext(b.Filename)
	body, err := ioutil.ReadFile(filepath.Join(bundlePath, "stencils", b.Uid) + ext)
	if err != nil {
		return nil, err
	}

	return &Stencil{
		Uid:              b.Uid,
		Filename:         b.Filename,
		TemplateFilename: b.TemplateFilename,
		ContextID:        b.ContextID,
		Status:           b.Status,
		Tags:             b.Tags,
		Body:             string(body),
		Sequence:         b.Sequence,
	}, nil
}
