#### 安装依赖中间件

    # 安装 redis
    docker run -d --name redis -p 6379:6379 hub.randjlighting.com/library/redis:6.2.14

    # 安装 mysql
    docker run -d --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 hub.randjlighting.com/library/mysql:8.0.38



#### 构建项目

    go mod tidy
    cd api
    go build -o fca -t
    ./fca

    under windows
    go mod tidy
    cd api
    go build -x

  
  
#### 验证

    http://127.0.0.1:8888/v1/hello/you
    http://127.0.0.1:8888/v1/hello/me

 
#### 使用


    # 新增请求和响应结构
    vim api/fca.api
  
    # 生成routes，types，handler 及 logic
    goctl api go -api ./api/fca.api -dir ./api
  
    # 根据ddl 语句生成模型
    goctl model mysql ddl -src ./model/user.sql -dir ./model -c

    # 根据datasource 语句生成模型
    goctl model mysql datasource -url "username:password@tcp(flopscloud:3306)/flops_cloud" -t user -dir ./model


#### 热加载（可选）

    # 安装全局命令
    go install github.com/air-verse/air@latest
    # 验证
    air -v
    # 进入项目（可执行 go run/build的目录）
    cd flopsServer/api
    # 初始化配置文件，会在当前目录生成 .air.toml
    air init
    # 启动服务
    air
    
> 参考地址：`go install github.com/air-verse/air@latest`
