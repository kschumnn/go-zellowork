package types

type ZelloUser struct {
	Name          string
	EMail         string
	Admin         bool
	LimitedAccess bool `json:"limited_access"`
	Job           string
	FullName      string `json:"full_name"`
	Channels      []string
}

type ZelloChannel struct {
	Name        string `json:"name"`
	Created     string `json:"created"`
	Count       string `json:"count"`
	IsShared    int    `json:"is_shared"`
	IsInvisible int    `json:"is_invisible"`
}

type ZelloChannelRoleSettings struct {
	ListenOnly   bool     `json:"listen_only"`
	NoDisconnect bool     `json:"no_disconnect"`
	To           []string `json:"to"`
	AllowAlerts  bool     `json:"allow_alerts"`
}
type ZelloChannelRole struct {
	Name     string                   `json:"name"`
	Settings ZelloChannelRoleSettings `json:"settings"`
}
