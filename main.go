package main

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

var config = []string{"kubelet.yml"}

type kubeletConf struct {
	Kubelet_opts string `yaml:"cluster_hostname"`
}

func main(){
	fmt.Println("Hello World!")
	fmt.Println(config[0])
	var data,_ = ioutil.ReadFile(config[0])
	kc := &kubeletConf{}
	yaml.Unmarshal(data,kc)
  	fmt.Printf("----test\n %s",(*kc).Kubelet_opts)

}


