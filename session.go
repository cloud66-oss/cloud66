package cloud66

import (
	"fmt"
)

// Session indicates a container based service session
type Session struct {
	UID            string `json:"uid"`
	Title          string `json:"title"`
	ServiceName    string `json:"service_name"`
	Command        string `json:"command"`
	Namespace      string `json:"namespace"`
	DeploymentName string `json:"deployment_name"`
	PodName        string `json:"pod_name"`
	ContainerName  string `json:"container_name"`
}

// StartRemoteSession starts a session via API
func (c *Client) StartRemoteSession(stackUID string, serviceName string) (*AsyncResult, error) {
	params := struct {
		ServiceName string `json:"service_name"`
	}{
		ServiceName: serviceName,
	}
	url := fmt.Sprintf("/stacks/%s/sessions.json", stackUID)
	req, err := c.NewRequest("POST", url, params, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}

// FetchRemoteSession fetches a session via API
func (c *Client) FetchRemoteSession(stackUID string, sessionUID, serviceName *string) (*Session, error) {
	queryStrings := make(map[string]string)
	queryStrings["page"] = "1"
	if serviceName != nil {
		queryStrings["service_name"] = *serviceName
	}
	if sessionUID != nil {
		queryStrings["session_id"] = *sessionUID
	}

	url := fmt.Sprintf("/stacks/%s/sessions/fetch.json", stackUID)
	req, err := c.NewRequest("GET", url, nil, queryStrings)
	if err != nil {
		return nil, err
	}
	var session *Session
	return session, c.DoReq(req, &session, nil)
}
