package main

import (
	"fmt"
	"github.com/LeBronQ/Mobility"
	"net/http"
	"github.com/gin-gonic/gin"
	consulapi "github.com/hashicorp/consul/api"
)

type DiscoveryConfig struct {
	ID      string
	Name    string
	Tags    []string
	Port    int
	Address string
}

var consulAddress = "127.0.0.1:8500"

type ReqParams struct {
	Node    Mobility.Node    `json:"node"`
}

type Respond struct {
	Node    Mobility.Node    `json:"node"`
}

func RegisterService(dis DiscoveryConfig) error {
	fmt.Println("---------")
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Printf("create consul client : %v\n", err.Error())
	}
	registration := &consulapi.AgentServiceRegistration{
		ID:      dis.ID,
		Name:    dis.Name,
		Port:    dis.Port,
		Tags:    dis.Tags,
		Address: dis.Address,
	}
	// 启动tcp的健康检测，注意address不能使用127.0.0.1或者localhost，因为consul-agent在docker容器里，如果用这个的话，
	// consul会访问容器里的port就会出错，一直检查不到实例
	check := &consulapi.AgentServiceCheck{}
	check.TCP = fmt.Sprintf("%s:%d", registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "60s"
	registration.Check = check

	if err := client.Agent().ServiceRegister(registration); err != nil {
		fmt.Printf("register to consul error: %v\n", err.Error())
		return err
	}
	return nil
}

func main() {
	router := gin.Default()
	dis := DiscoveryConfig{
		ID:      "23",
		Name:    "Default_Mobility",
		Tags:    []string{"a", "b"},
		Port:    8888,
		Address: "172.16.232.131", //通过ifconfig查看本机的eth0的ipv4地址
	}
	RegisterService(dis)
	router.POST("/mobility", func(c *gin.Context) {
		var param ReqParams
		if err := c.Bind(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var res Respond
		n := Mobility.UpdatePosition(param.Node)
		res.Node = n
		fmt.Println(n)
		c.JSON(http.StatusOK, res)
	})
	router.Run(":8888")
	



}

