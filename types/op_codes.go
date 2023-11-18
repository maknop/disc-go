package opcodes

/*
Required for initial gateway connection
- OP_10, OP_1, OP_11, OP_2
*/
type OP10Hello struct {
	OP *int          `json:"op"`
	D  OP10HelloData `json:"d"`
	S  *int          `json:"s"`
	T  *string       `json:"t"`
}

type OP10HelloData struct {
	HeartbeatInterval int `json:"heartbeatInterval"`
}

type OP1Heartbeat struct {
	OP int              `json:"op"`
	D  OP1HeartbeatData `json:"d"`
}

type OP1HeartbeatData struct {
	Sequence *int `json:"seq"`
}

type OP11HeartbeatACK struct {
	OP int `json:"op"`
}

type OP2Identity struct {
	OP int             `json:"op"`
	D  OP2IdentityData `json:"d"`
}

type OP2IdentityData struct {
	Token      string                `json:"token"`
	Intents    int                   `json:"intents"`
	Properties OP2IdentityProperties `json:"properties"`
}

type OP2IdentityProperties struct {
	OS      string `json:"$os"`
	Browser string `json:"$browser"`
	Device  string `json:"$device"`
}

type OP0Ready struct {
	T  string       `json:"t"`
	S  int          `json:"s"`
	OP int          `json:"op"`
	D  OP0ReadyData `json:"d"`
}

type OP0ReadyData struct {
	V                int32          `json:"v"`
	User             OP0User        `json:"user"`
	Guilds           []string       `json:"guilds"`
	SessionID        string         `json:"session_id"`
	ResumeGatewayURL string         `json:"resume_gateway_url"`
	Application      OP0Application `json:"application"`
}

type OP0User struct {
	Verified      bool   `json:"verified"`
	Username      string `json:"username"`
	MFAEnabled    bool   `json:"mfa_enabled"`
	ID            string `json:"id"`
	Flags         int32  `json:"flags"`
	Email         *int32 `json:"email"`
	Discriminator string `json:"discriminator"`
	Bot           bool   `json:"bot"`
	Avatar        string `json:"avatar"`
}

type OP0Application struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}
