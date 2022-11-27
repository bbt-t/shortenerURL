package app

type CreateShortURLRequest struct {
	URL string `json:"url"`
}

type Resp struct {
	Result string `json:"result"`
}
