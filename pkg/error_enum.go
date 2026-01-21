package pkg

import "errors"

var ErrorEnum = struct {
	ErrUsernameInvalid error
	ErrPasswordInvalid error
}{
	ErrUsernameInvalid: errors.New("用户名无效"),
	ErrPasswordInvalid: errors.New("密码无效"),
}
