package config

const (
	EnvConfigHome = "CONFIG_HOME"
	EnvLogHome    = "LOG_HOME"
)

const (
	BotToken    = "bot.token"
	ApiEndpoint = "bot.endpoint"
	MediaDelay  = "bot.media.delay"

	ChannelEnable      = "bot.channel.enable"
	ChannelId          = "bot.channel.id"
	ChannelFilter      = "bot.channel.filter"
	ChannelReplyEnable = "bot.channel.reply.enable"
	ChannelReplyDelay  = "bot.channel.reply.delay"

	PrivateEnable      = "bot.private.enable"
	PrivateSender      = "bot.private.sender"
	PrivateReplyEnable = "bot.private.reply.enable"
	PrivateReplyDelay  = "bot.private.reply.delay"
)

const (
	LogLocal = "log.local"
	LogHome  = "log.home"
	LogFile  = "log.file"
	LogColor = "log.color"
	LogDebug = "log.debug"
)

const (
	PublisherType = "publisher.type"
)

const (
	HaloHost        = "halo.host"
	HaloToken       = "halo.token"
	HaloImageGroup  = "halo.image.group"
	HaloImagePolicy = "halo.image.policy"
)
