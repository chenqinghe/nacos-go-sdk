# nacos-go-sdk
go sdk for nacos. provide friendly api to use nacos.



## client 
client represent a nacos client. All services need nacos client. 

to create a nacos client:
```go

baseUrl:="http://127.0.0.1:8848"
client:= nacos.NewNacosClient(baseUrl)
```


## config

#### crate config service 

```go

configService:= nacos.NewConfigService(client)

```

#### get config
```go
var (
	namespace = "your namespace id"
	group = "your group"
	dataId = "your dataId"
)
configData,err:= configService.GetConfig(namespace,group,dataId)
if err!=nil {
	// handle error
}
```


#### publish config
```go

var data = `
{
    "addr":"127.0.0.1",
    "port": 8080
}
` 

err:= configService.PublishConfig(namespace,group,dataId,[]byte(data))
if err!=nil {
	// handle error ....
}

```



## service 