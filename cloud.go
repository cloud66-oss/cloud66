package cloud66

type Cloud struct {
	Id          string            `json:"name"`
	Name        string            `json:"display_name"`
	Regions     []CloudRegion     `json:"regions"`
	ServerSizes []CloudServerSize `json:"server_sizes"`
}

type CloudServerSize struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CloudRegion struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) GetCloudInfo(cloudName string) (*Cloud, error) {
	req, err := c.NewRequest("GET", "/clouds/"+cloudName+".json", nil, nil)
	if err != nil {
		return nil, err
	}
	var cloudRes *Cloud
	return cloudRes, c.DoReq(req, &cloudRes, nil)
}
