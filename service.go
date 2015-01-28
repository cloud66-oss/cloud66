package cloud66

type Service struct {
	Name        string      `json:"name"`
	Containers  []Container `json:"containers"`
	WrapCommand string      `json:"wrap_command"`
}

func (c *Client) GetServices(stackUid string, serverUid *string) ([]Service, error) {
	var params interface{}
	if serverUid == nil {
		params = nil
	} else {
		params = struct {
			ServerUid string `json:"server_uid"`
		}{
			ServerUid: *serverUid,
		}
	}
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/services.json", params)
	if err != nil {
		return nil, err
	}
	var serviceRes []Service
	return serviceRes, c.DoReq(req, &serviceRes)
}

func (c *Client) GetService(stackUid string, serviceName string, serverUid *string, wrapCommand *string) (*Service, error) {
	var params interface{}
	if serverUid == nil {
		params = nil
	} else {
		if wrapCommand == nil {
			params = struct {
				ServerUid string `json:"server_uid"`
			}{
				ServerUid: *serverUid,
			}
		} else {
			params = struct {
				ServerUid   string `json:"server_uid"`
				WrapCommand string `json:"wrap_command"`
			}{
				ServerUid:   *serverUid,
				WrapCommand: *wrapCommand,
			}
		}
	}
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/services/"+serviceName+".json", params)
	if err != nil {
		return nil, err
	}
	var servicesRes *Service
	return servicesRes, c.DoReq(req, &servicesRes)
}

func (c *Client) StopService(stackUid string, serviceName string, serverUid *string) (*AsyncResult, error) {
	var params interface{}
	if serverUid == nil {
		params = nil
	} else {
		params = struct {
			ServerUid string `json:"server_uid"`
		}{
			ServerUid: *serverUid,
		}
	}
	req, err := c.NewRequest("DELETE", "/stacks/"+stackUid+"/services/"+serviceName+".json", params)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes)
}

func (c *Client) ScaleService(stackUid string, serviceName string, serverCount map[string]int) (*AsyncResult, error) {
	params := struct {
		ServiceName string         `json:"service_name"`
		ServerCount map[string]int `json:"server_count"`
	}{
		ServiceName: serviceName,
		ServerCount: serverCount,
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/services.json", params)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes)
}

func (s *Service) ServerContainerCountMap() map[string]int {
	var serverMap = make(map[string]int)
	for _, container := range s.Containers {
		if _, present := serverMap[container.ServerName]; present {
			serverMap[container.ServerName] = serverMap[container.ServerName] + 1
		} else {
			serverMap[container.ServerName] = 1
		}
	}
	return serverMap
}

func (c *Client) InvokeStackServiceAction(stackUid string, serviceName string, serverUid *string, action string) (*AsyncResult, error) {
	var params interface{}
	if serverUid == nil {
		params = struct {
			Command     string `json:"command"`
			ServiceName string `json:"service_name"`
		}{
			Command:     action,
			ServiceName: serviceName,
		}
	} else {
		params = struct {
			Command     string `json:"command"`
			ServiceName string `json:"service_name"`
			ServerUid   string `json:"server_uid"`
		}{
			Command:     action,
			ServiceName: serviceName,
			ServerUid:   *serverUid,
		}
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/actions.json", params)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes)
}
