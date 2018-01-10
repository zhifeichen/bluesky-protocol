
## 代理程序


### 主要业务

    1. 代理程序,主要业务是: 解析不同厂商的设备数据,封装成通用协议分发给各个服务端
    2. 检测设备状态,本机状态并汇报
    
    
    
    
       

### 依赖库

#### 1. 包管理工具 glide

    安装  brew install glide
    
    
### 初始化

#### 初始化工程
    glide init


### 代码目录结构

    ROOT
        |__ doc                                             文档目录
        |__ logs                                            日志目录
        |__ readme.md                                       README文件
        |__ cfg                                             配置
        |__ agent                                           代理程序模块
            |__ servers                                     各种服务
                |__ blueSkyProtocol                         bluesky协议服务,包括tcp/udp侦听以及handler
                
                
        |__ common                                          通用模块
            |__ logger                                      日志模块
                |__ readme.md                               使用说明
            |__ models                                      模型
            |__ utils                                       通用帮组类
            
            |__ chains                                      流处理组件
                |__ readme.md                               使用说明
            |__ tcpServer                                   tcp/udp 通讯框架
            |   __ readme.md                                使用说明
            
            |__ protocol
                |__ bluesky                                 bluesky协议,飞哥用
                    |__ server                              bluesky协议tcp/udp服务类以及codec类
                    ...                                     协议内容
                    
                |__ jianchi                                 
                    

                  
### 库
                    
#### xlogger 日志库

[日志框架说明](common/logger/readme.md)

#### tcp/udp server framework   socket服务框架

[tcp/udp 框架说明](common/tcpServer/readme.md)
                
#### chains framework  流处理框架

[流处理框架说明](common/chains/readme.md)


### 流程

### 安装

### 协议参照


### changelog