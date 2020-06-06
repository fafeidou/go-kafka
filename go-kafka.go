package main

import (
	"bufio"
	"fmt"
	"github.com/Shopify/sarama"
	"go-kafka/setting"
	gpool "go-kafka/util"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// kafka consumer
func init() {
	setting.Setup()
}

var (
	pool *gpool.Pool
)

func main() {

	pool = gpool.New(len(setting.App.Kafka))
	kafka := setting.App.Kafka
	for serviceName, value := range kafka {
		pool.Add(1)
		go start_log(value.Brokers, value.Topic, serviceName, value.Basedir)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	for {
		s := <-quit
		switch s {
		case syscall.SIGINT | syscall.SIGTERM:
			pool.Done()
			log.Println("Shutdown Server ...")
		default:
			log.Println("default Shutdown Server ...")

		}
	}
	pool.Wait()
}

func start_log(brokers []string, topic string, serviceName string, baseDir string) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}

	}()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			os.Mkdir(baseDir+serviceName, os.ModePerm)
			timeStr := time.Now().Format("2006-01-02")
			file, _ := os.OpenFile(baseDir+serviceName+"/"+serviceName+timeStr+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 755)
			writer := bufio.NewWriter(file)
			_, err = writer.WriteString(string(msg.Value[:]) + "\n") //将数据先写入缓存
			writer.Flush()
			file.Sync()
			if err != nil {
				fmt.Println("write file failed, err:", err)
			}
			log.Printf("Consumed message offset %d,msg value %s\n", msg.Offset, msg.Value)
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Printf("Consumed: %d\n", consumed)
}
