package errors

type BizError struct {
    msg  string
    Code int
}

const (
    NotFoundErrorCode    = 404
    InvalidFormErrorCode = 400
    ForbiddenErrorCode   = 403
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

func Forbidden(msg string) error {
    return &BizError{
        Code: ForbiddenErrorCode,
        msg:  msg,
    }
}

func IsForbidden(err error) bool {
    bizError, ok := err.(*BizError)
    if !ok {
        return false
    }
    return bizError.Code == ForbiddenErrorCode
}