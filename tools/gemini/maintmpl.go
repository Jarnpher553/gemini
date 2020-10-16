package main

const mainTmpl = `
package main

import (
	"github.com/Jarnpher553/gemini/config"
	"github.com/Jarnpher553/gemini/repo"
	"github.com/Jarnpher553/gemini/router"
	"github.com/Jarnpher553/gemini/server"
	"github.com/Jarnpher553/gemini/service"
	"github.com/Jarnpher553/gemini/redis"
	"{{name}}/services"
	"os"

	//_ "{{name}}/validators"
	//_ "{{name}}/error"
	//_ "{{name}}/schedules"
)

func main() {
	env := config.DeployEnv(os.Args)

	// 获取对应项的配置
	serverCf := config.Conf().Sub("server")
	mysqlCf := config.Conf().Sub("mysql")
	redisCf := config.Conf().Sub("redis")

	// 实例化服务器
	srv := server.Default(server.Name("api"), server.Env(env), server.Addr(serverCf.GetString("addr")), server.Startup(func(s *server.DefaultServer) error {
		// 数据库实例
		db := repo.New(repo.DbName(mysqlCf.GetString("dbName")), repo.Addr(mysqlCf.GetString("addr")), repo.Pwd(mysqlCf.GetString("password")), repo.UserName(mysqlCf.GetString("username")))
		// 在此处迁移初始化数据库
		db.Migrate(nil, nil)
	
		// redis实例化
		rd := redis.New(redis.Pwd(redisCf.GetString("password")), redis.PoolSize(redisCf.GetInt("poolSize")), redis.DB(redisCf.GetInt("db")), redis.Addr(redisCf.GetString("addr")))

		//初始化定时任务
		//scheduler.Bind(scheduler.Repo(db))

		//初始化email组件
		//email.Bind(email.Host("..."))

		// 实例化服务
		{{range .}}{{ .Name }} := service.NewService(&services.{{ title .Name }}{}, service.Repository(db), service.RedisClient(rd)){{ end }}

		// 实例化路由
		r := router.New()
		// 将服务注册进路由
		r.Assign({{services .}})
		
		s.Serve(r)

		return nil
	}))

	// 运行服务器
	srv.Run()
}
`
