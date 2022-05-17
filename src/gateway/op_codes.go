package gateway

/*
opcode 10 hello
- Server reponse after client sends initial request to the server
*/
type OP_10_Hello struct {
	OP *int    `json:"op"`
	D  Data    `json:"d"`
	S  *int    `json:"s"`
	T  *string `json:"t"`
}

type Data struct {
	Heartbeat_Interval int  `json:"heartbeat_interval"`
	Sequence           *int `json:"seq"`
}

type OP_1_Heartbeat struct {
	OP int  `json:"op"`
	D  *int `json:"d"`
}

type OP_11_Heartbeat_ACK struct {
	OP int `json:"op"`
}
