package errors

type AuthError string

func (e AuthError) Error() string {
	return string(e)
}

const (
	ErrInvalidToken  = AuthError("invalid_authorization_token")
	ErrServiceDown   = AuthError("service_unavailable")
	ErrInvalidHeader = AuthError("invalid_authorization_header")
	ErrMissingHeader = AuthError("missing_authorization_header")
	ErrUnauthorised  = AuthError("missing_permission")
)
