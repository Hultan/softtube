package softtube

import "fmt"

type ErrYtdlp struct {
	Message string
	Reason  string
	Err     error
}

func (e ErrYtdlp) Error() string {
	return fmt.Sprintf("%s. Reason:\n%s\n\n%s\n", e.Message, e.Reason, e.Err)
}

func newErrYtdlp(message string, reason string, err error) error {
	return ErrYtdlp{
		Message: message,
		Reason:  reason,
		Err:     err,
	}
}
