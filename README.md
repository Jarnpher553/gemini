[![Go Report Card](https://goreportcard.com/badge/github.com/Jarnpher553/gemini)](https://goreportcard.com/report/github.com/Jarnpher553/gemini)
![GitHub](https://img.shields.io/badge/license-MIT-brightgreen)

# gemini

## Introduction

GEMINI is a microservice framework

## Quick Start

```go
package main

import (
    _ "demo/schedule"
    _ "demo/task"
    
    
    "os"
    "demo/services"

    "github.com/Jarnpher553/gemini/config"
    "github.com/Jarnpher553/gemini/redis"
    "github.com/Jarnpher553/gemini/repo"
    "github.com/Jarnpher553/gemini/router"
    "github.com/Jarnpher553/gemini/scheduler"
    "github.com/Jarnpher553/gemini/server"
    "github.com/Jarnpher553/gemini/service"
    "github.com/Jarnpher553/gemini/task"
    "github.com/Jarnpher553/gemini/task/delay"
)

func main() {
    env := config.DeployEnv(os.Args)
    
    // 获取对应项的配置
    serverCf := config.Conf().Sub("server")
    mysqlCf := config.Conf().Sub("mysql")
    redisCf := config.Conf().Sub("redis")
    
    // 数据库实例
    db := repo.New(repo.DbName(mysqlCf.GetString("dbName")), repo.Addr(mysqlCf.GetString("addr")), repo.Pwd(mysqlCf.GetString("password")), repo.UserName(mysqlCf.GetString("username")), repo.LogMode(false))
    rd := redis.New(redis.Pwd(redisCf.GetString("password")), redis.PoolSize(redisCf.GetInt("poolSize")), redis.DB(redisCf.GetInt("db")), redis.Addr(redisCf.GetString("addr")))
    
    
    // 实例化服务
    userService := service.NewService(&services.UserService{}, service.Repository(db), service.RedisClient(rd))
    
    scheduler.Bind(scheduler.Repo(db), scheduler.Redis(rd))
    delay.Bind(true, task.Redis(rd), task.Repo(db))
    
    // 实例化路由
    r := router.New()
    // 将服务注册进路由
    r.InjectSlice(userService)
    
    // 实例化服务器
    srv := server.Default(server.Name("demo"), server.Env(env), server.RunMode(serverCf.GetString("runMode")), server.Router(r), server.Addr(serverCf.GetString("addr")))
    
    // 运行服务器
    srv.Run()
}
```
