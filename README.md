# plugin-gb28181
gb28181 plugin for monibuca


- [] 单端口
invite的时候设备会返回一个sdp信息，里面包含了媒体端口
然后我们需要建立一个映射，设备IP和媒体端口->streamPath
server的媒体端口就需要一开始就监听了

- [] 回放流自动停止

- [] 自动删除设备  1个小时没有更新的