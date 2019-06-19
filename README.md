# 《分布式缓存-原理、架构及Go语言实现》

## 第一章
- 一个缓存接口
- 一个内存缓存实现
- 一个HTTP Server服务
- 一个缓存处理器，通过Sever，操作Server中的Cache实现类
- 一个状态处理器，通过Sever，操作Server中的Cache实现类的状态类

## 第二章
- 将通信协议由HTTP协议转换为使用TCP协议
- set格式 `S<klen><SP><vlen><SP><key><value> // <SP>为空格`
