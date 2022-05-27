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

type OP_0_Ready struct {
	T  string          `json:"t"`
	S  int32           `json:"s"`
	OP string          `json:"op"`
	D  OP_0_Ready_Data `json:"d"`
}

type OP_0_Ready_Data struct {
	V                    int32            `json:"v"`
	UserSettings         []string         `json:"user_settings"`
	User                 OP_0_User        `json:"user"`
	SessionType          string           `json:"session_type"`
	SessionId            string           `json:"session_id"`
	Relationships        []string         `json:"relationships"`
	PrivateChannels      []string         `json:"private_channels"`
	Presences            []string         `json:"presences"`
	Guilds               []string         `json:"guilds"`
	GuildJoinRequests    []string         `json:"guild_join_requests"`
	GeoOrderedRTCRegions []string         `json:"geo_ordered_rtc_regions"`
	Application          OP_0_Application `json:"application"`
	Trace                []string         `json:"_trace"`
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
	Id    string `json:"id"`
	Flags string `json:"flags"`
}
