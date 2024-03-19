package googlebooks

import "net/http"

// Response base response for Razorpay response
type Response struct {
	gatewayError bool
	Error        error
	StatusCode   int
}

// SuccessResponse returns Success Response
func SuccessResponse(response *http.Response) *Response {
	return &Response{
		Error:        nil,
		StatusCode:   response.StatusCode,
		gatewayError: false,
	}
}

// ErrorResponse returns Error Response
func ErrorResponse(err error, gatewayError bool) *Response {
	return &Response{
		Error:        err,
		StatusCode:   http.StatusInternalServerError,
		gatewayError: gatewayError,
	}
}

// IsError returns true if there is an error in the response
func (r *Response) IsError() bool {
	if r.Error != nil {
		return true
	}

	if r.StatusCode >= 300 {
		return true
	}

	return false
}

// IsGatewayError returns true if error is from razorpay gateway
func (r *Response) IsGatewayError() bool {
	return r.gatewayError
}

// WebResponse extends *http.Response and adds few custom methods
type WebResponse struct {
	*http.Response
}

// NewWebResponse returns new instance of WebResponse
func NewWebResponse(response *http.Response) *WebResponse {
	return &WebResponse{Response: response}
}

// IsError returns true if StatusCode is greater than equal 300
func (w *WebResponse) IsError() bool {
	return w.StatusCode >= 300
}
