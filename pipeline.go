package cloud66

import (
	"encoding/json"
	"strconv"
)

type WorkflowWrapper struct {
	Workflow json.RawMessage `json:"pipeline"`
}

func (c *Client) GetWorkflow(stackUid, formationUid, snapshotUID string, useLatest bool, workflowName string) (*WorkflowWrapper, error) {
	queryStrings := make(map[string]string)
	queryStrings["page"] = "1"

	queryStrings["snapshot_uid"] = snapshotUID
	queryStrings["use_latest"] = strconv.FormatBool(useLatest)
	queryStrings["workflow"] = workflowName
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/formations/"+formationUid+"/pipeline.json", nil, queryStrings)
	if err != nil {
		return nil, err
	}
	var workflowWrapper *WorkflowWrapper
	return workflowWrapper, c.DoReq(req, &workflowWrapper, nil)
}
