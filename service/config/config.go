package config

// Config 组合全部配置模型
type Config struct {
	Server    Server    `mapstructure:"server"`
	Cron      Cron      `mapstructure:"cron"`
	Zookeeper Zookeeper `mapstructure:"zookeeper"`
}

// Server 服务启动端口号配置
type Server struct {
	Port string `mapstructure:"port"`
}

// Cron 定时任务配置
type Cron struct {
	Enable bool `mapstructure:"enable"`
}

// Zookeeper Zk连接配置
type Zookeeper struct {
	Servers        []string `mapstructure:"servers"`
	SessionTimeout int      `mapstructure:"sessionTimeout"`
}
