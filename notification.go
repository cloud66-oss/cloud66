package cloud66

type NotificationUploadParams struct {
	Alerts               []Notification `json:"alerts"`
	ApplicationGroupName string         `json:"application_group_name,omitempty"`
}

type NotificationSubscription struct {
	Channel   string `json:"channel"`
	SlackUrl  string `json:"slack_url,omitempty" yaml:"slack_url,omitempty"`
	WebookUrl string `json:"webhook_url,omitempty"  yaml:"webhook_url,omitempty"`
}

type Notification struct {
	Name          string                     `json:"alert_name"`
	Subscriptions []NotificationSubscription `json:"subscriptions,omitempty"`
}

type NotificationResponse struct {
	Alerts []string `json:"alerts"`
	Count  int      `json:"count"`
}

type NotificationResponseFailureBody struct {
	Alert  string `json:"alert_name"`
	Reason string `json:"reason"`
}

type NotificationResponseFailure struct {
	Alerts []NotificationResponseFailureBody `json:"alerts"`
	Count  int                               `json:"count"`
}

type NotificationResponseBody struct {
	Successes     NotificationResponse        `json:"successes"`
	NotApplicable NotificationResponseFailure `json:"not_applicable"`
	Failures      NotificationResponseFailure `json:"failures"`
}

func (c *Client) NotificationDownload(stackUid string) ([]Notification, error) {
	var notifications []Notification

	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/notifications", nil, nil)
	if err != nil {
		return nil, err
	}

	err = c.DoReq(req, &notifications, nil)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (c *Client) NotificationUploadStack(targetStackUid string, alerts []Notification) (*NotificationResponseBody, error) {
	var notificationUploadParams NotificationUploadParams
	notificationUploadParams.Alerts = alerts
	return c.NotificationUpload(notificationUploadParams, &targetStackUid)
}

func (c *Client) NotificationUploadApplicationGroup(targetUid string, alerts []Notification) (*NotificationResponseBody, error) {
	var notificationUploadParams NotificationUploadParams
	notificationUploadParams.Alerts = alerts
	notificationUploadParams.ApplicationGroupName = targetUid
	return c.NotificationUpload(notificationUploadParams, nil)
}

func (c *Client) NotificationUpload(notification NotificationUploadParams, targetStackUid *string) (*NotificationResponseBody, error) {
	var notifications NotificationResponseBody
	var requestPath string

	if targetStackUid != nil {
		requestPath = "/stacks/" + *targetStackUid + "/notifications"
	} else {
		requestPath = "/application_groups/notifications"
	}
	req, err := c.NewRequest("PATCH", requestPath, notification, nil)
	if err != nil {
		return nil, err
	}

	err = c.DoReq(req, &notifications, nil)
	if err != nil {
		return nil, err
	}
	return &notifications, nil
}
