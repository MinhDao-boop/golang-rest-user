package main

import (
	"golang-rest-user/provider/jwtProvider"
	"golang-rest-user/provider/mySqlProvider"
	"golang-rest-user/provider/redisProvider"
	"golang-rest-user/provider/routesProvider"
	"golang-rest-user/provider/serviceProvider"
	"golang-rest-user/provider/tenantProvider"

	"github.com/gin-gonic/gin"
)

func main() {

	redisProvider.Init()
	mySqlProvider.Init()
	jwtProvider.Init()
	r := gin.Default()

	serviceProvider.Init()

	tenantConnectSvc := serviceProvider.GetInstance().TenantConnectService

	tenantProvider.Init(tenantConnectSvc)

	routesProvider.Init(r)

	err := r.Run(":8080")
	if err != nil {
		return
	}

}
