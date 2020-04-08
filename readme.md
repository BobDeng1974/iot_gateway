- 小度技能平台网址 https://dueros.baidu.com/dbp/debug/index#/audio/?botId=38657853-65da-c980-50f9-6fa6309c9756

- ssh -gR :8000:localhost:8000 root@aliyun

proto 文件编译方法
protoc --go_out=plugins=grpc:. parser.proto
protoc --go_out=. gateway.proto

rsync -zvr ./protobuf_demo root@aliyun:/root/temp
rsync -zvr ./protobuf_demo/*.c root@aliyun:/root/temp/protobuf_demo/