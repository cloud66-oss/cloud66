package cloud66

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type FormationBundle struct {
	Version        string                 `json:"version"`
	Metadata       *Metadata              `json:"metadata"`
	Uid            string                 `json:"uid"`
	Name           string                 `json:"name"`
	StencilGroups  []*BundleStencilGroup  `json:"stencil_groups"`
	BaseTemplates  []*BundleBaseTemplates `json:"base_template"`
	Tags           []string               `json:"tags"`
	HelmReleases   []*BundleHelmRelease   `json:"helm_releases"`
	Configurations []string               `json:"configuration"`
}

type BundleHelmRelease struct {
	Uid           string `json:"uid"`
	ChartName     string `json:"chart_name"`
	DisplayName   string `json:"display_name"`
	Version       string `json:"version"`
	RepositoryURL string `json:"repository_url"`
	ValuesFile    string `json:"values_file"`
}

type BundleBaseTemplates struct {
	Name         string               `json:"name"`
	Repo         string               `json:"repo"`
	Branch       string               `json:"branch"`
	Stencils     []*BundleStencil     `json:"stencils"`
	Policies     []*BundlePolicy      `json:"policies"`
	Transformers []*BundleTransformer `json:"transformers"`
}

type Metadata struct {
	App         string    `json:"app"`
	Timestamp   time.Time `json:"timestamp"`
	Annotations []string  `json:"annotations"`
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
	Sequence int      `json:"sequence"`
	Tags     []string `json:"tags"`
}

type BundleTransformer struct { // this is just a placeholder for now
	Uid  string   `json:"uid"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func CreateFormationBundle(formation Formation, app string, configurations []string) *FormationBundle {
	bundle := &FormationBundle{
		Version: "1",
		Metadata: &Metadata{
			App:         app,
			Timestamp:   time.Now().UTC(),
			Annotations: make([]string, 0), //just a placeholder before creating the real method
		},
		Uid:            formation.Uid,
		Name:           formation.Name,
		Tags:           formation.Tags,
		BaseTemplates:  createBaseTemplates(formation),
		StencilGroups:  createStencilGroups(formation.StencilGroups),
		Configurations: configurations,
		HelmReleases:   createHelmReleases(formation.HelmReleses), //just a placeholder before creating the real method
	}
	return bundle
}

func createBaseTemplates(formation Formation) []*BundleBaseTemplates {
	baseTemplate := &BundleBaseTemplates{
		Repo:         formation.BaseTemplate.GitRepo,
		Branch:       formation.BaseTemplate.GitBranch,
		Stencils:     createStencils(formation.Stencils),
		Policies:     createPolicies(formation.Policies),
		Transformers: make([]*BundleTransformer, 0),
	}
	return append(make([]*BundleBaseTemplates, 0), baseTemplate)
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
			Sequence: st.Sequence,
			Tags:     st.Tags,
		}
	}

	return result
}

func (b *BundleStencil) AsStencil(bundlePath string) (*Stencil, error) {
	filePath := filepath.Join(filepath.Join(bundlePath, "stencils"), b.Filename)
	body, err := ioutil.ReadFile(filePath)

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

func (b *BundlePolicy) AsPolicy(bundlePath string) (*Policy, error) {
	filePath := filepath.Join(filepath.Join(bundlePath, "policies"), b.Uid+".cop")
	body, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	return &Policy{
		Uid:      b.Uid,
		Name:     b.Name,
		Selector: b.Selector,
		Sequence: b.Sequence,
		Body:     string(body),
		Tags:     b.Tags,
	}, nil
}

func createHelmReleases(helmReleases []HelmRelease) []*BundleHelmRelease {
	result := make([]*BundleHelmRelease, len(helmReleases))
	for idx, hr := range helmReleases {
		filename := hr.ChartName + "-values.yml"
		result[idx] = &BundleHelmRelease{
			ChartName:     hr.ChartName,
			DisplayName:   hr.DisplayName,
			Version:       hr.Version,
			RepositoryURL: hr.RepositoryURL,
			ValuesFile:    filename,
		}
	}

	return result
}

func (b *BundleHelmRelease) AsRelease(bundlePath string) (*HelmRelease, error) {
	filePath := filepath.Join(filepath.Join(bundlePath, "helm_releases"), b.ValuesFile)
	_, err := os.Stat(filePath)
	var body []byte
	if err != nil {
		body = nil
	} else {
		body, err = ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
	}

	return &HelmRelease{
		Uid:           b.Uid,
		ChartName:     b.ChartName,
		DisplayName:   b.DisplayName,
		RepositoryURL: b.RepositoryURL,
		Version:       b.Version,
		Body:          string(body),
	}, nil
}

/*
func (b *BundleStencilGroup) AsStencilGroup(bundlePath string) (*StencilGroup, error) {
	ext := filepath.Ext(b.Name)
	body, err := ioutil.ReadFile(filepath.Join(bundlePath, "stencil_groups", b.Uid) + ext)
	if err != nil {
		return nil, err
	}

	return &StencilGroup{
		Uid:      b.Uid,
		Name:     b.Name,
		Tags:     b.Tags,
		Body:     string(body),
		Sequence: b.Sequence,
	}, nil
}
*/
