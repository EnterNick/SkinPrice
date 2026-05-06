package errx

type Code string

const (
	CodeUnknown         Code = "unknown"
	CodeInvalidArgument Code = "invalid_argument"
	CodeNotFound        Code = "not_found"
	CodeAlreadyExists   Code = "already_exists"
	CodeConflict        Code = "conflict"
	CodeUnauthorized    Code = "unauthorized"
	CodeForbidden       Code = "forbidden"
	CodeUnavailable     Code = "unavailable"
	CodeTimeout         Code = "timeout"
	CodeExternal        Code = "external"
	CodeInternal        Code = "internal"
	CodeEmailInvalid    Code = "invalid_email"
	WeakPassword        Code = "weak_password"
)
