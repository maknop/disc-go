package gateway

/*
opcode 10 hello
- Server reponse after client sends initial request to the server
*/
type OP_10 struct {
	OP *int    `json:"op"`
	D  Data    `json:"d"`
	S  *int    `json:"s"`
	T  *string `json:"t"`
}

type Data struct {
	Heartbeat_Interval int  `json:"heartbeat_interval"`
	Sequence           *int `json:"seq"`
}

type OP_1 struct {
	OP int  `json:"op"`
	D  *int `json:"d"`
}

type OP_11 struct {
	OP int `json:"op"`
}
