package config

type crowemi struct {
	ClientName      string            `json:"client_name"`
	ClientId        string            `json:"client_id"`
	ClientSecretKey string            `json:"client_secret_key"`
	Uri             map[string]string `json:"uri"`
	Env             string            `json:"env"`
	Debug           bool              `json:"debug"`
}

type googleCloudCredential struct{}

type googleCloud struct {
	ProjectId  string                `json:"project_id"`
	Topic      string                `json:"topic"`
	Credential googleCloudCredential `json:"credentials"`
}

type bot struct {
	ChannelId    string `json:"channel_id"`
	Token        string `json:"token"`
	AllowedUsers []int  `json:"allowed_users"`
}

type Webhooks struct {
	App         string         `json:"app"`
	Crowemi     crowemi        `json:"crowemi"`
	BotConfig   map[string]bot `json:"bot"`
	GoogleCloud googleCloud    `json:"google_cloud"`
}
