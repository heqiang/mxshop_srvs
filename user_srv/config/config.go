package config

type MysqlConfig struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     int    `json:"port" mapstructure:"port"`
	Name     string `json:"name" mapstructure:"name"`
	User     string `json:"user" mapstructure:"user"`
	Password string `json:"password" mapstructure:"password"`
	Db       string `json:"db" mapstructure:"db"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type ServerConfig struct {
	Name          string `mapstructure:"name"`
	*MysqlConfig  `mapstructure:"mysql"`
	*LogConfig    `mapstructure:"log"`
	*ConsulConfig `mapstructure:"consul"`
}
