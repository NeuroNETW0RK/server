package configs

import (
	_ "embed"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"time"
)

var config = new(Config)

const (
	LocalConfigFile = "neuronet-server"
)

type Config struct {
	Mysql    Mysql    `mapstructure:"mysql"`
	Redis    Redis    `mapstructure:"redis"`
	Log      Log      `mapstructure:"log"`
	Server   Server   `mapstructure:"server"`
	FilePath FilePath `mapstructure:"filePath"`
}

type Mysql struct {
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Addr            string        `mapstructure:"addr"`
	Port            int           `mapstructure:"port"`
	Database        string        `mapstructure:"database"`
	MaxOpenConn     int           `mapstructure:"maxOpenConn"`
	MaxIdleConn     int           `mapstructure:"maxIdleConn"`
	ConnMaxLifeTime time.Duration `mapstructure:"connMaxLifeTime"`
	LogLevel        int           `mapstructure:"logLevel"`
}

type Redis struct {
	Addr       string `mapstructure:"addr"`
	Port       int    `mapstructure:"port"`
	Db         int    `mapstructure:"db"`
	MaxReTries int    `mapstructure:"maxRetries"`
	Password   string `mapstructure:"password"`
	PoolSize   int    `mapstructure:"poolSize"`
}

type Log struct {
	DisableCaller     bool     `mapstructure:"disableCaller"`
	DisableStacktrace bool     `mapstructure:"disableStacktrace"`
	Level             string   `mapstructure:"level"`
	Format            string   `mapstructure:"format"`
	OutputPaths       []string `mapstructure:"outputPaths"`
}

type Server struct {
	Address   string `mapstructure:"address"`
	Port      int    `mapstructure:"port"`
	Mode      string `mapstructure:"mode"`
	MachineIp string `mapstructure:"machineIp"`
}

type FilePath struct {
	AgentRootPath  string `mapstructure:"agentRootPath"`
	MountPath      string `mapstructure:"mountPath"`
	MountSubPath   string `mapstructure:"mountSubPath"`
	KubeConfigPath string `mapstructure:"kubeConfigPath"`
}

func InitConfig(configName string) {

	viper.SetConfigType("yaml")
	viper.SetConfigName(configName)
	viper.AddConfigPath("configs")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.Set("machineIp", os.Getenv("machineIp"))

	if err := viper.Unmarshal(config); err != nil {
		panic(err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})
}

func Get() Config {
	return *config
}
