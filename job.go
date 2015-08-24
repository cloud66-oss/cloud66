package cloud66

import "strconv"

var JobStatus = map[int]string{
	0: "Updated", // ST_UPDATED
	1: "Started", // ST_STARTED
	2: "Success", // ST_SUCCESS
	3: "Failed",  // ST_FAILED
}

type Job struct {
	Id     int               `json:"id"`
	Uid    string            `json:"uid"`
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Cron   string            `json:"cron"`
	Status int               `json:"status"`
	Params map[string]string `json:"params"`
}

func (c *Client) GetJobs(stackUid string, serverUid *string) ([]Job, error) {
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
	var result []Job
	var jobRes []Job

	for {
		req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/jobs.json", params, query_strings)
		if err != nil {
			return nil, err
		}

		jobRes = nil
		err = c.DoReq(req, &jobRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, jobRes...)
		if p.Current < p.Next {
			query_strings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}

	}

	return result, nil
}

func (c *Client) GetJob(stackUid string, jobUid string) (*Job, error) {
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/jobs/"+jobUid+".json", nil, nil)
	if err != nil {
		return nil, err
	}
	var jobRes *Job
	return jobRes, c.DoReq(req, &jobRes, nil)
}
