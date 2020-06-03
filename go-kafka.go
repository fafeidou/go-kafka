package main

import (
	"bufio"
	"fmt"
	"github.com/Shopify/sarama"
	"go-kafka/setting"
	"log"
	"os"
	"os/signal"
	"sync"
)

// kafka consumer
func init() {
	setting.Setup()
}

var wg *sync.WaitGroup

func main() {
	kafka := setting.App.Kafka
	wg = &sync.WaitGroup{}
	for serviceName, value := range kafka {
		wg.Add(1)
		go start_log(value.Brokers, value.Topic, serviceName)
	}
	wg.Wait()

}

func start_log(brokers []string, topic string, serviceName string) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
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
			os.Mkdir(serviceName, os.ModePerm)
			file, _ := os.OpenFile(serviceName+"/"+serviceName+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 755)
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
