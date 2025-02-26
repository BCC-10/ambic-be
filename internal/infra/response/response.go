package response

type Err struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Res struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Payload    interface{} `json:"payload"`
}
