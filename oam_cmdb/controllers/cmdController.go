package controllers

import (
	"OAM/cmd"
	"net"
	"strings"
	"time"

	fn "OAM/util"

	"github.com/google/uuid"
)

type CmdController struct {
	BaseController
}

func (c *CmdController) PingHost() {
	ip := c.GetString("ip")
	if ip == "" {
		c.JsonParamError("请输入IP地址")
	}
	ip = strings.TrimSpace(ip)
	if net.ParseIP(ip) == nil {
		c.JsonParamError("请输入正确的IP地址")
	}
	reqToken := uuid.New().String()
	go func() {
		stats := cmd.Ping(ip, 10, time.Second*15)
		fn.GetPublicCache().PutWithExpireTime(reqToken, stats, time.Minute*5)
	}()
	localIP := fn.LocalIp()
	r := map[string]string{
		"localIp": localIP,
		"token":   reqToken,
	}
	c.JsonOk(r)
}

func (c *CmdController) GetCmdResult() {
	token := c.GetString("token")
	if token == "" {
		c.JsonParamError("参数错误")
	}
	pingStats := fn.GetPublicCache().GetString(token)
	if pingStats == "" {
		c.JsonStatus(203, "执行中")
	} else {
		c.JsonOk(pingStats)
	}
}

func (c *CmdController) PortTest() {
	ip := c.GetString("ip")
	if ip == "" {
		c.JsonParamError("请输入IP地址")
	}
	ip = strings.TrimSpace(ip)
	if net.ParseIP(ip) == nil {
		c.JsonParamError("请输入正确的IP地址")
	}
	port, err := c.GetInt("port")
	if err != nil {
		c.JsonParamError("请输入正确的端口号")
	}
	c.JsonOk(cmd.IsPortAvailable(ip, port))
}
