package service

import (
	"context"
	"time"

	"github.com/capyflow/housekeeper/internal/conf"
	"github.com/capyflow/housekeeper/internal/model"
	"github.com/capyflow/housekeeper/pkg"
	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	ctx       context.Context
	JwtOption *conf.JwtOptions
	Admin     *conf.Admin
}

func NewUserService(ctx context.Context, config *conf.Config) *UserService {
	return &UserService{
		ctx:       ctx,
		JwtOption: config.Jwt,
		Admin:     config.Admin,
	}
}

func (us *UserService) Login(ctx context.Context, req *model.LoginReq) (*model.LoginResp, error) {
	// 验证用户名和密码
	if req.Username != us.Admin.Username {
		return nil, pkg.ErrorEnum.ErrUsernameInvalid
	}

	if req.Password != us.Admin.Password {
		return nil, pkg.ErrorEnum.ErrPasswordInvalid
	}

	// 生成token
	token, err := us.generateToken(req.Username)
	if err != nil {
		return nil, err
	}

	return &model.LoginResp{
		Token: token,
	}, nil
}

// 生成jwt token
func (us *UserService) generateToken(username string) (string, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(time.Duration(us.JwtOption.Expire) * time.Hour)

	// 创建JWT claims
	claims := jwt.MapClaims{
		"username": username,
		"exp":      expirationTime.Unix(), // 过期时间
		"iat":      time.Now().Unix(),     // 签发时间
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名token
	tokenString, err := token.SignedString([]byte(us.JwtOption.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
