package config

type DbConfig struct {
	File                   string `yaml:"file"`
	MaxIdleConns           int    `yaml:"MaxIdleConns"`
	MaxOpenConns           int    `yaml:"MaxOpenConns"`
	ConnMaxIdleTimeSeconds int    `yaml:"ConnMaxIdleTimeSeconds"`
	ConnMaxLifetimeSeconds int    `yaml:"ConnMaxLifetimeSeconds"`
}

type NeovideoConfig struct {
	Port     int    `yaml:"port"`
	LogLevel string `yaml:"log_level"`
	Db       DbConfig
}
