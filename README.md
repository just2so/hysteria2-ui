* 申请SSL证书

```shell
bash <(curl -Ls https://raw.githubusercontent.com/just2so/hysteria2-ui/main/acme/acme_2.0.sh)
```

* 安装H-UI服务

```shell
bash <(curl -Ls https://raw.githubusercontent.com/just2so/hysteria2-ui/main/h-ui/install_hui.sh)
```

* 卸载H-UI服务

```shell
bash <(curl -Ls https://raw.githubusercontent.com/just2so/hysteria2-ui/main/h-ui/uninstall_hui.sh)
```

# 微信推送

* 申请公众号测试账户,使用微信扫码即可 https://mp.weixin.qq.com/debug/cgi-bin/sandbox?t=sandbox/login
* 进入页面以后我们来获取到这四个值:appID appSecret openId template_id
* 模板内容

```bash
用户名：{{username.DATA}} 
```






