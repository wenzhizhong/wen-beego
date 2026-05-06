package constant

type MqNameType string

const (
	MQ_TEST_MSG           MqNameType = "Test.ActionTestMsg"    // 测试消息
	MQ_API_LOG_SAVE_TO_DB MqNameType = "ApiLog.ActionSaveToDb" // 接口日志保存到数据库
)
