> 该项目致力于微服务日志消息消费中间件，聚合各个服务的日志，简版的logstash，logstash由于是java写的，频繁的GC导致并发吞吐量和性能下降，
所以想搞一个简单的日志收集器

# 打包命令

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build go-kafka.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build go-kafka.go
