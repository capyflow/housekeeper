package locale

import (
	"fmt"
	"strings"
)

type I18nStruct struct {
	Value    string
	mapValue map[string]string
}

func (i *I18nStruct) GetMap() map[string]string {
	return i.mapValue
}

func (i *I18nStruct) GetValue(key string) string {
	return i.mapValue[key]
}

func (I18n *I18nStruct) GetLangValue(i18nKey string, lang string, defaultValue string) string {
	key := strings.ToLower(fmt.Sprintf("%s.%s", i18nKey, lang))
	if value, ok := I18n.mapValue[key]; ok {
		return value
	}
	return defaultValue
}

var I18nValue = I18nStruct{Value: "{\"PASSWORD_NOT_MATCH.en-us\":\"Password not match\",\"PASSWORD_NOT_MATCH.zh-cn\":\"密码不匹配\",\"USERNAME_NOT_MATCH.en-us\":\"Username not match\",\"USERNAME_NOT_MATCH.zh-cn\":\"用户名不匹配\"}", mapValue: map[string]string{"PASSWORD_NOT_MATCH.zh-cn": "密码不匹配", "USERNAME_NOT_MATCH.zh-cn": "用户名不匹配", "PASSWORD_NOT_MATCH.en-us": "Password not match", "USERNAME_NOT_MATCH.en-us": "Username not match"}}

var K = struct {
	PASSWORD_NOT_MATCH string
	USERNAME_NOT_MATCH string
}{
	USERNAME_NOT_MATCH: "USERNAME_NOT_MATCH",
	PASSWORD_NOT_MATCH: "PASSWORD_NOT_MATCH",
}
