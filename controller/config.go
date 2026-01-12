package controller

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/housepower/ckman/config"
	"github.com/housepower/ckman/model"
)

type ConfigController struct {
	Controller
	config *config.CKManConfig
	signal chan os.Signal
}

func NewConfigController(ch chan os.Signal, config *config.CKManConfig, wrapfunc Wrapfunc) *ConfigController {
	cf := &ConfigController{}
	cf.signal = ch
	cf.config = config
	cf.wrapfunc = wrapfunc
	return cf
}

// 该接口不暴露给用户
func (controller *ConfigController) GetVersion(c *gin.Context) {
	version := strings.Split(config.GlobalConfig.Version, "-")[0]
	controller.wrapfunc(c, model.E_SUCCESS, version)
}

func (controller *ConfigController) GetInstances(c *gin.Context) {
	var instances []string
	config.ClusterMutex.RLock()
	for _, instance := range config.ClusterNodes {
		instances = append(instances, net.JoinHostPort(instance.Ip, fmt.Sprint(instance.Port)))
	}
	config.ClusterMutex.RUnlock()

	if len(instances) == 0 {
		instances = append(instances, net.JoinHostPort(controller.config.Server.Ip, fmt.Sprint(controller.config.Server.Port)))
	}

	controller.wrapfunc(c, model.E_SUCCESS, instances)
}
