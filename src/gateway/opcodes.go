package gateway

/*
opcode 10 hello
- Server reponse after client sends initial request to the server
*/
type Payload struct {
	OP *int    `json:"op"`
	D  Data    `json:"d"`
	S  *int    `json:"s"`
	T  *string `json:"t"`
}

type Data struct {
	Heartbeat_Interval *int `json:"heartbeat_interval"`
	Sequence           *int `json:"seq"`
}
