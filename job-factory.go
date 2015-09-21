package cloud66

import "encoding/json"

type Job interface {
	GetBasicJob() BasicJob
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

// -------- DockerHostTaskJob --------

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

// -------- DockerServiceTaskJob --------

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
