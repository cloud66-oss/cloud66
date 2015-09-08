package cloud66

import "strconv"
import "encoding/json"
import "fmt"

var JobStatus = map[int]string{
	0: "Updated", // ST_UPDATED
	1: "Started", // ST_STARTED
	2: "Success", // ST_SUCCESS
	3: "Failed",  // ST_FAILED
}

type Job interface {
	GetBasicJob() BasicJob
}

type basicJob struct {
	Id        int             `json:"id"`
	Uid       string          `json:"uid"`
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Cron      string          `json:"cron"`
	Status    int             `json:"status"`
	ParamsRaw json.RawMessage `json:"params"`
	Params    map[string]string
}

type BasicJob struct {
	*basicJob
}

func (bj *BasicJob) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &bj.basicJob); err == nil {
		if err = json.Unmarshal(bj.ParamsRaw, &bj.basicJob.Params); err != nil {
			return err
		}
		return nil
	} else {
		return err
	}

}

func (job BasicJob) GetBasicJob() BasicJob {
	return job
}

// ----------------

type dockerHostTaskJob struct {
	Command string `json:"command"`
}

type DockerHostTaskJob struct {
	*BasicJob
	*dockerHostTaskJob
}

func (job *DockerHostTaskJob) UnmarshalJSON(b []byte) error {
	var bj BasicJob
	if err := json.Unmarshal(b, &bj); err == nil {
		var j dockerHostTaskJob
		if err := json.Unmarshal(bj.basicJob.ParamsRaw, &j); err != nil {
			return err
		}
		*job = DockerHostTaskJob{BasicJob: &bj, dockerHostTaskJob: &j}
		return nil
	} else {
		return err
	}

}

func (job DockerHostTaskJob) GetBasicJob() BasicJob {
	return *job.BasicJob
}

// ----------------

type dockerServiceTaskJob struct {
	Task        string `json:"task"`
	ServiceName string `json:"service_name"`
	PrivateIp   string `json:"private_ip"`
}

type DockerServiceTaskJob struct {
	*BasicJob
	*dockerServiceTaskJob
}

func (job *DockerServiceTaskJob) UnmarshalJSON(b []byte) error {
	var bj BasicJob
	if err := json.Unmarshal(b, &bj); err == nil {
		var j dockerServiceTaskJob
		if err := json.Unmarshal(bj.basicJob.ParamsRaw, &j); err != nil {
			return err
		}
		*job = DockerServiceTaskJob{BasicJob: &bj, dockerServiceTaskJob: &j}
		return nil
	} else {
		return err
	}

}

func (job DockerServiceTaskJob) GetBasicJob() BasicJob {
	return *job.BasicJob
}

// ----------------

func (c *Client) GetJobs(stackUid string, serverUid *string) ([]Job, error) {
	fmt.Printf("")
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

	var p Pagination
	var result []Job
	var jobRes []*json.RawMessage

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

		for _, j := range jobRes {
			var job *Job
			if job, err = JobFactory(*j); err != nil {
				return nil, err
			}
			result = append(result, *job)
		}

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
	var jobRes *json.RawMessage
	err = c.DoReq(req, &jobRes, nil)
	if err != nil {
		return nil, err
	}
	return JobFactory(*jobRes)
}

func JobFactory(jobRes json.RawMessage) (*Job, error) {
	var T = struct {
		Type string `json:"type"`
	}{}

	if err := json.Unmarshal(jobRes, &T); err != nil {
		return nil, err
	}

	var job Job

	switch T.Type {
	case "DockerHostTaskJob":
		job = new(DockerHostTaskJob)
	case "DockerServiceTaskJob":
		job = new(DockerServiceTaskJob)
	default:
		job = new(BasicJob)
	}

	if err := json.Unmarshal(jobRes, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

func (c *Client) RunJobNow(stackUid string, jobUid string, jobVars *string) (*AsyncResult, error) {
	var params interface{}
	if jobVars == nil {
		params = nil
	} else {
		params = struct {
			JobVars string `json:"job_vars"`
		}{
			JobVars: *jobVars,
		}
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/jobs/"+jobUid+"run_now.json", params, nil)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes, nil)
}
