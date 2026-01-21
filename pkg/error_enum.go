package pkg

import "errors"

var ErrorEnum = struct {
	ErrUsernameInvalid error
	ErrPasswordInvalid error
	ErrTokenInvalid    error
	ErrTokenExpired    error
}{
	ErrUsernameInvalid: errors.New("用户名无效"),
	ErrPasswordInvalid: errors.New("密码无效"),
	ErrTokenInvalid:    errors.New("Token无效"),
	ErrTokenExpired:    errors.New("Token已过期"),
}
