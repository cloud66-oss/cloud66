package cloud66

import (
	"encoding/json"
)

type Pipeline struct {
	Version string          `json:"version"`
	Steps   json.RawMessage `json:"steps"`
}

func (c *Client) GetPipeline(stackUid, formationUid, snapshotUID string) (*Pipeline, error) {
	params := struct {
		SnapshotUID string `json:"snapshot_uid"`
	}{
		SnapshotUID: snapshotUID,
	}
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/formations/"+formationUid+"/pipeline.json", params, nil)
	if err != nil {
		return nil, err
	}
	var pipelineRes *Pipeline
	return pipelineRes, c.DoReq(req, &pipelineRes, nil)
}
