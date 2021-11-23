package main

import (
	"flag"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
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
	"os"
	"os/signal"
	"syscall"
)

func main() {
	initlialize.InitConfig()
	initlialize.Initdb()
	//initlialize.Initlogger(global.ServerConfig.LogConfig)
	initlialize.Initlogger()

	IP := flag.String("ip", global.ServerConfig.ConsulConfig.ServerHost, "ip地址")
	Port := flag.Int("port", 0, "端口")
	flag.Parse()
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	server, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic(err)
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

	registeration := &consulapi.AgentServiceRegistration{
		Name: global.ServerConfig.Name,
		//ID:      global.ServerConfig.Name,
		ID:      fmt.Sprintf("%s", uuid.NewV4()),
		Port:    *Port,
		Tags:    []string{"hq", "srv"},
		Address: global.ServerConfig.ConsulConfig.ServerHost,
		Check: &consulapi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.ConsulConfig.ServerHost, *Port),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
			Notes:                          "健康检查",
		},
	}
	err = client.Agent().ServiceRegister(registeration)
	if err != nil {
		panic(err)
	}

	go func() {
		err = g.Serve(server)
		if err != nil {
			return
		}
	}()
	//终止信息的接受
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := client.Agent().ServiceDeregister(registeration.ID); err != nil {
		zap.L().Info("服务注销失败")
	}
	zap.L().Info("注销成功")
}
