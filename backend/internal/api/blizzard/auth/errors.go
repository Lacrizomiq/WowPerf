package auth

// ResponseCode represents the response codes for the Battle.net API
type ResponseCode string

const (
	// Success codes
	CodeLinkSuccessful   ResponseCode = "link_successful"
	CodeUnlinkSuccessful ResponseCode = "unlink_successful"

	// Error codes
	CodeInvalidState      ResponseCode = "invalid_state"
	CodeMissingCode       ResponseCode = "missing_code"
	CodeTokenExchangeFail ResponseCode = "token_exchange_failed"
	CodeLinkFailed        ResponseCode = "link_failed"
	CodeUnlinkFailed      ResponseCode = "unlink_failed"
	CodeStatusCheckFailed ResponseCode = "status_check_failed"
	CodeAccountNotLinked  ResponseCode = "account_not_linked"
	CodeInsufficientScope ResponseCode = "insufficient_scope"
)

// ErrorResponse structure for error responses
type ErrorResponse struct {
	Error   string       `json:"error"`
	Code    ResponseCode `json:"code"`
	Details string       `json:"details,omitempty"`
}
