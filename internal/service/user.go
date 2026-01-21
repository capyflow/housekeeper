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
	ctx             context.Context
	jwtOption       *conf.JwtOptions
	admin           *conf.Admin
	SkipExpireCheck bool
}

func NewUserService(ctx context.Context, config *conf.Config) *UserService {
	return &UserService{
		ctx:       ctx,
		jwtOption: config.Jwt,
		admin:     config.Admin,
	}
}

func (us *UserService) Login(ctx context.Context, req *model.LoginReq) (*model.LoginResp, error) {
	// 验证用户名和密码
	if req.Username != us.admin.Username {
		return nil, pkg.ErrorEnum.ErrUsernameInvalid
	}

	if req.Password != us.admin.Password {
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
	expirationTime := time.Now().Add(time.Duration(us.jwtOption.Expire) * time.Hour)

	// 创建JWT claims
	claims := jwt.MapClaims{
		"username": username,
		"exp":      expirationTime.Unix(), // 过期时间
		"iat":      time.Now().Unix(),     // 签发时间
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名token
	tokenString, err := token.SignedString([]byte(us.jwtOption.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 解析token的方法
// ignoreExpired: true 表示忽略过期检查，false 表示严格检查过期时间
func (us *UserService) ParseToken(tokenString string, ignoreExpired bool) (*jwt.MapClaims, error) {
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, pkg.ErrorEnum.ErrTokenInvalid
		}
		return []byte(us.jwtOption.Secret), nil
	}, jwt.WithoutClaimsValidation()) // 先不验证claims，手动控制

	if err != nil {
		return nil, pkg.ErrorEnum.ErrTokenInvalid
	}

	// 获取claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, pkg.ErrorEnum.ErrTokenInvalid
	}

	// 如果不忽略过期，则检查exp字段
	if !ignoreExpired {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, pkg.ErrorEnum.ErrTokenExpired
			}
		}
	}

	return &claims, nil
}
