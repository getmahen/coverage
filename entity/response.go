package entity

// Response is the container for what should be returned from all public facing endpoints
type Response struct {
	Result interface{} `json:"Result,omitempty"`
	Errors []Error     `json:"Errors,omitempty"`
}

// Error is used to report errors to the caller
type Error struct {
	Message string `json:"message"`
	Path    string `json:"path,omitempty"`
}

type CoverageCheckResponse struct {
	IsCovered bool
}

type CsaResponse struct {
	CsaFound bool
	Csa      string
}
