package main

const toml = `
[dev]
[dev.server]
addr = ":8080"
[dev.mysql]
addr = "127.0.0.1:3306"
username = "root"
password = "password"
dbName = "demo"
[dev.redis]
addr = "127.0.0.1:6379"
password = "password"
db = 0
poolSize = 100
#自定义配置项

[test]
[test.server]
addr = ":8080"
[test.mysql]
addr = "127.0.0.1:3306"
username = "root"
password = "password"
dbName = "demo"
[test.redis]
addr = "127.0.0.1:6379"
password = "password"
db = 0
poolSize = 100
#自定义配置项

[pre]
[pre.server]
addr = ":8080"
[pre.mysql]
addr = "127.0.0.1:3306"
username = "root"
password = "password"
dbName = "demo"
[pre.redis]
addr = "127.0.0.1:6379"
password = "password"
db = 0
poolSize = 100
#自定义配置项

[prod]
[prod.server]
addr = ":8080"
[prod.mysql]
addr = "127.0.0.1:3306"
username = "root"
password = "password"
dbName = "demo"
[prod.redis]
addr = "127.0.0.1:6379"
password = "password"
db = 0
poolSize = 100
#自定义配置项
`
