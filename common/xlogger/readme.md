## 日志服务

    Golang 日志框架.


### 使用说明


#### 初始化

    xlogger.New("./logs/access.log", xlogger.DEBUG, true)
        
#### 使用
    
    xlogger.Debug(v...)
    xlogger.Info(v...)
	xlogger.Warn(v...)
	xlogger.Error(v...)
	
	xlogger.Debugf(fomat,v...)
	xlogger.Infof(fomat,v...)
    xlogger.Warnf(fomat,v...)
    xlogger.Errorf(fomat,v...)
    	
		
#### 关闭
		
	defer xlogger.Close()
		