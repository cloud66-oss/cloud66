package cloud66

import (
	"encoding/json"
)

type WorkflowWrapper struct {
	Pipeline   json.RawMessage `json:"pipeline"`
}

func (c *Client) GetWorkflow(stackUid, formationUid, snapshotUID string) (*WorkflowWrapper, error) {
	params := struct {
		SnapshotUID string `json:"snapshot_uid"`
	}{
		SnapshotUID: snapshotUID,
	}
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/formations/"+formationUid+"/pipeline.json", params, nil)
	if err != nil {
		return nil, err
	}
	var workflowWrapper *WorkflowWrapper
	return workflowWrapper, c.DoReq(req, &workflowWrapper, nil)
}
