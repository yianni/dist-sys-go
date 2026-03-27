package broadcast

type broadcastRequest struct {
	Type    string `json:"type"`
	Message int    `json:"message"`
}

type broadcastResponse struct {
	Type string `json:"type"`
}

type readRequest struct {
	Type string `json:"type"`
}

type readResponse struct {
	Type     string `json:"type"`
	Messages []int  `json:"messages"`
}

type topologyRequest struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

type topologyResponse struct {
	Type string `json:"type"`
}

type gossipRequest struct {
	Type     string `json:"type"`
	Messages []int  `json:"messages"`
}
