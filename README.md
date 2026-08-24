# route-origin-guardian

编号16，纯Go BGP/BMP路由采集与RPKI源验证服务。数据面提供两个受限协议采集端口：BGP默认`:1790`、BMP默认`:11019`。为了支持离线冒烟，采集器也接受每行一个JSON事件；正式协议适配器可通过领域接口替换，不影响RIB、RPKI和分发层。

## 快速启动

```sh
go run ./cmd/guardian
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

除健康、就绪和指标端点外，API需要`X-API-Key: dev-key`。生产环境请设置`ROG_API_KEY`并使用TLS终止层。

## 冒烟流程

导入ROA快照：

```sh
curl -X POST -H 'X-API-Key: dev-key' -H 'Content-Type: application/json' \
  --data @examples/roa-snapshot.json http://localhost:8080/v1/rpki/snapshots
```

发布一条路由：

```sh
curl -X POST -H 'X-API-Key: dev-key' -H 'Content-Type: application/json' \
  --data @examples/route-announce.json http://localhost:8080/v1/routes
curl -H 'X-API-Key: dev-key' 'http://localhost:8080/v1/rpki/validate?prefix=203.0.113.0/24&asn=64500'
curl -H 'X-API-Key: dev-key' 'http://localhost:8080/v1/rpki/validate?prefix=203.0.113.0/24&asn=64599'
curl -H 'X-API-Key: dev-key' 'http://localhost:8080/v1/rpki/validate?prefix=192.0.2.0/24&asn=64500'
curl -H 'X-API-Key: dev-key' 'http://localhost:8080/v1/events?after=0'
```

撤回路由时将`operation`改为`withdraw`。也可以向`:1790`或`:11019`发送JSON行，例如`printf '%s\n' "$(cat examples/route-announce.json)" | nc localhost 1790`。

## 目录与边界

`internal`按session、bgp、bmp、rib、rpki、validation、analysis、distribution、tenant和platform领域拆分；核心领域分别包含domain、application、adapter和infrastructure层。validation链路保留可持久化的路由源验证失败分类，协议解析使用固定消息上限，所有后台循环接受context取消，内存RIB用于离线模式，`migrations/001_init.sql`提供PostgreSQL事件和检查点表供生产仓储实现。

## 验证

```sh
gofmt -w $(rg --files -g '*.go')
go vet ./...
go test -race ./...
go build ./...
./scripts/check-lines.sh
```

源码不依赖外部数据库即可启动；数据库、gRPC服务契约、Docker部署文件和OpenAPI文档已随项目提供。服务收到SIGINT/SIGTERM后停止监听并释放采集端口。
