# chatgpt-wecom

企业微信接入 GPT3 接口。可以创建企业微信应用，体验 ChatGPT。

相较于官方提供的 `CreateCompletion` 接口，该项目增加了会话管理功能，能够较好地提供多轮对话能力。

## 快速安装

### 0. 前置步骤
* 登录企业微信的[管理后台](https://work.weixin.qq.com/wework_admin/loginpage_wx)创建机器人，名字随意取
  * 点击【应用管理】-【创建应用】，填写完毕保存。

### 1. 配置

修改 `conf/chatgpt.conf` 文件，主要涉及企业微信应用配置、GPT3 API Key、会话管理数据库配置等。

- 企业微信应用配置
  - corp_id：在企业微信后台【我的企业】-【企业信息】处获取【企业ID】
  - corp_secret：在企业微信后台【应用管理】处获取【Secret】
  - agent_id：在企业微信后台【应用管理】处获取【AgentId】
  - encoding_aes_key：企业微信后台 【接收消息】- 【API 接收消息】获取【EncodingAESKey】，可以随机生成
  - token：企业微信后台 【接收消息】- 【API 接收消息】获取【Token】，可以随机生成
- OpenAI Key
  - 需要自行申请

### 2. 运行
* **选择1：Docker运行（推荐）**

```shell
docker compose up -d
```

启动完毕，执行 `docker compose ps` 确认程序存活即可。进入步骤3 。

* 选择2：本地运行（需要手动配置MySQL）
  * [点击下载安装包](https://github.com/yijia2413/chatgpt-wecom/releases) 和 [配置文件](https://github.com/yijia2413/chatgpt-wecom/releases/download/v0.1.0/chatgpt.conf)
  * 修改`chatgpt.conf`, mysql 相关的配置
  * 执行 `./chatgpt-wecom -conf=chatgpt.conf -initdb`
  * 然后执行 `./chatgpt-wecom -conf=chatgpt.conf`
```
### 3. 配置企业微信
配置企业微信应用。在企业微信后台 【接收消息】- 【API 接收消息】配置接收消息服务器配置。

* URL 配置格式：`http://ip:port/wecom/receive`
* 在企业微信后台，添加可信IP地址

### 聊天
![img](/png/example.jpg)