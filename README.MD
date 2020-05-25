# project is developing

Base on [miniaudio](https://github.com/dr-soft/miniaudio)  

There's no need to install any dependencies

## run

1. go run rpc_server.go
2. go run ui_client.go   [  -dir directory ] ( default  ./ )

## Attention
### To avoid panic: argument has Go pointer to Go pointer
export GODEBUG=cgocheck=0



## Todo

- Support mutil-diretory
- Change a few opearation
- Support Minio
