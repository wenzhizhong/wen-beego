[中文](readme-zh.md) | [English](readme.md)
## Project Introduction
  #### Basic information
  A web system based on beego, which aims to achieve basic functions and facilitate rapid project secondary development.
  1. Login and permission management
  2. Menu management
  3. User management
  4. Role management
  5. Log management
  6. Document management
  7. Configuration management
  8. Scheduled task management
  9. Code generation
  If the business model is B2C (pure self support: enterprise ->consumer), use the admin_org module,
  If the business model is B2B2C (platform mode: platform ->(merchant ->consumer)), use the admin_plat module

  #### Directory structure
  ```
  ├─apps                    # Applications
  │  ├─admin_org            # Application: merchant management
  │  ├─admin_plat            # Application: platform management
  │  │  ├─controllers       # Controllers
  │  │  ├─models            # Models   
  │  │  ├─services          # Business logic    
  │  │  └─views             # Views     
  │  ├─common               # Public modules    
  │  │  ├─controller        # base controller
  │  │  ├─global            # global
  │  │  ├─helper            # helper
  │  │  ├─middleware        # middleware
  │  │  └─models            # public models    
  │  └─index                # Applications: Other custom applications
  │      ├─controllers
  │      ├─models
  │      ├─services
  │      └─views
  ├─beego-docker            # container configuration
  ├─conf                    # configuration
  ├─doc                     # documentation
  ├─routers                 # routing
  ├─static                  # static
  ├─temp                    # temporary file
  │  └─logs 
  ├─tests                   # Test
  ```

## Install & Run
  #### environment
  go 1.24.4
  postgresql 17.5
  redis 7.0.10

  #### Clone Code
  ```
  git clone https://github.com/wenzhizhong/wen-beego.git
  ```

  #### Install the bee command-line tool
  ```
  go install github.com/beego/bee/v2@latest
  ```

  #### Set environment variables
  conf/app.yaml is the default configuration of the development environment, and the development environment skips this step.
  
  For the test environment, you can set the environment variable to `BEEGO_RUNMODE` `test`, and the official environment to set the variable to `BEEGO_RUNMODE` `prod`
  ```
  # Linux/macOS
  export BEEGO_RUNMODE=test
  
  # Windows
  set BEEGO_RUNMODE=test
  ```

  #### run
  Run postgresql, redis.

  I'm using dockers here, starting directly in the terminal, and you can automatically build the container `service:postgresql``redis`
  ```
  cd beego-docker
  docker-compose up -d
  ```

  1. development environment, which can be run directly.
  If you use vender mode, you can add more parameters to the command. `bee run-vendor=true`
  ```
  go mod tidy
  bee run 
  # bee dlv -port 8013 # vscode debug
  ```
  2. test environment, commands can be added and then runbee run-runmode=test
  ```
  go mod tidy
  go mod vendor
  bee run -runmode=test
  ```
  3. Build environment，[production environment deployment](doc/生产环境部署.md)


## bee command adds instructions
  #### 1. Generate API documentation
  Tests have found that the bee generate docs command doesn't work with the current project's routing(),
  so use swaggo/swag instead of generate documentation:routers/xxx_router.go
  ```
  go install github.com/swaggo/swag/cmd/swag@latest

  swag init -g cmd/http/main.go  --output apps/swagger 
  ```
  After the command is successfully executed, the swagger.json and swagger.yaml files are generated in the apps/swagger directory, and
  the document access address is http://127.0.0.1:8031/swagger/index.html

# System manual 
  [System manual](doc/系统手册.md)
