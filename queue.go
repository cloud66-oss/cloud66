package cloud66

import "encoding/json"

type Queue struct {
	Name string `json:"name"`
}

type AgentRegistration struct {
	LogKey    string `json:"log_key"`
	AgentUUID string `json:"agent_uuid"`
}

func (c *Client) RegisterAgent() (*AgentRegistration, error) {
	var payload = struct {
		Hostname  string `json:"host_name"`
	}{
		Hostname:  c.Hostname,
	}

	req, err := c.NewRequest("POST", "/queues/register.json", payload, nil)
	if err != nil {
		return nil, err
	}
	
	var queueRes AgentRegistration

	err = c.DoReq(req, &queueRes, nil)
	if err != nil {
		return nil, err
	}

	return &queueRes, nil 
}

func (c *Client) GetQueues() ([]Queue, error) {
	req, err := c.NewRequest("GET", "/queues.json", nil, nil)
	if err != nil {
		return nil, err
	}

	var queueRes []Queue
	queueRes = nil
	err = c.DoReq(req, &queueRes, nil)
	if err != nil {
		return nil, err
	}

	return queueRes, nil
}

func (c *Client) PopQueue(queueName string) (json.RawMessage, error) {
	queryStrings := make(map[string]string)
	queryStrings["host_name"] = c.Hostname

	req, err := c.NewRequest("GET", "/queues/"+queueName+".json", nil, queryStrings)
	if err != nil {
		return nil, err
	}

	var queueRes json.RawMessage
	err = c.DoReq(req, &queueRes, nil)
	if err != nil {
		return nil, err
	}

	return queueRes, nil
}

func (c *Client) UpdateQueue(queueName string, taskUUID string, state string, runResult string, sessionID string) (json.RawMessage, error) {
	var payload = struct {
		TaskUUID  string `json:"task_uuid"`
		State     string `json:"state"`
		RunResult string `json:"run_result"`
		Hostname  string `json:"host_name"`
		SessionID string `json:"session_id"`
	}{
		TaskUUID:  taskUUID,
		State:     state,
		RunResult: runResult,
		Hostname:  c.Hostname,
		SessionID: sessionID,
	}

	req, err := c.NewRequest("PUT", "/queues/"+queueName+".json", payload, nil)
	if err != nil {
		return nil, err
	}

	var queueRes json.RawMessage
	err = c.DoReq(req, &queueRes, nil)
	if err != nil {
		return nil, err
	}

	return queueRes, nil
}
