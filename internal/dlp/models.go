package dlp

type CheckRequest struct {
	TestFile      string
	TestURL       string
	HTTPMethod    string
	FileExtension string // file extension to add to file_name
}

type CheckResponse struct {
	StatusCode int
	StatusText string
}

type Result struct {
	IsDLPActive bool
	StatusText  string
}
