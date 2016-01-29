package cloud66

import "strconv"

type Process struct {
	Name               string         `json:"name"`
	Md5                string         `json:"md5"`
	Command            string         `json:"command"`
	ServerProcessCount map[string]int `json:"servers"`
}

func (c *Client) GetProcesses(stackUid string, serverUid *string) ([]Process, error) {
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

	query_strings := make(map[string]string)
	query_strings["page"] = "1"

	var p Pagination
	var result []Process
	var processRes []Process

	for {
		req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/processes.json", params, query_strings)
		if err != nil {
			return nil, err
		}

		processRes = nil
		err = c.DoReq(req, &processRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, processRes...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}
	}
	return result, nil
}

func (c *Client) GetProcess(stackUid string, name string, serverUid *string) (*Process, error) {
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
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/processes/"+name+".json", params, nil)
	if err != nil {
		return nil, err
	}
	var processRes *Process
	return processRes, c.DoReq(req, &processRes, nil)
}

func (c *Client) ScaleProcess(stackUid string, processName string, serverCount map[string]int) (*AsyncResult, error) {
	params := struct {
		ProcessName string         `json:"process_name"`
		ServerCount map[string]int `json:"server_count"`
	}{
		ProcessName: processName,
		ServerCount: serverCount,
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/processes.json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}
