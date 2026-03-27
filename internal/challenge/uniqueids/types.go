package uniqueids

type generateRequest struct {
	Type string `json:"type"`
}

type generateResponse struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}
