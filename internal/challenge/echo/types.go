package echo

type echoRequest struct {
	Type  string `json:"type"`
	MsgID int    `json:"msg_id,omitempty"`
	Echo  string `json:"echo"`
}

type echoResponse struct {
	Type string `json:"type"`
	Echo string `json:"echo"`
}
