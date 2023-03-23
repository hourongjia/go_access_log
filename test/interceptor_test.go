package test

import (
	"errors"
	"fmt"
	"gin_accesslog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/go-playground/assert.v1"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func init() {
	// 启动项目
	go start()
}

func start() {
	r := gin.Default()
	log, _ := zap.NewDevelopment()
	r.Use(gin_accesslog.Ginzap(log, time.RFC3339, true))

	// 正常无请求参数，有返回参数的请求
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// 有请求参数，无返回参数的请求

	r.GET("/v1", v1)

	// 有请求参数，有返回参数的请求
	r.GET("/v2", v2)

	// panic 方法
	r.GET("/panic", v5)

	// test post
	r.POST("/v3", v3)

	// test error
	r.POST("/v4", v4)

	r.Run()
}

func v1(c *gin.Context) {
	type Params struct {
		OS string `form:"os"`
	}
	param := &Params{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(201, "error")
	}
}

func v2(c *gin.Context) {
	type Params struct {
		OS string `form:"os"`
	}
	param := &Params{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(201, "error")
	}
	c.JSON(200, "success")

}

func v4(c *gin.Context) {

	c.Error(errors.New("hourongjia is error"))
	c.JSON(200, "success")

}

func v5(c *gin.Context) {
	panic("An unexpected error happen!")
}

func TestV1(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/v1?os=1") // url
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	assert.Equal(t, string(body), "")
}

// TestV2 有返回的
func TestV2(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/v2?os=1") // url
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	str := string(body)
	fmt.Println(str)
}

func TestPanic(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/panic?os=1") // url
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	str := string(body)
	fmt.Println(str)
}

func TestPost(t *testing.T) {
	resp, err := http.Post("http://localhost:8080/v3", // url
		"text/json;charset=utf-8",                // contentType 内容类型
		strings.NewReader("{\"name\":\"jack\"}")) // body请求体
	if err != nil {
		fmt.Printf("post请求失败 error: %+v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取Body失败 error: %+v", err)
		return
	}
	fmt.Println(string(body))

}

func TestError(t *testing.T) {
	resp, err := http.Post("http://localhost:8080/v4", // url
		"text/json;charset=utf-8",                // contentType 内容类型
		strings.NewReader("{\"name\":\"jack\"}")) // body请求体
	if err != nil {
		fmt.Printf("post请求失败 error: %+v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取Body失败 error: %+v", err)
		return
	}
	fmt.Println(string(body))

}

func v3(c *gin.Context) {

	body, err := c.GetRawData()
	if err != nil {
		c.JSON(201, "error")
	}
	fmt.Println(string(body))

	type Params struct {
		OS string `form:"os"`
	}
	param := &Params{}
	err = c.Bind(param)
	if err != nil {
		c.JSON(201, "error")
	}
	c.JSON(200, "success")
}
