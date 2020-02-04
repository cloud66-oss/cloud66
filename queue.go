package cloud66

import "encoding/json"

type Queue struct {
	Name string `json:"name"`
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

func (c *Client) PopQueue(queueName string, blocking bool) (json.RawMessage, error) {
	queryStrings := make(map[string]string)
	if blocking {
		queryStrings["blocking"] = "1"
	}

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
