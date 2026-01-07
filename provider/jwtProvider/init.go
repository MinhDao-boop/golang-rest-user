package jwtProvider

import "golang-rest-user/security"

var instance *security.Manager

func Init() {
	jwtConfig := security.LoadJWTConfig()
	instance = security.NewManager(jwtConfig)
}

func GetInstance() *security.Manager {
	return instance
}
