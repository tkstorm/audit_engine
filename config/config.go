package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tkstorm/audit_engine/msyql"
	"github.com/tkstorm/audit_engine/rabbit"
	"github.com/tkstorm/audit_engine/tool"
)

type EngineInfo struct {
	Name    string
	Version string
}

type CFG struct {
	cmd        CmdArgs
	EInfo      EngineInfo
	Test       bool
	ConfigFile string
	RabbitMq   rabbit.Config
	Mysql      mysql.Config
}

//version info
func (cfg *CFG) GetVersion(egi EngineInfo) string {
	return fmt.Sprintf("%s, %s", egi.Name, egi.Version)
}

//config init
func (cfg *CFG) Init(cmd CmdArgs) {
	//read config file
	viper.SetConfigFile(cmd.Cfg)
	if err := viper.ReadInConfig(); err != nil {
		tool.ErrorPanic(err, "viper read config error")
	}

	//test
	cfg.cmd = cmd
	cfg.Test = cmd.T
	cfg.ConfigFile = cmd.Cfg

	//version
	cfg.EInfo = EngineInfo{
		Name:    viper.GetString("name"),
		Version: viper.GetString("version"),
	}

	//init rabbitmq config
	cfg.RabbitMq = rabbit.Config{
		Host: viper.GetString("rabbitmq.host"),
		Port: viper.GetInt("rabbitmq.port"),
		User: viper.GetString("rabbitmq.user"),
		Pass: viper.GetString("rabbitmq.pass"),
	}

	//init mysql config
	cfg.Mysql = mysql.Config{
		Host: viper.GetString("mysql.host"),
		Port: viper.GetInt("mysql.port"),
		User: viper.GetString("mysql.user"),
		Pass: viper.GetString("mysql.pass"),
	}

}

func (cfg *CFG) PrintEnv() {
	//print config
	tool.PrettyPrint(cfg.GetVersion(cfg.EInfo))
	tool.PrettyPrint("config_file:", cfg.ConfigFile)
	tool.PrettyPrint("testing:", cfg.Test)

	//cmd print
	tool.PrettyPrint("cmd input", fmt.Sprintf("%+v", cfg.cmd))
}

func (cfg *CFG) PrintVersion() {
	tool.PrettyPrint(cfg.GetVersion(cfg.EInfo))
}

func (cfg *CFG) PrintHelpInfo() {
	pflag.PrintDefaults()
}
