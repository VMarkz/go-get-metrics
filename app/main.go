package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type path struct {
	Path string `json:"path"`
}

type metrics struct {
	Latency int64 `json:"latency"`
}

func string_body(body []byte) string {
	return string(body)
}

func get_body(resp *http.Response) []byte {
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	return body
}

func make_get_request(path string) *http.Response {
	resp, err := http.Get(path)

	if err != nil {
		log.Fatalln(err)
	}

	return resp
}

type fn_def func(string) string

func get_path(path string) string {
	return string_body(get_body(make_get_request(path)))
}

func time_to_elapse(fn fn_def, path string) int64 {
	init_time := time.Now().UnixMilli()
	fn(path)
	end_time := time.Now().UnixMilli()

	return end_time - init_time
}

func get_reponse_time(context *gin.Context) {
	var path path

	if err := context.BindJSON(&path); err != nil {
		return
	}

	var metrics metrics = metrics{Latency: time_to_elapse(get_path, path.Path)}

	context.IndentedJSON(http.StatusOK, metrics)
}

func main() {
	router := gin.Default()
	router.GET("/go-get-metrics/metrics", get_reponse_time)
	router.Run("localhost:8080")
}
