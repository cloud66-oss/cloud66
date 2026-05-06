package cloud66

import (
	"strconv"
	"time"
)

// OperationLogSeverity values mirror the server-side LogUtils scale.
const (
	OperationLogSeverityTrace  = 0
	OperationLogSeverityDebug  = 1
	OperationLogSeverityInfo   = 2
	OperationLogSeverityNotice = 3
	OperationLogSeverityWarn   = 4
	OperationLogSeverityError  = 5
)

// OperationLogEntry is a single log entry returned by the operations logs endpoint.
// json keys are intentionally terse — they mirror the server's OperationLogEntity shape.
type OperationLogEntry struct {
	// severity level (0..5); see OperationLogSeverity* constants
	Severity int `json:"v"`
	// human-readable log message
	Message string `json:"m"`
	// timestamp in RFC 3339 (e.g. 2025-01-15T10:00:05Z)
	Timestamp time.Time `json:"t"`
	// source operation when nested under a parent; nil for root entries
	Source *string `json:"s"`
	// echelon — depth from the root operation (0 = root)
	Echelon int `json:"e"`
	// uid of the child operation that produced the entry; only present when
	// include_children=true and the entry came from a descendant, nil otherwise
	OperationUid *string `json:"operation_uid,omitempty"`
}

// OperationLogs returns all log entries for the operation identified by its UID.
// Works for any async operation (deployments, ssl, load balancer, archive/restore, etc.).
//
// minSeverity, when non-nil, filters out entries below the given severity (0..5);
// pass nil for "no filter". includeChildren aggregates logs from child operations.
// auto-paginates through every page server-side.
func (c *Client) OperationLogs(operationUid string, minSeverity *int, includeChildren bool) ([]OperationLogEntry, error) {
	queryStrings := make(map[string]string)
	queryStrings["page"] = "1"
	// only set min_severity when caller asked for one — server treats absent as "no filter"
	if minSeverity != nil {
		queryStrings["min_severity"] = strconv.Itoa(*minSeverity)
	}
	// only set include_children when true; server defaults to false
	if includeChildren {
		queryStrings["include_children"] = "true"
	}

	var p Pagination
	var result []OperationLogEntry
	var operationLogRes []OperationLogEntry

	for {
		req, err := c.NewRequest("GET", "/operations/"+operationUid+"/logs.json", nil, queryStrings)
		if err != nil {
			return nil, err
		}

		operationLogRes = nil
		err = c.DoReq(req, &operationLogRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, operationLogRes...)
		// stop when there's no next page (current page == last page)
		if p.Current < p.Next {
			queryStrings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}
	}

	return result, nil
}
