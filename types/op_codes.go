package opcodes

type OP_0_Ready struct {
	T  string          `json:"t"`
	S  int             `json:"s"`
	OP int             `json:"op"`
	D  OP_0_Ready_Data `json:"d"`
}

type OP_0_Ready_Data struct {
	V                int32            `json:"v"`
	User             OP_0_User        `json:"user"`
	Guilds           []string         `json:"guilds"`
	SessionId        string           `json:"session_id"`
	ResumeGatewayUrl string           `json:"resume_gateway_url"`
	Application      OP_0_Application `json:"application"`
}

type OP_0_User struct {
	Verified      bool   `json:"verified"`
	Username      string `json:"username"`
	MFAEnabled    bool   `json:"mfa_enabled"`
	Id            string `json:"id"`
	Flags         int32  `json:"flags"`
	Email         *int32 `json:"email"`
	Discriminator string `json:"discriminator"`
	Bot           bool   `json:"bot"`
	Avatar        string `json:"avatar"`
}

type OP_0_Application struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

type OP_1_Heartbeat struct {
	OP int                 `json:"op"`
	D  OP_1_Heartbeat_Data `json:"d"`
}

type OP_1_Heartbeat_Data struct {
	Sequence *int `json:"seq"`
}

type OP_2_Identity struct {
	OP int                `json:"op"`
	D  OP_2_Identity_Data `json:"d"`
}

type OP_2_Identity_Data struct {
	Token      string                   `json:"token"`
	Intents    int                      `json:"intents"`
	Properties OP_2_Identity_Properties `json:"properties"`
}

type OP_2_Identity_Properties struct {
	OS      string `json:"$os"`
	Browser string `json:"$browser"`
	Device  string `json:"$device"`
}

type OP_6_Reconnect struct {
	OP string `json:"op"`
	D  string `json:"d"`
}

type OP_9_Invalid_Session struct {
	OP string `json:"op"`
	D  string `json:"d"`
}

type OP_10_Hello struct {
	OP *int             `json:"op"`
	D  OP_10_Hello_Data `json:"d"`
	S  *int             `json:"s"`
	T  *string          `json:"t"`
}

type OP_10_Hello_Data struct {
	Heartbeat_Interval int `json:"heartbeat_interval"`
}

type OP_11_Heartbeat_ACK struct {
	OP int `json:"op"`
}
