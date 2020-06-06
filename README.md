#前言
简版的logstash，logstash由于是java写的，频繁的GC导致并发吞吐量和性能下降，所以想搞一个简单的日志收集器，熟悉下golang





# 打包命令

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build go-kafka.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build go-kafka.go
