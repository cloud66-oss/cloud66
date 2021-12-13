package cloud66

import (
	"strings"
	"time"
)

type CurrentStackType int

const (
	STACK_PRIMARY   CurrentStackType = 1
	STACK_SECONDARY CurrentStackType = 2
)

type FailoverGroup struct {
	Uid                string           `json:"uid"`
	Address            string           `json:"address"`
	PrimaryStackName   string           `json:"primary_stack_name"`
	PrimaryStackUid    string           `json:"primary_stack_uid"`
	SecondaryStackName string           `json:"secondary_stack_name"`
	SecondaryStackUid  string           `json:"secondary_stack_uid"`
	CurrentStack       CurrentStackType `json:"current_stack"`
	BusyToggling       bool             `json:"busy_toggling"`
	Readonly           bool             `json:"readonly"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}

func (c *Client) FailoverGroupList() ([]FailoverGroup, error) {
	queryStrings := make(map[string]string)

	var result []FailoverGroup

	req, err := c.NewRequest("GET", "/failover_groups.json", nil, queryStrings)
	if err != nil {
		return nil, err
	}

	err = c.DoReq(req, &result, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) AddFailoverGroup(primaryStack *string, secondaryStack *string, currentStack *CurrentStackType) error {
	params := struct {
		PrimaryStack   *string           `json:"primary_stack_uid"`
		SecondaryStack *string           `json:"secondary_stack_uid"`
		CurrentStack   *CurrentStackType `json:"current_stack"`
	}{
		PrimaryStack:   primaryStack,
		SecondaryStack: secondaryStack,
		CurrentStack:   currentStack,
	}

	req, err := c.NewRequest("POST", "/failover_groups", params, nil)
	if err != nil {
		return err
	}
	return c.DoReq(req, nil, nil)
}

func (c *Client) UpdateFailoverGroup(failoverGroupUid string, primaryStack *string, secondaryStack *string, currentStack *CurrentStackType) error {
	params := struct {
		PrimaryStack   *string           `json:"primary_stack_uid"`
		SecondaryStack *string           `json:"secondary_stack_uid"`
		CurrentStack   *CurrentStackType `json:"current_stack"`
	}{
		PrimaryStack:   primaryStack,
		SecondaryStack: secondaryStack,
		CurrentStack:   currentStack,
	}

	req, err := c.NewRequest("PUT", "/failover_groups/"+failoverGroupUid, params, nil)
	if err != nil {
		return err
	}

	return c.DoReq(req, nil, nil)
}

func (c *Client) DeleteFailoverGrouop(failoverGroupUid string) error {
	req, err := c.NewRequest("DELETE", "/failover_groups/"+failoverGroupUid, nil, nil)
	if err != nil {
		return err
	}
	return c.DoReq(req, nil, nil)
}

func (currentStackType *CurrentStackType) String() string {
	switch *currentStackType {
	case STACK_PRIMARY:
		return "Primary"
	case STACK_SECONDARY:
		return "Secondary"
	}
	return "Wrong CurrentStackType"
}

func ParseCurrentStack(param string) CurrentStackType {
	if strings.EqualFold(param, "primary") {
		return STACK_PRIMARY
	}

	if strings.EqualFold(param, "secondary") {
		return STACK_SECONDARY
	}
	return STACK_PRIMARY
}
