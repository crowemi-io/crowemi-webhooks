package main

type CrowemiConfig struct {
	ClientName      string            `json:"client_name"`
	ClientId        string            `json:"client_id"`
	ClientSecretKey string            `json:"client_secret_key"`
	Uri             map[string]string `json:"uri"`
	GcpProjectId    string            `json:"gcp_project_id"`
	GcpLogTopic     string            `json:"gcp_log_topic"`
	Env             string            `json:"env"`
	Debug           bool              `json:"debug"`
}

type BotConfig struct {
	Name         string `json:"name"`
	Token        string `json:"token"`
	AllowedUsers []int  `json:"allowed_users"`
}

type Config struct {
	App       string        `json:"app"`
	Crowemi   CrowemiConfig `json:"crowemi"`
	BotConfig []BotConfig   `json:"bot_config"`
}
