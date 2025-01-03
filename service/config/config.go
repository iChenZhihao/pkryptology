package config

import "time"

// Config 组合全部配置模型
type Config struct {
	Server    Server    `mapstructure:"server"`
	Cron      Cron      `mapstructure:"cron"`
	Zookeeper Zookeeper `mapstructure:"zookeeper"`
	GG20      GG20      `mapstructure:"gg20" json:"gg20"`
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

type GG20 struct {
	// 稳定窗口时间 单位: 毫秒
	StableTimeWindow time.Duration `mapstructure:"stableTimeWindow" json:"stableTimeWindow"`
}

func (g *GG20) GetStableTimeWindow() time.Duration {
	if g.StableTimeWindow > 0 {
		return g.StableTimeWindow * time.Millisecond
	}
	return time.Duration(5000) * time.Millisecond
}
