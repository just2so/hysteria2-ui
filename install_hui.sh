#!/bin/bash
#创建目录
mkdir -p /usr/local/h-ui/
#下载最新版h-ui应用程序，将下载的文件保存到/usr/local/h-ui/h-ui，在下载完成后，给下载的文件添加执行权限
curl -fsSL https://github.com/jonssonyan/h-ui/releases/latest/download/h-ui-linux-amd64 -o /usr/local/h-ui/h-ui && chmod +x /usr/local/h-ui/h-ui
#下载h-ui的systemd 服务文件，并将其保存到 /etc/systemd/system/h-ui.service
curl -fsSL https://raw.githubusercontent.com/jonssonyan/h-ui/main/h-ui.service -o /etc/systemd/system/h-ui.service
#重新加载 systemd 服务管理器，以便识别新添加或修改的服务文件
systemctl daemon-reload
#启动h-ui服务
systemctl enable h-ui
#重启h-ui服务
systemctl restart h-ui
echo
echo "登录方式：IP:PORT"
echo "面板端口：8081"
echo "用户名/密码：sysadmin"
echo
