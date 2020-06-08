package job

import (
	"github.com/robfig/cron"
	"go-kafka/gzip"
	"go-kafka/setting"
	"log"
	"os"
	"strings"
)

type TestJob struct {
}

type GzipJob struct {
}

func (this GzipJob) Run() {
	kafka := setting.App.Kafka
	for serviceName, value := range kafka {
		//start_log(value.Brokers, value.Topic, serviceName, value.Basedir)
		filePath := strings.Join([]string{value.Basedir, serviceName}, string(os.PathSeparator))
		fileMap := gzip.GetAllFile(filePath)
		for fileName, path := range fileMap {
			index := strings.LastIndex(fileName, ".gz")
			if index == -1 {
				f1, _ := os.Open(path)
				files := []*os.File{f1}
				gzip.Compress(files, path+string(os.PathSeparator)+fileName+".gz")
			}
		}
	}
}

//启动多个任务
func Job() {
	i := 0
	c := cron.New()

	//AddFunc
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		i++
		log.Println("cron running:", i)
	})

	c.AddJob(spec, GzipJob{})
	//启动计划任务
	c.Start()

	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer c.Stop()

	select {}
}
