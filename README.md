# Golang项目基础库

## 配置管理

封装Viper库，可直接解析项目当前或父级conf目录下面的config.yaml配置文件；

#### Example

simple demo auto parse config file

```go
import "github.com/deng00/go-base/config"
configIns := config.GetConfigFromLocal()
configIns.GetString("level")
```

## 日志管理

封装Zap日志库，减少不必要的配置项；

提供输出日志到多种渠道： Console、本地文件、日志中心（待完成）；

增加实时告警，当前支持DingDing告警；

可自动根据配置文件初始化日志服务，只需在配置中配置告警渠道即可实现自动告警；

#### 日志示例

simple demo show log on console

```go
import "github.com/deng00/go-base/logging"
logConfig := &logging.LogConfig{}
logConfig.EnableDebug()
logConfig.EnableHandlerFile()
logger := logging.GetLogger("testServiceName", logConfig).Sugar()
logger.Info("info", logging.String("url", "https://www.baidu.com"))

```

Send alert realtime and save log to file, files will create under current dir `logs` named as the service name you set.
Single log file max size up to 1 GB .

```go
logConfig := &logging.LogConfig{}
logConfig.EnableHandlerStd()
logConfig.SetAlertChannel(logging.NewDingDingAlertChanel("XXX"))
logConfig.SetAlertLevel(logging.ErrorLevel)
logger := logging.GetLogger("testServiceName", logConfig).Sugar()
logger.Errorf("info key %s", "value")
```

根据使用配置实例自动初始化：

```go
configIns := config.GetConfigFromLocal()
LogConfig := logging.GetLogConfig(configIns)
logger := logging.GetLogger("testServiceName", logConfig).Sugar()
logger.Infof("info key %s", "value")
```


## 服务注册

封装consul服务，提供自动服务注册和发现；

提供动态配置文件管理，当consul的配置(key)修改后本地服务config实例自动同步；

## 日志管理Manager

All in One

```go
var configIns = &config.Config{}
configManager := config.Manager{}
_ = configManager.Init("go-base-demo")
configIns = configManager.GetIns()
LogConfig := logging.GetLogConfig(configIns)
logger := logging.GetLogger("test", LogConfig)
logger.Error("test send alert")
```

配置告警使用的场景，配置文件格式格式：

```yaml
log:
  handlerStd: true
  handlerFile: true
  fileDir: logs
  level: info
  alert:
    level: error
    channel: dingding
    channel_param: 机器人Token
```