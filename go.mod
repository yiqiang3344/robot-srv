module robot-srv

go 1.13

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

//本地调试时可通过此方法引入本地代码
//replace github.com/yiqiang3344/go-lib => /Users/xinfei/docker/code/go/src/github.com/yiqiang3344/go-lib

require (
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/garyburd/redigo v1.6.2
	github.com/google/uuid v1.1.2 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jmoiron/sqlx v1.3.3
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/registry/etcdv3/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/trace/opentracing/v2 v2.9.1
	github.com/uber/jaeger-client-go v2.28.0+incompatible
	github.com/yiqiang3344/go-lib v0.0.7
	go.etcd.io/bbolt v1.3.5 // indirect
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
