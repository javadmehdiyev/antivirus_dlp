package dlp

type CheckRequest struct {
	TestFile   string
	TestURL    string
	HTTPMethod string
}

type CheckResponse struct {
	StatusCode int
	StatusText string
}

type Result struct {
	IsDLPActive bool
	StatusText  string
}
