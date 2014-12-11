package cloud66

import (
	"errors"
	"strings"
	"time"
)

var stackStatus = map[int]string{
	0: "Pending analysis",      //STK_QUEUED
	1: "Deployed successfully", //STK_SUCCESS
	2: "Deployment failed",     //STK_FAILED
	3: "Analyzing",             //STK_ANALYSING
	4: "Analyzed",              //STK_ANALYSED
	5: "Queued for deployment", //STK_QUEUED_FOR_DEPLOYING
	6: "Deploying",             //STK_DEPLOYING
	7: "Unable to analyze",     //STK_TERMINAL_FAILURE
}

var healthStatus = map[int]string{
	0: "Unknown",  //HLT_UNKNOWN
	1: "Building", //HLT_BUILDING
	2: "Impaired", //HLT_PARTIAL
	3: "Healthy",  //HLT_OK
	4: "Failed",   //HLT_BROKEN
}

type Stack struct {
	Uid             string     `json:"uid"`
	Name            string     `json:"name"`
	Git             string     `json:"git"`
	GitBranch       string     `json:"git_branch"`
	Environment     string     `json:"environment"`
	Cloud           string     `json:"cloud"`
	Fqdn            string     `json:"fqdn"`
	Language        string     `json:"language"`
	Framework       string     `json:"framework"`
	StatusCode      int        `json:"status"`
	HealthCode      int        `json:"health"`
	MaintenanceMode bool       `json:"maintenance_mode"`
	HasLoadBalancer bool       `json:"has_loadbalancer"`
	RedeployHook    *string    `json:"redeploy_hook"`
	LastActivity    *time.Time `json:"last_activity_iso"`
	UpdatedAt       time.Time  `json:"updated_at_iso"`
	CreatedAt       time.Time  `json:"created_at_iso"`
	DeployDir       string     `json:"deploy_directory"`
}

type StackSetting struct {
	Key      string      `json:"key"`
	Value    interface{} `json:"value"`
	Readonly bool        `json:"readonly"`
}

type StackEnvVar struct {
	Key      string      `json:"key"`
	Value    interface{} `json:"value"`
	Readonly bool        `json:"readonly"`
}

func (s Stack) Status() string {
	return stackStatus[s.StatusCode]
}

func (s Stack) Health() string {
	return healthStatus[s.HealthCode]
}

func (c *Client) StackList() ([]Stack, error) {
	req, err := c.NewRequest("GET", "/stacks.json", nil)
	if err != nil {
		return nil, err
	}

	var stacksRes []Stack
	return stacksRes, c.DoReq(req, &stacksRes)
}

func (c *Client) StackListWithFilter(filter filterFunction) ([]Stack, error) {
	req, err := c.NewRequest("GET", "/stacks.json", nil)
	if err != nil {
		return nil, err
	}

	var stacksRes []Stack
	err = c.DoReq(req, &stacksRes)
	if err != nil {
		return nil, err
	}

	var result []Stack
	for _, item := range stacksRes {
		if filter(item) {
			result = append(result, item)
		}
	}
	return result, nil
}

func (c *Client) StackInfo(stackName string) (*Stack, error) {
	stack, err := c.FindStackByName(stackName, "")
	if err != nil {
		return nil, err
	}

	uid := stack.Uid
	req, err := c.NewRequest("GET", "/stacks/"+uid+".json", nil)
	if err != nil {
		return nil, err
	}

	var stacksRes *Stack
	return stacksRes, c.DoReq(req, &stacksRes)
}

func (c *Client) StackInfoWithEnvironment(stackName, environment string) (*Stack, error) {
	stack, err := c.FindStackByName(stackName, environment)
	if err != nil {
		return nil, err
	}

	uid := stack.Uid
	req, err := c.NewRequest("GET", "/stacks/"+uid+".json", nil)
	if err != nil {
		return nil, err
	}

	var stacksRes *Stack
	return stacksRes, c.DoReq(req, &stacksRes)
}

func (c *Client) StackSettings(uid string) ([]StackSetting, error) {
	req, err := c.NewRequest("GET", "/stacks/"+uid+"/settings.json", nil)
	if err != nil {
		return nil, err
	}

	var settingsRes []StackSetting
	return settingsRes, c.DoReq(req, &settingsRes)
}

func (c *Client) StackEnvVars(uid string) ([]StackEnvVar, error) {
	req, err := c.NewRequest("GET", "/stacks/"+uid+"/environments.json", nil)
	if err != nil {
		return nil, err
	}

	var envVarsRes []StackEnvVar
	return envVarsRes, c.DoReq(req, &envVarsRes)
}

