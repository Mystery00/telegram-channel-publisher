# telegram-channel-publisher

将Telegram频道的内容发布到其他地方，例如Halo的瞬间

## 使用方法

```shell
docker run -d \
    -e CONFIG_HOME=/app/etc \
    -v /path/to/config.yaml:/app/etc/config.yaml \
    mystery0/telegram-channel-publisher:20231030-0b40
```
> 备注：后面会把版本号改成数字的形式，而不是现在这个格式

环境变量是指定服务使用的配置文件的目录，配置文件的名称必须是`config.yaml`，如果运行时找不到配置文件，会报错退出。

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

创建令牌的时候，需要授予的权限如下：
- 附件查看：查询附件的信息以获取图片地址
- 附件管理：上传图片需要这个权限
- 瞬间管理：发布瞬间需要这个权限
- 瞬间查看：瞬间管理会默认授予这个权限

### 获取Halo的图片存储策略

![获取policy](img/get_policy.png)