package main

import (
	"flag"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/handle"
	"mxshop_srvs/user_srv/initlialize"
	"mxshop_srvs/user_srv/proto"
	"mxshop_srvs/user_srv/utils"
	"net"
)

func main() {
	IP := flag.String("ip", "192.168.31.101", "ip地址")
	Port := flag.Int("port", 0, "端口")

	flag.Parse()
	initlialize.InitConfig()
	initlialize.Initdb()
	initlialize.Initlogger(global.ServerConfig.LogConfig)
	zap.L().Info("")
	server, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic(err)
	}
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.L().Info(fmt.Sprintf("服务监听地址为：%s:%d", *IP, *Port))
	g := grpc.NewServer()
	proto.RegisterUserServer(g, &handle.UserService{})
	reflection.Register(g)
	//注册grpc的服务健康检查
	grpc_health_v1.RegisterHealthServer(g, health.NewServer())
	//grpc服务注册
	cfg := consulapi.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port)

	client, err := consulapi.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	check := &consulapi.AgentServiceCheck{
		GRPC:                           "https://www.baidu.com",
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
		Notes:                          "健康检查",
	}
	registeration := &consulapi.AgentServiceRegistration{
		Name:    global.ServerConfig.Name,
		ID:      global.ServerConfig.Name,
		Port:    *Port,
		Tags:    []string{"hq", "srv"},
		Address: global.ServerConfig.ConsulConfig.Host,
		Check:   check,
	}
	err = client.Agent().ServiceRegister(registeration)
	if err != nil {
		panic(err)
	}
	zap.L().Info(fmt.Sprintf("注册的grpc服务地址为=>%s:%d", global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port))
	err = g.Serve(server)
	if err != nil {
		return
	}

}
