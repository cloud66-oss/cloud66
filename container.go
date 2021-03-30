package cloud66

import (
	"strconv"
	"time"
)

type Port struct {
	Container int `json:"container"`
	Http      int `json:"http"`
	Https     int `json:"https"`
}

const (
	CNT_HEALTH_UNVERIFIED = 0
	CNT_HEALTH_UP         = 1
	CNT_HEALTH_DOWN       = 2
	CNT_HEALTH_NA         = 3
)

type Container struct {
	Uid             string    `json:"uid"`
	Name            string    `json:"name"`
	ServerUid       string    `json:"server_uid"`
	ServerName      string    `json:"server_name"`
	ServiceName     string    `json:"service_name"`
	Image           string    `json:"image"`
	Command         string    `json:"command"`
	Ports           []Port    `json:"ports"`
	PrivateIP       string    `json:"private_ip"`
	DockerIP        string    `json:"docker_ip"`
	HealthState     int       `json:"health_state"`
	HealthMessage   string    `json:"health_message"`
	HealthSource    string    `json:"health_source"`
	CaptureOutput   bool      `json:"capture_output"`
	RestartOnDeploy bool      `json:"restart_on_deploy"`
	StartedAt       time.Time `json:"started_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (c *Client) GetContainers(stackUid string, serverUid *string, serviceName *string) ([]Container, error) {
	queryStrings := make(map[string]string)
	queryStrings["page"] = "1"

	if serverUid != nil {
		queryStrings["server_uid"] = *serverUid
	}
	if serviceName != nil {
		queryStrings["service_name"] = *serviceName
	}

	var p Pagination
	var result []Container
	var containerRes []Container
	for {
		req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/containers.json", nil, queryStrings)
		if err != nil {
			return nil, err
		}

		containerRes = nil
		err = c.DoReq(req, &containerRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, containerRes...)
		if p.Current < p.Next {
			queryStrings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	return result, nil

}

func (c *Client) GetContainer(stackUid string, containerUid string) (*Container, error) {
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/containers/"+containerUid+".json", nil, nil)
	if err != nil {
		return nil, err
	}
	var containerRes *Container
	return containerRes, c.DoReq(req, &containerRes, nil)
}

func (c *Client) StopContainer(stackUid string, containerUid string) (*AsyncResult, error) {
	req, err := c.NewRequest("DELETE", "/stacks/"+stackUid+"/containers/"+containerUid+".json", nil, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}

func (c *Client) InvokeStackContainerAction(stackUid string, containerUid string, action string) (*AsyncResult, error) {
	params := struct {
		Command      string `json:"command"`
		ContainerUid string `json:"container_id"`
	}{
		Command:      action,
		ContainerUid: containerUid,
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/actions.json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}
