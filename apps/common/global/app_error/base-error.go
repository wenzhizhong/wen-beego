package app_error

import "WenBeego/apps/common/helper"

type BaseError struct {
	Err   error
	Code  int
	Trace string
}

func (b *BaseError) Error() string {
	if b == nil {
		return ""
	} else if b.Err != nil {
		return b.Err.Error()
	}
	return ""
}

func baseError(err error, code int) error {
	if err == nil {
		return nil
	}
	return &BaseError{
		Err:   err,
		Code:  code,
		Trace: helper.GetTraceStr(),
	}
}
