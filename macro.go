package zabbix

import "github.com/AlekSi/reflector"

// Macro represent Zabbix User MAcro object
// https://www.zabbix.com/documentation/3.2/manual/api/reference/usermacro/object
type Macro struct {
	MacroID   string `json:"hostmacroids,omitempty"`
	HostID    string `json:"hostid"`
	MacroName string `json:"macro"`
	Value     string `json:"value"`
}

// Macros is an array of Macro
type Macros []Macro

// MacroGet Wrapper for usermacro.get
// https://www.zabbix.com/documentation/3.2/manual/api/reference/usermacro/get
func (api *API) MacroGet(params Params) (res Macros, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("usermacro.get", params)
	if err != nil {
		return
	}
	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}

// MacroGetByID Get macro by macro ID if there is exactly 1 matching macro
func (api *API) MacroGetByID(id string) (res *Macro, err error) {
	triggers, err := api.MacroGet(Params{"hostmacroids": id})
	if err != nil {
		return
	}

	if len(triggers) == 1 {
		res = &triggers[0]
	} else {
		e := ExpectedOneResult(len(triggers))
		err = &e
	}
	return
}

// MacroCreate Wrapper for usermacro.create
// https://www.zabbix.com/documentation/3.2/manual/api/reference/usermacro/create
func (api *API) MacroCreate(macros Macros) error {
	response, err := api.CallWithError("usermacro.create", macros)
	if err != nil {
		return err
	}

	result := response.Result.(map[string]interface{})
	macroids := result["hostmacroids"].([]interface{})
	for i, id := range macroids {
		macros[i].HostID = id.(string)
	}
	return nil
}

// MacroUpdate Wrapper for usermacro.update
// https://www.zabbix.com/documentation/3.2/manual/api/reference/usermacro/update
func (api *API) MacroUpdate(macros Macros) (err error) {
	_, err = api.CallWithError("usermacro.create", macros)
	return
}

// MacroDeleteByID Wrapper for usermacro.delete
// Cleans MacroId in all macro elements if call succeed.
//https://www.zabbix.com/documentation/3.2/manual/api/reference/usermacro/delete
func (api *API) MacroDeleteByID(ids []string) (err error) {
	response, err := api.CallWithError("usermacro.delete", ids)

	result := response.Result.(map[string]interface{})
	hostmacroids := result["hostmacroids"].([]interface{})
	if len(ids) != len(hostmacroids) {
		err = &ExpectedMore{len(ids), len(hostmacroids)}
	}
	return
}

// MacroDelete Wrapper for usermacro.delete
// https://www.zabbix.com/documentation/3.2/manual/api/reference/usermacro/delete
func (api *API) MacroDelete(macros Macros) (err error) {
	ids := make([]string, len(macros))
	for i, macro := range macros {
		ids[i] = macro.MacroID
	}

	err = api.MacroDeleteByID(ids)
	if err == nil {
		for i := range macros {
			macros[i].MacroID = ""
		}
	}
	return
}
