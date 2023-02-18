# chatgpt-wecom

企业微信接入 GPT3 接口。可以创建企业微信应用，体验 ChatGPT。

相较于官方提供的 `CreateCompletion` 接口，该项目增加了会话管理功能，能够较好地提供多轮对话能力。

## 快速开始

1. 修改配置

修改 `conf/online.conf` 文件，主要涉及企业微信应用配置、GPT3 API Key、会话管理数据库配置等。

- 企业微信应用配置
  - corp_id：在企业微信后台【我的企业】-【企业信息】处获取【企业ID】
  - corp_secret：在企业微信后台【应用管理】处获取【Secret】
  - agent_id：在企业微信后台【应用管理】处获取【AgentId】
  - encoding_aes_key：企业微信后台 【接收消息】- 【API 接收消息】获取【EncodingAESKey】，可以随机生成
  - token：企业微信后台 【接收消息】- 【API 接收消息】获取【Token】，可以随机生成
- Open AI Key
  - 需要自行申请
- 数据库
  - ~数据库需要自行创建，数据表的创建可以通过命令行方式执行。~
  - 数据库支持 sqlite3，可以通过修改配置使用。如果使用 MySQL，需要自行创建数据库。
  - **数据表在程序启动时自动创建。**

2. `Docker` 运行

```shell
docker-compose up -d
```

3. ~初始化数据表~

数据表在程序启动时自动创建。


4. 配置企业微信应用。在企业微信后台 【接收消息】- 【API 接收消息】配置接收消息服务器配置。

- URL 配置格式：`http[s]://ip:port/wecom/receive`


## Changelog

### v0.1.1

- 修复 prompt 过长导致接口调用失败问题
- 支持 sqlite3
- 自动初始化数据库
- 支持企业微信进入事件，可以配置进入事件回复语
- 支持关闭会话功能

### v0.1.0

- 项目初始化