func (c *Client) StackEnvVarNew(stackUid string, key string, value string) (*AsyncResult, error) {
	params := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{
		Key:   key,
		Value: value,
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/environments.json", params)
	if err != nil {
		return nil, err
	}
	var asyncResult *AsyncResult
	return asyncResult, c.DoReq(req, &asyncResult)
}

func (c *Client) StackEnvVarSet(stackUid string, key string, value string) (*AsyncResult, error) {
	params := struct {
		Value string `json:"value"`
	}{
		Value: value,
	}
	req, err := c.NewRequest("PUT", "/stacks/"+stackUid+"/environments/"+key+".json", params)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes)
}

func (c *Client) FindStackByName(stackName, environment string) (*Stack, error) {
	stacks, err := c.StackList()

	for _, b := range stacks {
		if (strings.ToLower(b.Name) == strings.ToLower(stackName)) && (environment == "" || environment == b.Environment) {
			return &b, err
		}
	}

	return nil, errors.New("Stack not found")
}

func (c *Client) ManagedBackups(uid string) ([]ManagedBackup, error) {
	req, err := c.NewRequest("GET", "/stacks/"+uid+"/backups.json", nil)
	if err != nil {
		return nil, err
	}

	var managedBackupsRes []ManagedBackup
	return managedBackupsRes, c.DoReq(req, &managedBackupsRes)
}

func (c *Client) Set(uid string, key string, value string) (*AsyncResult, error) {
	key = strings.Replace(key, ".", "-", -1)
	params := struct {
		Value string `json:"value"`
	}{
		Value: value,
	}
	req, err := c.NewRequest("PUT", "/stacks/"+uid+"/settings/"+key+".json", params)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes)
}

func (c *Client) Lease(uid string, ipAddress *string, timeToOpen *int, port *int, serverUid *string) (*AsyncResult, error) {
	var (
		theIpAddress  *string
		theTimeToOpen *int
		thePort       *int
		theServerUid  *string
	)
	// set defaults
	if ipAddress == nil {
		var value = "AUTO"
		theIpAddress = &value
	} else {
		theIpAddress = ipAddress
	}
	if timeToOpen == nil {
		var value = 20
		theTimeToOpen = &value
	} else {
		theTimeToOpen = timeToOpen
	}
	if port == nil {
		var value = 22
		thePort = &value
	} else {
		thePort = port
	}
	if serverUid == nil {
		var value = ""
		theServerUid = &value
	} else {
		theServerUid = serverUid
	}

	params := struct {
		TimeToOpen *int    `json:"ttl"`
		IpAddress  *string `json:"from_ip"`
		Port       *int    `json:"port"`
		ServerUid  *string `json:"server_id"`
	}{
		TimeToOpen: theTimeToOpen,
		IpAddress:  theIpAddress,
		Port:       thePort,
		ServerUid:  theServerUid,
	}
	req, err := c.NewRequest("POST", "/stacks/"+uid+"/firewalls.json", params)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes)
}

func (c *Client) LeaseSync(stackUid string, ipAddress *string, timeToOpen *int, port *int, serverUid *string) (*GenericResponse, error) {
	asyncRes, err := c.Lease(stackUid, ipAddress, timeToOpen, port, serverUid)
	if err != nil {
		return nil, err
	}
	genericRes, err := c.WaitStackAsyncAction(asyncRes.Id, stackUid, 2*time.Second, 5*time.Minute, false)
	if err != nil {
		return nil, err
	}
	return genericRes, err
}

func (c *Client) RedeployStack(stackUid string, gitRef string) (*GenericResponse, error) {
	params := struct {
		GitRef string `json:"git_ref"`
	}{
		GitRef: gitRef,
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/deployments.json", params)
	if err != nil {
		return nil, err
	}
	var stacksRes *GenericResponse
	return stacksRes, c.DoReq(req, &stacksRes)
}

func (c *Client) InvokeStackAction(stackUid string, action string) (*AsyncResult, error) {
	params := struct {
		Command string `json:"command"`
	}{
		Command: action,
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/actions.json", params)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes)
}

func (c *Client) InvokeDbStackAction(stackUid string, serverUid string, dbType *string, action string) (*AsyncResult, error) {
	var params interface{}
	if dbType == nil {
		params = struct {
			Command   string `json:"command"`
			ServerUid string `json:"server_uid"`
		}{
			Command:   action,
			ServerUid: serverUid,
		}
	} else {
		params = struct {
			Command   string `json:"command"`
			ServerUid string `json:"server_uid"`
			DbType    string `json:"db_type"`
		}{
			Command:   action,
			ServerUid: serverUid,
			DbType:    *dbType,
		}
	}
	req, err := c.NewRequest("POST", "/stacks/"+stackUid+"/actions.json", params)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes)
}
