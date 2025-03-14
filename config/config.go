package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	Server     ServerConfig          `mapstructure:"server"`
	Database   DatabaseConfig        `mapstructure:"database"`
	Redis      ServerInfo            `mapstructure:"redis"`
	RabbitMq   ServerInfo            `mapstructure:"rabbit-mq" json:"rabbitMq"`
	Mongo      ServerInfo            `mapstructure:"mongo" json:"mongo"`
	Log        LogConfig             `mapstructure:"log"`
	RateLimit  RateLimitConfig       `mapstructure:"rate-limit" json:"rateLimit"`
	ThirdParty ThirdPartyConfig      `mapstructure:"third_party"`
	RpcServer  map[string]ServerInfo `mapstructure:"rpc-server" json:"rpcServer"`
}

type ServerConfig struct {
	LiveApi APIConfig `mapstructure:"liveapi"`
	BSApi   APIConfig `mapstructure:"bsapi"`
	H5Api   APIConfig `mapstructure:"h5api"`
	Payment APIConfig `mapstructure:"payment"`
}

type APIConfig struct {
	Port      int    `mapstructure:"port"`
	Mode      string `mapstructure:"mode"`
	UrlPrefix string `mapstructure:"url-prefix" json:"urlPrefix"`
	AesKey    string `mapstructure:"aes-key" json:"aesKey"`
}

type DatabaseConfig struct {
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	Database string `mapstructure:"database" json:"database"`
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	LogMode  bool   `mapstructure:"log-mode" json:"logMode"`
}

type ServerInfo struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	User string `mapstructure:"user" json:"user"`
	Pwd  string `mapstructure:"pwd" json:"pwd"`
	Db   string `mapstructure:"db" json:"db"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
	File   struct {
		Path       string `mapstructure:"path"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxAge     int    `mapstructure:"max_age"`
		MaxBackups int    `mapstructure:"max_backups"`
	} `mapstructure:"file"`
}

type RateLimitConfig struct {
	FillInterval int64 `mapstructure:"fill-interval" json:"fillInterval"`
	Capacity     int64 `mapstructure:"capacity" json:"capacity"`
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	configPath := fmt.Sprintf("configs/%s/config.yaml", env)
	return loadConfigFromFile(configPath)
}

func loadConfigFromFile(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 替换环境变量
	for _, key := range viper.AllKeys() {
		val := viper.GetString(key)
		if strings.HasPrefix(val, "${") && strings.HasSuffix(val, "}") {
			envKey := strings.TrimSuffix(strings.TrimPrefix(val, "${"), "}")
			if envVal := os.Getenv(envKey); envVal != "" {
				viper.Set(key, envVal)
			}
		}
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
