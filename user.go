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
	UserID          string     `json:"userid,omitempty"`
	Username        string     `json:"username,omitempty"`
	Name            string     `json:"name,omitempty"`
	Surname         string     `json:"surname,omitempty"`
	Password        string     `json:"passwd,omitempty"`
	CurrentPassword string     `json:"current_passwd,omitempty"`
	Url             string     `json:"url,omitempty"`
	Autologin       int        `json:"autologin,string,omitempty"`
	Autologout      string     `json:"autologout,omitempty"`
	Lang            string     `json:"lang,omitempty"`
	Refresh         string     `json:"refresh,omitempty"`
	Theme           string     `json:"theme,omitempty"`
	AttemptFailed   int        `json:"attempt_failed,string,omitempty"`
	AttemptIP       string     `json:"attempt_ip,omitempty"`
	AttemptClock    int        `json:"attempt_clock,string,omitempty"`
	RowsPerPage     int        `json:"rows_per_page,string,omitempty"`
	Timezone        string     `json:"timezone,omitempty"`
	RoleID          string     `json:"roleid,omitempty"`
	UserDirectoryID string     `json:"userdirectoryid,omitempty"`
	TsProvisioned   int        `json:"ts_provisioned,string,omitempty"`
	Provisioned     int        `json:"provisioned,string,omitempty"`
	Medias          Medias     `json:"medias,omitempty"`
	UsrGrps         UserGroups `json:"usrgrps,omitempty"`

	// Legacy field for backward compatibility
	Alias string   `json:"alias,omitempty"`
	Type  UserType `json:"type,string,omitempty"`
}

// Media represents Zabbix user media object
// https://www.zabbix.com/documentation/current/manual/api/reference/user/object#media
type Media struct {
	MediaID              string      `json:"mediaid,omitempty"`
	MediaTypeID          string      `json:"mediatypeid"`
	SendTo               interface{} `json:"sendto"`
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

// UserCreate Wrapper for user.create
// https://www.zabbix.com/documentation/current/manual/api/reference/user/create
func (api *API) UserCreate(users Users) (err error) {
	response, err := api.CallWithError("user.create", users)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	userids := result["userids"].([]interface{})
	for i, userid := range userids {
		users[i].UserID = userid.(string)
	}
	return
}

// UserUpdate Wrapper for user.update
// https://www.zabbix.com/documentation/current/manual/api/reference/user/update
func (api *API) UserUpdate(users Users) (err error) {
	_, err = api.CallWithError("user.update", users)
	return
}

// UserDelete Wrapper for user.delete
// https://www.zabbix.com/documentation/current/manual/api/reference/user/delete
func (api *API) UserDelete(userids []string) (err error) {
	response, err := api.CallWithError("user.delete", userids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	if deletedids, ok := result["userids"].([]interface{}); ok && len(deletedids) != len(userids) {
		err = &Error{Code: -32602, Message: "User delete failed", Data: "Could not delete all users"}
	}
	return
}
