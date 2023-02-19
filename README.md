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
- 数据库
  - ~数据库需要自行创建，数据表的创建可以通过命令行方式执行。~
  - 数据库支持 sqlite3，可以通过修改配置使用。如果使用 MySQL，需要自行创建数据库。
  - **数据表在程序启动时自动创建。**

### 2. 运行
* **选择1：Docker运行（推荐）**

```shell
git clone https://github.com/yijia2413/chatgpt-wecom.git
cd chatgpt-wecom
# 构建镜像
make dockerenv
# 运行带sqlite的镜像，运行前确认chatgpt.conf修改完毕
docker run -it -d --name chatgpt --restart=always \
  -v $(pwd)/conf/chatgpt.conf:/home/works/program/chatgpt.conf chatgpt-wecom:0.1.1
```

* **选择2：本地运行**
* 下载对应的二进制，[chatgpt-wecom](https://github.com/yijia2413/chatgpt-wecom/releases)
* 执行命令 `./chatgpt-wecom -conf=conf/chatgpt.conf` 即可，同理需要确认`chatgpt.conf`配置完毕

### 3. 配置企业微信

* URL 配置格式：`http://ip:port/wecom/receive`
* 在企业微信后台，添加可信IP地址

## FAQ

**怎么创建数据库**

- v0.1.1 版本中支持 sqlite3 数据库，只需要修改配置文件的配置，程序启动后便会初始化数据库和数据表，不需要额外的操作。

- 如果使用的是 MySQL，则需要自行创建数据库，建库 SQL 可以参考下面的命令：[init.sql](/init.sql)


**数据库连接失败**

- 首先检查数据库配置是否正确
- 如果使用 docker 部署服务，需要确认容器内能否连接到数据库。最常见的一个问题是，在宿主机部署了 MySQL，但是在容器内配置 `127.0.0.1`，这种情况需要配置宿主机的 IP

**数据库配置说明**

v0.1.1 版本可以支持 MySQL、SQLite、PostgreSQL。常见的配置如下：

MySQL:

```toml
[database]
# mysql
driver="mysql"
dataSource="root:12345678@tcp(127.0.0.1:3306)/chatgpt?parseTime=True"
```

SQLite

```toml
[database]
# sqlite3
driver="sqlite3"
dataSource="file:chatgpt?_fk=1&parseTime=True"
```

## Changelog

### v0.1.1

- 修复 prompt 过长导致接口调用失败问题
- 支持 sqlite3
- 自动初始化数据库
- 支持企业微信进入事件，可以配置进入事件回复语。需要先在企业微信上配置【上报进入应用事件】
- 支持关闭会话功能

### v0.1.0

- 项目初始化

### Happy Chatting
![img](/png/example.jpg)