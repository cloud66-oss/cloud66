package cloud66

import "time"

type Container struct {
	Uid                 string    `json:"uid"`
	ServerUid           string    `json:"server_uid"`
	ServiceName         string    `json:"service_name"`
	Image               string    `json:"image"`
	PortList            string    `json:"port_list"`
	Command             string    `json:"command"`
	StartedAt           time.Time `json:"started_at"`
	WebPorts            string    `json:"web_ports"`
	CaptureStdoutStderr bool      `json:"capture_stdout_stderr"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (c *Client) GetContainers(stackUid string) ([]Container, error) {
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/containers.json", nil)
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
