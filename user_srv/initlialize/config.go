package initlialize

import (
	"fmt"

	"github.com/spf13/viper"
	"mxshop_srvs/user_srv/global"
)

func GerEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GerEnvInfo("MXSHOP_DEBUG")
	ConfigFilePrefix := "config"
	configFileName := fmt.Sprintf("%s_pro.yaml", ConfigFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("%s_debug.yaml", ConfigFilePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
}
