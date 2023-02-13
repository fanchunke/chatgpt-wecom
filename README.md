
# chatgpt-wecom

企业微信接入 GPT3 接口。可以创建企业微信应用，体验 ChatGPT。

相较于官方提供的 `CreateCompletion` 接口，该项目增加了会话管理功能，能够较好地提供多轮对话能力。

## 快速开始

1. 修改配置

修改 `conf/online.conf` 文件，主要涉及企业微信应用配置、GPT3 API Key、会话管理数据库配置等。

数据库需要自行创建，数据表的创建可以通过命令行方式执行。

2. `Docker` 运行

```shell
docker-compose up -d
```

3. 初始化数据表

```shell
# 进入容器
docker exec -it chatgpt-wecom bash

# 执行命令
./app -conf=conf/online.conf -init-ent
```

4. 配置企业微信应用后，即可体验。