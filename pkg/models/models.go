package models

type Configuration struct {
	Configs               []Config `json:"url"`
	QuantityOfBadRequests int      `json:"quantity_of_bad_requests"`
}

type Config struct {
	Url          string   `json:"url"`
	Checks       []string `json:"checks"`
	MinChecksCnt int      `json:"min_checks_cnt"`
}

type Response struct {
	Url        string   `json:"url,omitempty"`
	StatusCode int      `json:"statusCode,omitempty"`
	Text       string   `json:"text,omitempty"`
	Checks     []string `json:"checks,omitempty"`
}

type Params struct {
	Url         string `json:"url"`
	HeaderKey   string `json:"headerKey"`
	HeaderValue string `json:"headerValue"`
}
