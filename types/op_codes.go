package opcodes

/*
Required for initial gateway connection
- OP_10, OP_1, OP_11, OP_2
*/
type OP_10_Hello struct {
	OP *int             `json:"op"`
	D  OP_10_Hello_Data `json:"d"`
	S  *int             `json:"s"`
	T  *string          `json:"t"`
}

type OP_10_Hello_Data struct {
	Heartbeat_Interval int `json:"heartbeat_interval"`
}

type OP_1_Heartbeat struct {
	OP int                 `json:"op"`
	D  OP_1_Heartbeat_Data `json:"d"`
}

type OP_1_Heartbeat_Data struct {
	Sequence *int `json:"seq"`
}

type OP_11_Heartbeat_ACK struct {
	OP int `json:"op"`
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

type Ready struct {
	V         int32             `json:"v"`
	User      User              `json:"user"`
	Guilds    Unavailable_Guild `json:"guilds"`
	SessionId string            `json:"session_id"`
}

type User struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Fiscriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
}

type Unavailable_Guild struct {
	Id          string `json:"id"`
	Unavailable string `json:"unavailable"`
}

type Application struct {
	Id    string `json:"id"`
	Flags string `json:"flags"`
}
