package main

import (
	"miniproject/controller"
	"miniproject/ds"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.10.199"})

	ds.NewDataSource()
	controller.Inject(router)
	router.Run("localhost:8080")
}
