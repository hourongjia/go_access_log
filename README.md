# 简介
gin_accesslog 是gin框架一个中间件，主要作用是记录accesslog日志，打印接口执行时间，请求入参，返回信息等内容

# ginzap缺陷
ginzap也是一个记录accesslog的gin插件，但是目前来看，它无法做到打印请求body信息和返回的body信息，gin_access_log_interceptor解决了这个问题

# 使用方法


## 集成到gin中
``` 
r := gin.Default()
log, _ := zap.NewDevelopment()
r.Use(gin_accesslog.Ginzap(log, time.RFC3339, true))
``` 

# 单元测试

![img.png](img.png)

测试覆盖率：100%




