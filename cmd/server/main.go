package main

import (
	pm "TapMars/admin_gateway/pkg/productManager"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	ctx := context.Background()

	var addr string
	//var opts []grpc.DialOption
	var timeout time.Duration
	gin.SetMode(gin.ReleaseMode)
	timeoutSec := os.Getenv("TIMEOUT_SEC")
	if timeoutSec == "" {
		timeout = time.Second * 10
	} else {
		sec, err := strconv.Atoi(timeoutSec)
		if err != nil {
			log.Fatalf("Parsing TIMEOUT_SEC: %v", err)
		}
		timeout = time.Second * time.Duration(sec)
	}
	pmConn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Product Manager connection: %v", err)
	}
	pmProxy := pm.NewProductManager(pmConn)
	defer func(pm *pm.ProductManager) {
		err := pm.Close()
		if err != nil {
			log.Fatalf("Closing Proxy: %v", err)
		}
	}(pmProxy)

	router := gin.Default()
	router.Use(pm.RequestTimeoutWrapper(timeout))
	pmRouter := router.Group("/product-manager")
	pmProxy.AddRoutes(pmRouter)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("Failed to find Port: %s", port)
	}

	log.Fatal(http.ListenAndServe(":"+port, router))

}
