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

type BackupSegment struct {
	Ok            bool   `json:"ok"`
	Filename      string `json:"file_name"`
	Url           string `json:"url"`
	NextExtension string `json:"next_extension"`
}

func (c *Client) GetBackupSegment(backupId int, extension string) (*BackupSegment, error) {
	ext := ""
	if extension != "" {
		ext = "/" + extension
	}
	req, err := c.NewRequest("GET", "/backups/"+strconv.Itoa(backupId)+"/export"+ext+".json", nil)
	if err != nil {
		return nil, err
	}

	var backupSegmentRes *BackupSegment
	return backupSegmentRes, c.DoReq(req, &backupSegmentRes)
}
