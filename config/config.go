package config

// 废弃
// AppConfig 用于映射 YAML 配置
type AppConfig struct {
	Settings Settings `yaml:"settings"`
}

// Settings 包含应用程序的配置
type Settings struct {
	Application ApplicationConfig `yaml:"application"`
}

// ApplicationConfig 定义应用程序的相关配置
type ApplicationConfig struct {
	Mode string `yaml:"mode"`
	Host string `yaml:"host"`
}
