package main

import (
	"bufio"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
)

// kafka consumer

func main() {
	consumer, err := sarama.NewConsumer([]string{"106.14.227.31:9092"}, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("abklog_topic", 0, sarama.OffsetNewest)
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
			file, _ := os.OpenFile("./xx.txt", os.O_CREATE|os.O_APPEND, 0777)
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
