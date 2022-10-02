

#### 一键ssh跳转/上传(scp)/下载(scp)工具

##### 优点
* 支持谷歌验证码

##### 工具使用  
1. 从release下载最新go-ssh-tool.mac(或者go-ssh-tool.linux) 后缀是所使用的环境
2. 下载文件config.yaml(放在上面工具的相同目录)并填入自己的主机信息
3. 运行: 
   1. 查看工具说明: `./go-ssh-tool.mac`  
   2. 查看主机清单: `./go-ssh-tool.mac host`  
   3. 登录一个主机: `./go-ssh-tool.mac term -m 1`   
        > -m: 主机序号  
   4. 上传文件: `./ssh-tool.mac up -m 1 -f xxxx -d /tmp`  
        > -f: 是本地的一个文件  
        > -d: 上传到主机的哪个目录  
   5. 下载文件:  `down -m 1 -s /root/anaconda-ks.cfg`  
        > -s: 服务器上的文件  
        > -t: 本地目录，可不传(默认本机)
   6. 远程执行一个命令: `./ssh-tool.mac exe -m 1 -c "cat /etc/redhat-release"`  
        > -c: 是命令，若包含空格，用引号引起来
*功能*

 - [x] 列出主机
 - [x] 选择主机并执行一个shell命令
 - [x] 从指定主机下载文件
 - [x] 从上传文件至指定主机
 - [x] 登录执行主机

