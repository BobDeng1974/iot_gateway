- 小度技能平台网址 https://dueros.baidu.com/dbp/debug/index#/audio/?botId=38657853-65da-c980-50f9-6fa6309c9756

- ssh -gR :8000:localhost:8000 root@aliyun

正常来说，这里应该只做网络层内容。应用层应该是另外的app server来实现

一种设备类型对应一套解析器。linux分离思想

upx ./bin/open

proto 文件编译方法
protoc --go_out=plugins=grpc:. parser.proto
protoc --go_out=. gateway.proto

c 版本解析器编译方法
protoc-c --c_out=. parser.proto

rsync -zvr ./protobuf_demo root@aliyun:/root/temp
rsync -zvr ./protobuf_demo/*.c root@aliyun:/root/temp/protobuf_demo/