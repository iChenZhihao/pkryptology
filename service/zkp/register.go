package zkp

import (
	"fmt"
	"github.com/coinbase/kryptology/service/global"
	"github.com/go-zookeeper/zk"
	"github.com/golang/glog"
	"net"
)

const ServerRegisterPath = "/gg20/nodes"

func Register() {
	ip, err := GetLocalIP()
	if err != nil {
		glog.Error(err.Error())
	}
	address := fmt.Sprintf("%s:%s", ip, global.Config.Server.Port)
	err = registerService(global.GetZkManager().GetConn(), ServerRegisterPath, address)
	if err != nil {
		glog.Error("Sign Node Register To Zookeeper Error: ", err.Error())
	} else {
		glog.Info(address, " Register To Zk Success! ")
	}
}

// 注册服务到 Zookeeper
func registerService(conn *zk.Conn, basePath, address string) error {
	// 确保根节点存在
	if exists, _, err := conn.Exists(basePath); err != nil {
		glog.Error(err.Error())
		return err
	} else if !exists {
		if _, err := conn.Create(basePath, nil, 0, zk.WorldACL(zk.PermAll)); err != nil {
			return err
		}
	}
	// 创建临时节点
	nodePath := fmt.Sprintf("%s/%s", basePath, address)
	_, err := conn.Create(nodePath, []byte(address), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	return nil
}

// GetLocalIP 获取本机 IP 地址
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no suitable IP address found")
}
