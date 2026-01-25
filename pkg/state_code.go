package pkg

import (
	"github.com/capyflow/housekeeper/locale"
	vhttp "github.com/capyflow/vortexv3/server/http"
)

var StateCodes = struct {
	PasswordNotMatch vhttp.SubCode // 密码不匹配
	UsernameNotMatch vhttp.SubCode // 用户名不匹配
}{
	PasswordNotMatch: vhttp.SubCode{Code: 10401, MsgKey: locale.K.PASSWORD_NOT_MATCH},
	UsernameNotMatch: vhttp.SubCode{Code: 10402, MsgKey: locale.K.USERNAME_NOT_MATCH},
}
