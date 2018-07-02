package cloud66

import (
	"strconv"
	"strings"
	"time"
)

type Snapshot struct {
	Uid         string    `json:"uid"`
	CreatedAt   time.Time `json:"created_at_iso"`
	UpdatedAt   time.Time `json:"updated_at_iso"`
	Action      string    `json:"action"`
	TriggeredAt time.Time `json:"triggered_at"`
	TriggeredBy string    `json:"triggered_by"`
	Tags        []string  `json:"tags"`
}

type RenderError struct {
	Text    string `json:"text"`
	Link    string `json:"link"`
	Stencil string `json:"stencil"`
}

type Renders struct {
	Content        map[string]string `json:"content"`
	Errors         []RenderError     `json:"errors"`
	RequestedFiles []string          `json:"requested_files"`
	StencilGroup   string            `json:"stencil_group`
}

func (c *Client) Snapshots(stackUid string) ([]Snapshot, error) {
	query_strings := make(map[string]string)
	query_strings["page"] = "1"

	var p Pagination
	var result []Snapshot
	var snapshotRes []Snapshot

	for {
		req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/snapshots.json", nil, query_strings)
		if err != nil {
			return nil, err
		}

		snapshotRes = nil
		err = c.DoReq(req, &snapshotRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, snapshotRes...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	return result, nil
}

func (c *Client) RenderSnapshot(stackUid string, snapshotUid string, formationUid string, requestFiles []string, useLatest bool, stencilGroup string) (*Renders, error) {
	query_strings := make(map[string]string)
	query_strings["requested_files"] = strings.Join(requestFiles, ",")
	if !useLatest {
		// default is true on the server
		query_strings["use_latest"] = "false"
	}
	if stencilGroup != "" {
		query_strings["stencil_group"] = stencilGroup
	}

	var result *Renders

	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/snapshots/"+snapshotUid+"/formation/"+formationUid, nil, query_strings)
	if err != nil {
		return nil, err
	}

	result = nil
	err = c.DoReq(req, &result, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}
