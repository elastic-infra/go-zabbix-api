package zabbix

type (
	// Role type.
	// "type" in https://www.zabbix.com/documentation/current/manual/api/reference/role/object
	RoleType int
)

const (
	RoleTypeUser       RoleType = 1
	RoleTypeAdmin      RoleType = 2
	RoleTypeSuperAdmin RoleType = 3
)

// Role represents Zabbix role object
// https://www.zabbix.com/documentation/current/manual/api/reference/role/object
type Role struct {
	RoleID   string   `json:"roleid,omitempty"`
	Name     string   `json:"name"`
	Type     RoleType `json:"type,string,omitempty"`
	ReadOnly int      `json:"readonly,string,omitempty"`
}

// Roles is an array of Role
type Roles []Role

// RoleGet Wrapper for role.get
// https://www.zabbix.com/documentation/current/manual/api/reference/role/get
func (api *API) RoleGet(params Params) (res Roles, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	err = api.CallWithErrorParse("role.get", params, &res)
	return
}
