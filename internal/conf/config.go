package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/capyflow/allspark-go/ds"
)

type Config struct {
	Port      int64        `toml:"Port" json:"port"`            // 端口
	RdbConfig *ds.DsConfig `toml:"RdbConfig" json:"rdb_config"` // 数据库配置
}

// 加载配置文件
func LoadConfig(confPath string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(confPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
