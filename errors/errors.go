package errors

type BizError struct {
    msg  string
    Code int
}

const (
    NotFoundErrorCode    = 404
    InvalidFormErrorCode = 400
)

func (f *BizError) Error() string {
    return f.msg
}

func IsBizError(err error) bool {
    _, ok := err.(*BizError)
    return ok
}

func NotFound(msg string) error {
    return &BizError{
        Code: NotFoundErrorCode,
        msg:  msg,
    }
}

func IsNotFound(err error) bool {
    bizError, ok := err.(*BizError)
    if !ok {
        return false
    }
    return bizError.Code == NotFoundErrorCode
}

func InvalidForm(msg string) error {
    return &BizError{
        Code: InvalidFormErrorCode,
        msg:  msg,
    }
}

func IsInvalidForm(err error) bool {
    bizError, ok := err.(*BizError)
    if !ok {
        return false
    }
    return bizError.Code == InvalidFormErrorCode
}