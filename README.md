# 电商系统mxshop-srv

> 客户端api调用:https://github.com/heqiang/mxshop-api.git
>
> 由于阿里云的短信服务暂时未开通成功，所以用户相关手机短信注册等功能暂不能实现，逻辑已完善，

##### 一、api功能模块：  

+ 1 用户模块  
  + 用户登录
  + 用户注册
  + 用户信息修改 
  + 用户列表查看  
+ 用户权限验证
  + jwt

环境:

+ sdk:go1.17.3 windows/amd64  

+ 编辑器：goland 
+ 技术栈 go、gin、grpc、consul、nacos、docker
+ 数据库：mysql8、redis5.0

####  二、安装下载

2.1、下载:

api客户端

>  git@github.com:heqiang/mxshop-api.git
>
> go mod tidy

2.2 相关服务配置

grpc服务

> cd proto
>
> protoc --go_out=plugins=grpc:. *.proto

数据库服务

> ```console
> docker run --name some-mysql -e MYSQL_ROOT_PASSWORD=123456 -p 3306:3306 -d mysql:8.0
> docker run --name some-redis -p 6379:6379 -d redis
> ```

consul注册

> docker run -d --name=dev-consul  -p 8500:8500  -e CONSUL_BIND_INTERFACE=eth0 consul

配置中心nacos

> ```shell
> docker run --name nacos-quick -e MODE=standalone -p 8848:8848 -d nacos/nacos-server:2.0.2
> ```

进入nacos

step1 http://127.0.0.1:8848/nacos

step2 登录 账号密码都是nacos

step3 根据config.yaml对nacos进行相关配置，只需要配置命名空间及配置列表，配置信息如下

```json
{
  "name": "mxshop_srv",
  "mysql": {
    "host": "127.0.0.1",
    "port": 3306,
    "user": "root",
    "password": "142212",
    "db": "mxshop_user_srv"
  },
  "log": {
    "level": "debug",
    "filename": "web_app_log.log",
    "max_size": 200,
    "max_age": 30,
    "max_backups": 7
  },
  "consul": {
    "host": "127.0.0.1",
    "port": 8500,
      # 更换本机的ipv4地址
    "serverhost": "192.168.31.101"
  }
}
```

#### 三 运行：

> go  run main.go



