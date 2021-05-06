package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/yiqiang3344/go-lib/helper"
	robotsrv "github.com/yiqiang3344/go-lib/proto/robot-srv"
	"log"
	"robot-srv/handler"
)

func init() {
	helper.InitCfg()
	helper.InitLogger()
}

func main() {
	project := helper.GetCfgString("project")
	//配置微服务链路追踪
	jaegerTracer, closer, err := helper.InitJaegerTracer("go.micro.service."+project, helper.GetCfgString("jaeger.address"))
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	reg := etcdv3.NewRegistry(
		registry.Addrs(helper.GetCfgString("etcd.address")),
	)

	service := micro.NewService(
		micro.Name("go.micro.service."+project),
		micro.Version("latest"),
		micro.Registry(reg),
		micro.WrapHandler(opentracing.NewHandlerWrapper(jaegerTracer)),
		micro.WrapHandler(helper.LogWrapper),
	)

	service.Init()

	_ = robotsrv.RegisterRobotSrvHandler(service.Server(), new(handler.RobotSrv))

	go helper.CheckSignReload(project, service)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
