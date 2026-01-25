package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/capyflow/allspark-go/ds"
	vpkg "github.com/capyflow/vortexv3/pkg"
)

type Config struct {
	Port      int64           `toml:"Port" json:"port"`            // 端口
	RdbConfig *ds.DsConfig    `toml:"RdbConfig" json:"rdb_config"` // 数据库配置
	Jwt       *vpkg.JwtOption `toml:"Jwt" json:"jwt"`              // JWT配置
	Admin     *Admin          `toml:"Admin" json:"admin"`          // 管理员
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
