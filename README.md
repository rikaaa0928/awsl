# awsl
A WebSocket Linker  

## How To

windows10可使用wsl

1. 安装[golang](https://golang.org/doc/install "go")  
2. `git clone https://github.com/Evi1/awsl.git`
3. `cd awsl`
4. `chmod +x install_all.sh`
5. `./install_all.sh`
6. 运行`$GOPATH/bin/awsl -c $configfile`

没有`-c $configfile`时读取/etc/awsl/config.json
配置文件参考`test/server.config.json`和`test/client.config.json`

## WSL解决方案
1. 创建wsl_run.sh
```
#!/bin/bash
# service privoxy start >/dev/null 2>&1    # 可使用privoxy完成http -> socks5
$GOPATH/bin/awsl > $logfile 2>&1
```
2. `chmod +x wsl_run.sh`
3. `sudo visudo`  
添加  `$username ALL=(root) NOPASSWD: $path_to_wsl_run`
4. windows中新建`wsl.vbs`
```
Set ws = CreateObject("Wscript.Shell") 
ws.run "wsl -d Ubuntu-18.04 -e sudo $path_to_wsl_run.sh", vbhide
```
5. 使用任务计划程序添加vbs