package request

type Response struct {
	Resource string `json:"resource,omitempty"`
	Data     string `json:"data,omitempty"`
}

type result struct {
	Response *Response
	Err      error
}
