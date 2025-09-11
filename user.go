package zabbix

type (
	// Whether to pause escalation during maintenance periods or not.
	// "debug_mode" in https://www.zabbix.com/documentation/4.0/manual/api/reference/user/object#user_group
	UserType int
)

const (
	ZabbixUser       UserType = 0
	ZabbixAdmin      UserType = 1
	ZabbixSuperAdmin UserType = 2
)

type (
	// Media status.
	MediaStatus int
)

const (
	MediaStatusActive   MediaStatus = 0
	MediaStatusDisabled MediaStatus = 1
)

// User represent Zabbix user object
// https://www.zabbix.com/documentation/current/manual/api/reference/user/object
type User struct {
	UserID          string `json:"userid,omitempty"`
	Username        string `json:"username"`
	Name            string `json:"name,omitempty"`
	Surname         string `json:"surname,omitempty"`
	Url             string `json:"url,omitempty"`
	Autologin       int    `json:"autologin,string,omitempty"`
	Autologout      string `json:"autologout,omitempty"`
	Lang            string `json:"lang,omitempty"`
	Refresh         string `json:"refresh,omitempty"`
	Theme           string `json:"theme,omitempty"`
	AttemptFailed   int    `json:"attempt_failed,string,omitempty"`
	AttemptIP       string `json:"attempt_ip,omitempty"`
	AttemptClock    int    `json:"attempt_clock,string,omitempty"`
	RowsPerPage     int    `json:"rows_per_page,string,omitempty"`
	Timezone        string `json:"timezone,omitempty"`
	RoleID          string `json:"roleid,omitempty"`
	UserDirectoryID string `json:"userdirectoryid,omitempty"`
	TsProvisioned   int    `json:"ts_provisioned,string,omitempty"`
	Provisioned     int    `json:"provisioned,string,omitempty"`

	// Legacy field for backward compatibility
	Alias string   `json:"alias,omitempty"`
	Type  UserType `json:"type,string,omitempty"`
}

// Media represents Zabbix user media object
// https://www.zabbix.com/documentation/current/manual/api/reference/user/object#media
type Media struct {
	MediaID              string      `json:"mediaid,omitempty"`
	MediaTypeID          string      `json:"mediatypeid"`
	SendTo               []string    `json:"sendto"`
	Active               MediaStatus `json:"active,string,omitempty"`
	Severity             int         `json:"severity,string,omitempty"`
	Period               string      `json:"period,omitempty"`
	UserDirectoryMediaID string      `json:"userdirectory_mediaid,omitempty"`
	Provisioned          int         `json:"provisioned,string,omitempty"`
}

// Users is an array of User
type Users []User

// Medias is an array of Media
type Medias []Media

// UsersGet Wrapper for user.get
// https://www.zabbix.com/documentation/4.0/manual/api/reference/user/get
func (api *API) UsersGet(params Params) (res Users, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	err = api.CallWithErrorParse("user.get", params, &res)
	return
}
