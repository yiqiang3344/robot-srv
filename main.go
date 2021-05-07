package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/uber/jaeger-client-go"
	robotsrv "github.com/yiqiang3344/go-lib/proto/robot-srv"
	"github.com/yiqiang3344/go-lib/utils/build"
	"github.com/yiqiang3344/go-lib/utils/config"
	cLog "github.com/yiqiang3344/go-lib/utils/log"
	"github.com/yiqiang3344/go-lib/utils/trace"
	"log"
	"robot-srv/handler"
)

func init() {
	config.InitCfg()
	cLog.InitLogger(config.GetCfgString("project"), config.GetCfgBool("showLogToConsole", false))
}

func main() {
	project := config.GetCfgString("project")
	//配置微服务链路追踪
	jaegerTracer, closer, err := trace.InitJaegerTracer(
		"go.micro.service."+project,
		config.GetCfgString("jaeger.address"),
		jaeger.SamplerTypeConst,
		1,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	reg := etcdv3.NewRegistry(
		registry.Addrs(config.GetCfgString("etcd.address")),
	)

	service := micro.NewService(
		micro.Name("go.micro.service."+project),
		micro.Version("latest"),
		micro.Registry(reg),
		micro.WrapHandler(opentracing.NewHandlerWrapper(jaegerTracer)),
		micro.WrapHandler(cLog.LogWrapper),
	)

	service.Init()

	_ = robotsrv.RegisterRobotSrvHandler(service.Server(), new(handler.RobotSrv))

	go build.CheckSignReload(project, service)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
