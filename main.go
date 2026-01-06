package main

import (
	"golang-rest-user/config"
	"golang-rest-user/provider/mySqlProvider"
	"golang-rest-user/provider/routesProvider"
	"golang-rest-user/provider/serviceProvider"
	"golang-rest-user/provider/tenantProvider"

	"github.com/gin-gonic/gin"
)

func main() {

	config.InitRedis()
	mySqlProvider.Init()
	r := gin.Default()

	serviceProvider.Init()

	tntCntSvc := serviceProvider.GetInstance().TenantConnectService

	tenantProvider.Init(tntCntSvc)

	routesProvider.Init(r)

	err := r.Run(":8080")
	if err != nil {
		return
	}

}
