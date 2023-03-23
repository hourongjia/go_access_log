package test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func init() {
	// 启动项目
	go start()
}

// BenchmarkWithOutBodyTest 没有requestBody
func BenchmarkWithOutBodyTest(b *testing.B) {
	for n := 0; n < b.N; n++ {
		resp, err := http.Post("http://localhost:8080/v3", // url
			"text/json;charset=utf-8", // contentType 内容类型
			strings.NewReader(""))     // body请求体
		if err != nil {
			fmt.Printf("post请求失败 error: %+v", err)
			return
		}
		defer resp.Body.Close()
	}
}

// BenchmarkWithBodyTest 有requestBody
func BenchmarkWithBodyTest(b *testing.B) {
	for n := 0; n < b.N; n++ {
		resp, err := http.Post("http://localhost:8080/v3", // url
			"text/json;charset=utf-8",                // contentType 内容类型
			strings.NewReader("{\"name\":\"jack\"}")) // body请求体
		if err != nil {
			fmt.Printf("post请求失败 error: %+v", err)
			return
		}
		defer resp.Body.Close()
	}

}
