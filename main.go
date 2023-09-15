package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"zhihu/api"
)

func main() {
	router := fasthttprouter.New()
	router.GET("/", api.ShareCore)
	log.Fatal(fasthttp.ListenAndServe(":12345", router.Handler))
}
