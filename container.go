package cloud66

import "time"

type Port struct {
	Container int `json:"container"`
	Http      int `json:"http"`
	Https     int `json:"https"`
}

type Container struct {
	Uid             string    `json:"uid"`
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
	type Params struct {
		ServerUid   string `json:"server_uid"`
		ServiceName string `json:"service_name"`
	}
	var params Params
	if serverUid != nil && serviceName != nil {
		params = Params{
			ServerUid:   *serverUid,
			ServiceName: *serviceName,
		}
	} else if serverUid != nil {
		params = Params{
			ServerUid: *serverUid,
		}
	} else if serviceName != nil {
		params = Params{
			ServiceName: *serviceName,
		}
	}
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/containers.json", params)
	if err != nil {
		return nil, err
	}
	var containerRes []Container
	return containerRes, c.DoReq(req, &containerRes)
}

func (c *Client) GetContainer(stackUid string, containerUid string) (*Container, error) {
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/containers/"+containerUid+".json", nil)
	if err != nil {
		return nil, err
	}
	var containerRes *Container
	return containerRes, c.DoReq(req, &containerRes)
}

func (c *Client) StopContainer(stackUid string, containerUid string) (*AsyncResult, error) {
	req, err := c.NewRequest("DELETE", "/stacks/"+stackUid+"/containers/"+containerUid+".json", nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes)
}

func (c *Client) InvokeStackContainerAction(stackUid string, containerUid string, action string) (*AsyncResult, error) {
	params := struct {
		Command      string `json:"command"`
		ContainerUid string `json:"container_id"`
	}{
		Command:      action,
		ContainerUid: containerUid,
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/actions.json", params)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes)
}
