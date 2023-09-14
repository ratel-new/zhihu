package main

import (
	"climb/api"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	router := fasthttprouter.New()
	router.GET("/", api.ShareCore)
	log.Fatal(fasthttp.ListenAndServe(":12345", router.Handler))
}
