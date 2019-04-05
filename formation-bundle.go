package cloud66

import (
	"io/ioutil"
	"path/filepath"
	"time"
)

type FormationBundle struct {
	Version        string                   `json:"version"`
	Metadata       *Metadata                `json:"metadata"`
	Uid            string                   `json:"uid"`
	Name           string                   `json:"name"`
	StencilGroups  []*BundleStencilGroup    `json:"stencil_groups"`
	BaseTemplates  []*BundleBaseTemplates   `json:"base_template"`
	Policies       []*BundlePolicy          `json:"policies"`
	Tags           []string                 `json:"tags"`
	HelmReleases   []*BundleHelmReleases    `json:"helm_releases"`
	Configurations []string                 `json:"configuration"`
}


type BundleHelmReleases struct {
	Name             string `json:"repo"`
	Version          string `json:"version"`
	RepositoryURL    string `json:"repository_url"`
	Values           string `json:"values_file"`
}

type BundleConfiguration struct {
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
}

type BundleBaseTemplates struct {
	Repo     string `json:"repo"`
	Branch   string `json:"branch"`
	Stencils []*BundleStencil `json:"stencils"`
}

type Metadata struct {
	App         string     `json:"app"`
	Timestamp   time.Time  `json:"timestamp"`
	Annotations []string   `json:"annotations"`
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
			App:         app,
			Timestamp:   time.Now().UTC(),
			Annotations: make([]string, 0), //just a placeholder before creating the real method
		},
		Uid:  formation.Uid,
		Name: formation.Name,
		Tags: formation.Tags,
		BaseTemplates: createBaseTemplates(formation),
		StencilGroups: createStencilGroups(formation.StencilGroups),
		Policies:      createPolicies(formation.Policies),
		Configurations: make([]string, 0), //just a placeholder before creating the real method
		HelmReleases: make([]*BundleHelmReleases, 0), //just a placeholder before creating the real method
	}
	return bundle
}

func createBaseTemplates(formation Formation) []*BundleBaseTemplates {
	baseTemplate := &BundleBaseTemplates{
		Repo:     formation.BaseTemplate.GitRepo,
		Branch:   formation.BaseTemplate.GitBranch,
		Stencils: createStencils(formation.Stencils),
		}
	return append(make([]*BundleBaseTemplates,0), baseTemplate)
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
