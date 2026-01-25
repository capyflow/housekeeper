package service

import (
	"context"
	"time"

	"github.com/capyflow/housekeeper/internal/conf"
	"github.com/capyflow/housekeeper/internal/model"
	"github.com/capyflow/housekeeper/pkg"
	vpkg "github.com/capyflow/vortexv3/pkg"
	vhttp "github.com/capyflow/vortexv3/server/http"
	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	ctx             context.Context
	jwtOption       *vpkg.JwtOption
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

func (us *UserService) Login(ctx *vhttp.Context, req *model.LoginReq) (*model.LoginResp, error) {
	// 验证用户名和密码
	if req.Username != us.admin.Username {
		return nil, pkg.ErrorEnum.ErrUsernameInvalid
	}

	if req.Password != us.admin.Password {
		return nil, pkg.ErrorEnum.ErrPasswordInvalid
	}

	ip := ctx.GinContext().ClientIP()

	// 生成token
	token, err := us.generateToken(req.Username, ip)
	if err != nil {
		return nil, err
	}

	return &model.LoginResp{
		Token: token,
	}, nil
}

// 生成jwt token
func (us *UserService) generateToken(username, ip string) (string, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(time.Duration(us.jwtOption.Expire) * time.Hour)

	// 创建JWT claims
	claims := jwt.MapClaims{
		"uid": username,
		"exp": expirationTime.Unix(), // 过期时间
		"iat": time.Now().Unix(),     // 签发时间
		"ip":  ip,
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名token
	tokenString, err := token.SignedString([]byte(us.jwtOption.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
