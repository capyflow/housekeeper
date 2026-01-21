package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/capyflow/allspark-go/ds"
)

type Config struct {
	Port      int64        `toml:"Port" json:"port"`            // 端口
	RdbConfig *ds.DsConfig `toml:"RdbConfig" json:"rdb_config"` // 数据库配置
	Jwt       *JwtOptions  `toml:"Jwt" json:"jwt"`              // JWT配置
	Admin     *Admin       `toml:"Admin" json:"admin"`          // 管理员
}

type JwtOptions struct {
	Secret          string `toml:"Secret" json:"secret"`                     // JWT密钥
	Expire          int64  `toml:"Expire" json:"expire"`                     // JWT过期时间(小时)
	SkipExpireCheck bool   `toml:"SkipExpireCheck" json:"skip_expire_check"` // 是否跳过过期时间检查
}

type Admin struct {
	Username string `toml:"Username" json:"username"`
	Password string `toml:"Password" json:"password"`
}

// 加载配置文件
func LoadConfig(confPath string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(confPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
