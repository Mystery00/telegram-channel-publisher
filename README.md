# telegram-channel-publisher

将Telegram频道的内容发布到其他地方，例如Halo的瞬间

## 使用方法

TODO

## 配置文件

```yaml
bot:
  # 机器人的token
  token: "{{telegram.bot.token}}"
  # 机器人的接入点信息，如果部署的服务器无法访问Telegram，可以使用代理的形式来访问
  endpoint: ""

log:
  home: "logs"
  file: "publisher.log"
  color: false
  local: false
  # 调试模式
  debug: false

publisher:
  # 发布的类型，目前只支持log和halo，log是将内容输出到日志文件，halo是将内容发布到Halo的瞬间
  type: "halo"

halo:
  # Halo的访问域名，例如 https://blog.mystery0.vip
  host: "{{halo.host}}"
  # Halo的token 参考下图进行创建
  token: "{{halo.token}}"
  image:
    # 消息中的图片的分组信息，最好是在管理后台自己创建一个分组专门放瞬间的图片，方便管理，可为空
    group: ""
    # 消息中的图片的存储策略，参考后续内容进行获取
    policy: "{{halo.image.policy}}"
```

### 创建Halo的token
![创建halotoken](img/create-halo-token.png)

### 获取Halo的图片存储策略
TODO