package errx

import "encoding/json"

type APIError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	payload, err := json.Marshal(e)
	if err != nil {
		return `{"code":"internal","message":"internal error"}`
	}
	return string(payload)
}

func FromError(err error, fallbackMsg string) error {
	if err == nil {
		return nil
	}

	message := fallbackMsg
	if message == "" {
		message = err.Error()
	}

	return &APIError{
		Code:    CodeOf(err),
		Message: message,
	}
}
