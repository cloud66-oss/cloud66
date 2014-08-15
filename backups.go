package cloud66

import (
	"strconv"
	"time"
)

/*
 */
var BackupStatus = map[int]string{
	0: "Ok",     // BCK_OK
	1: "Failed", // BCK_FAILED
}

var RestoreStatus = map[int]string{
	0: "Not Restored",   // RST_NA
	1: "Restoring",      // RST_RESTORING
	2: "Restored",       // RST_OK
	3: "Restore Failed", // RST_FAILED
}

var VerifyStatus = map[int]string{
	0: "Not Verified",        // VRF_NA
	1: "Verifying",           // VRF_VERIFING
	2: "Verified",            // VRF_OK
	3: "Verification Failed", // VRF_FAILED
	4: "Unable to Verify",    // VRF_INTERNAL_ISSUE
}

type ManagedBackup struct {
	Id            int       `json:"id"`
	ServerUid     string    `json:"server_uid"`
	Filename      string    `json:"file_name"`
	DbType        string    `json:"db_type"`
	DatabaseName  string    `json:"database_name"`
	FileBase      string    `json:"file_base"`
	BackupDate    time.Time `json:"backup_date_iso"`
	BackupStatus  int       `json:"backup_status"`
	BackupResult  string    `json:"backup_result"`
	RestoreStatus int       `json:"restore_status"`
	RestoreResult string    `json:"restore_result"`
	CreatedAt     time.Time `json:"created_at_iso"`
	UpdatedAt     time.Time `json:"updated_at_iso"`
	VerifyStatus  int       `json:"verify_status"`
	VerifyResult  string    `json:"verify_result"`
	StoragePath   string    `json:"storage_path"`
	SkipTables    string    `json:"skip_tables"`
}

type BackupSegmentIndex struct {
	Filename  string `json:"name"`
	Extension string `json:"id"`
}

type BackupSegment struct {
	Ok  bool   `json:"ok"`
	Url string `json:"public_url"`
}

func (c *Client) GetBackupSegmentIndeces(stackUid string, backupId int) ([]BackupSegmentIndex, error) {
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/backups/"+strconv.Itoa(backupId)+"/files.json", nil)
	if err != nil {
		return nil, err
	}
	var backupSegIndex []BackupSegmentIndex
	return backupSegIndex, c.DoReq(req, &backupSegIndex)
}

func (c *Client) GetBackupSegment(stackUid string, backupId int, extension string) (*BackupSegment, error) {
	ext := ""
	if extension != "" {
		ext = "/" + extension
	}
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/backups/"+strconv.Itoa(backupId)+"/files/"+ext+".json", nil)
	if err != nil {
		return nil, err
	}
	var backupSegmentRes *BackupSegment
	return backupSegmentRes, c.DoReq(req, &backupSegmentRes)
}
