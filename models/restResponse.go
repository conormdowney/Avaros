package models

// The response object to set and return
// The routing package I used handles the rest api object.
// This is just deomnstrative of how it would look
type RestResponse struct {
	Result     interface{} // continas the object being returned
	Error      error       // contains any error
	StatusCode int         // contains the status code being returned
}
