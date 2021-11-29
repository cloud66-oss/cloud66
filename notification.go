package cloud66

type NotificationUploadParams struct {
	Alerts              []Notification `json:"alerts"`
	TargetStackUid      string         `json:"dest_stack_id,omitempty"`
	ApplicationGroupUid string         `json:"application_group_id,omitempty"`
}

type NotificationSubscription struct {
	Channel string `json:"channel"`
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
	Successes NotificationResponse        `json:"successes"`
	Failures  NotificationResponseFailure `json:"failures"`
}

func (c *Client) NotificationDownload(stackUid string) ([]Notification, error) {
	var notifications []Notification

	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/alerts", nil, nil)
	if err != nil {
		return nil, err
	}

	err = c.DoReq(req, &notifications, nil)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (c *Client) NotificationUploadStack(stackUid string, targetStackUid string, alerts []Notification) (*NotificationResponseBody, error) {
	var notification NotificationUploadParams
	notification.Alerts = alerts
	notification.TargetStackUid = targetStackUid
	return c.NotificationUpload(stackUid, notification)
}

func (c *Client) NotificationUploadAG(stackUid string, targetUid string, alerts []Notification) (*NotificationResponseBody, error) {
	var notification NotificationUploadParams
	notification.Alerts = alerts
	notification.ApplicationGroupUid = targetUid
	return c.NotificationUpload(stackUid, notification)
}

func (c *Client) NotificationUpload(stackUid string, notification NotificationUploadParams) (*NotificationResponseBody, error) {
	var notifications NotificationResponseBody

	req, err := c.NewRequest("PATCH", "/stacks/"+stackUid+"/alerts", notification, nil)
	if err != nil {
		return nil, err
	}

	err = c.DoReq(req, &notifications, nil)
	if err != nil {
		return nil, err
	}
	return &notifications, nil
}
