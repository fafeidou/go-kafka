package gzip

import (
	"go-kafka/setting"
	"os"
	"strings"
	"testing"
)

func init() {
	setting.Setup()
}

func TestCompress(t *testing.T) {
	kafka := setting.App.Kafka
	for serviceName, value := range kafka {
		//start_log(value.Brokers, value.Topic, serviceName, value.Basedir)
		filePath := strings.Join([]string{value.Basedir, serviceName}, string(os.PathSeparator))
		fileMap := GetAllFile(filePath)
		for fileName, path := range fileMap {
			index := strings.LastIndex(fileName, ".gz")
			if index == -1 {
				f1, _ := os.Open(path)
				files := []*os.File{f1}
				Compress(files, path+string(os.PathSeparator)+fileName+".gz")
			}
		}
	}

}
