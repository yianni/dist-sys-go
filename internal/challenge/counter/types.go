package counter

type addRequest struct {
	Type  string `json:"type"`
	Delta int    `json:"delta"`
}

type addResponse struct {
	Type string `json:"type"`
}

type readRequest struct {
	Type string `json:"type"`
}

type readResponse struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}
