package zabbix

type (
	// Media type.
	// "type" in https://www.zabbix.com/documentation/current/manual/api/reference/mediatype/object
	MediaTypeType int
)

const (
	MediaTypeEmail   MediaTypeType = 0
	MediaTypeScript  MediaTypeType = 1
	MediaTypeSMS     MediaTypeType = 2
	MediaTypeWebhook MediaTypeType = 4
)

// MediaType represents Zabbix media type object
// https://www.zabbix.com/documentation/current/manual/api/reference/mediatype/object
//
// Version compatibility notes:
// - Zabbix 7.0+ uses MessageFormat field
// - Zabbix 6.4 and earlier use ContentType field
// - Both fields are supported for maximum compatibility
// - Use GetMessageFormat() and SetMessageFormat() methods for version-agnostic access
type MediaType struct {
	MediaTypeID        string        `json:"mediatypeid,omitempty"`
	Type               MediaTypeType `json:"type,string,omitempty"`
	Name               string        `json:"name"`
	SMTPServer         string        `json:"smtp_server,omitempty"`
	SMTPHelo           string        `json:"smtp_helo,omitempty"`
	SMTPEmail          string        `json:"smtp_email,omitempty"`
	SMTPAuthentication int           `json:"smtp_authentication,string,omitempty"`
	MaxSessions        int           `json:"maxsessions,string,omitempty"`
	MaxAttempts        int           `json:"maxattempts,string,omitempty"`
	AttemptInterval    string        `json:"attempt_interval,omitempty"`
	MessageFormat      int           `json:"message_format,string,omitempty"`
	ContentType        int           `json:"content_type,string,omitempty"` // Legacy field for v6.4 and earlier compatibility
	Script             string        `json:"script,omitempty"`
	Timeout            string        `json:"timeout,omitempty"`
	ProcessTags        int           `json:"process_tags,string,omitempty"`
	ShowEventMenu      int           `json:"show_event_menu,string,omitempty"`
	EventMenuURL       string        `json:"event_menu_url,omitempty"`
	EventMenuName      string        `json:"event_menu_name,omitempty"`
	Description        string        `json:"description,omitempty"`
	Status             int           `json:"status,string,omitempty"`
	Provider           int           `json:"provider,string,omitempty"`

	// OAuth fields (for certain providers)
	RedirectionURL     string `json:"redirection_url,omitempty"`
	ClientID           string `json:"client_id,omitempty"`
	ClientSecret       string `json:"client_secret,omitempty"`
	AuthorizationURL   string `json:"authorization_url,omitempty"`
	TokenURL           string `json:"token_url,omitempty"`
	TokensStatus       int    `json:"tokens_status,string,omitempty"`
	AccessToken        string `json:"access_token,omitempty"`
	AccessTokenUpdated int    `json:"access_token_updated,string,omitempty"`
	AccessExpiresIn    int    `json:"access_expires_in,string,omitempty"`
	RefreshToken       string `json:"refresh_token,omitempty"`
}

// GetMessageFormat returns the message format value, preferring MessageFormat over ContentType
func (mt *MediaType) GetMessageFormat() int {
	if mt.MessageFormat != 0 {
		return mt.MessageFormat
	}
	return mt.ContentType
}

// SetMessageFormat sets both MessageFormat and ContentType for maximum compatibility
func (mt *MediaType) SetMessageFormat(format int) {
	mt.MessageFormat = format
	mt.ContentType = format
}

// MediaTypes is an array of MediaType
type MediaTypes []MediaType

// MediaTypeGet Wrapper for mediatype.get
// https://www.zabbix.com/documentation/current/manual/api/reference/mediatype/get
func (api *API) MediaTypeGet(params Params) (res MediaTypes, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	err = api.CallWithErrorParse("mediatype.get", params, &res)
	return
}
