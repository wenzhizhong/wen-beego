[中文](readme-zh.md) | [English](readme.md)
## 项目介绍
  #### 基本信息
  一个基于beego的web系统，目标实现基础功能，方便快速项目二次开发。
  1. 登录、权限管理
  2. 菜单管理
  3. 用户管理
  4. 角色管理
  5. 日志管理
  6. 文件管理
  7. 配置管理
  8. 定时任务管理
  9. 代码生成
  如果经营模式是B2C (纯自营:企业->消费者)，使用admin_org模块即可，
  如果经营模式是B2B2C (平台模式:平台->(商户->消费者))，使用admin_plat模块即可

  #### 目录结构
  [目录结构](doc/目录结构.md)

## 安装&运行
  #### 环境
  go 1.24.4
  postgresql 17.5
  redis 7.0.10

  #### 克隆代码
  ```
  git clone https://github.com/wenzhizhong/wen-beego.git
  ```

  #### 安装bee 命令行工具
  ```
  go install github.com/beego/bee/v2@latest
  ```

  #### 设置环境变量
  conf/app.yaml 默认配置就是开发环境，开发环境跳过这步骤。
  测试环境，可以设置环境变量`BEEGO_RUNMODE`为 `test` ，正式环境设置变量`BEEGO_RUNMODE`为 `prod`
  ```
  # Linux/macOS
  export BEEGO_RUNMODE=test
  
  # Windows
  set BEEGO_RUNMODE=test
  ```

  #### 运行
  启动postgresql, redis。

  我这里是使用dockers，直接在终端启动，即可自动构建`postgresql``redis`容器服务：
  ```
  cd beego-docker
  docker-compose up -d
  ```

  1. 开发环境，可以直接运行。
  如果使用vender 模式，`bee run`命令可以再添加参数 `-vendor=true`。
  ```
  go mod tidy
  bee run 
  # bee dlv -port 8013 # vscode debug
  ```
  2. 测试环境，`bee run`命令可以添加`-runmode=test`再运行
  ```
  go mod tidy
  go mod vendor
  bee run -runmode=test
  ```
  3. 生成环境，[生产环境部署](doc/生产环境部署.md)

## bee命令补充说明
  #### 1. 生成api文档
  测试发现bee generate docs 命令不适用当前项目的路由（`routers/xxx_router.go`）的代码写法，
  所以使用swaggo/swag代替生成文档：
  ```
  go install github.com/swaggo/swag/cmd/swag@latest

  swag init -g cmd/http/main.go  --output apps/swagger 
  ```
  命令运行成功后，会在apps/swagger目录下生成swagger.json和swagger.yaml文件，
  文档访问地址：http://127.0.0.1:8013/swagger/index.html

# 系统手册 
  [系统手册](doc/系统手册.md)
