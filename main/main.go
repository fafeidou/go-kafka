package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func main() {
	var c Conf
	conf := c.getConf()
	fmt.Println(conf.Kafka)
}

//profile variables
type Conf struct {
	Kafka map[string]LogInfo `yaml:"kafka"`
}

type LogInfo struct {
	Brokers    []string `yaml:"brokers"`
	Topic      string   `yaml:"topic"`
	Basedir    string   `yaml:"baseDir"`
	Group      string   `yaml:"group"`
	MaxSize    int      `yaml:"maxSize"`
	MaxBackups int      `yaml:"maxBackups"`
	MaxAge     int      `yaml:"maxAge"`
}

func (c *Conf) getConf() *Conf {
	yamlFile, err := ioutil.ReadFile("../conf/conf.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}
