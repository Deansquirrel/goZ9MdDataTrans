package object

type TaskKey string

//任务Key
const (
	TaskKeyRefreshHeartBeat TaskKey = "RefreshHeartBeat" //刷新心跳
	TaskKeyRefreshConfig    TaskKey = "RefreshConfig"    //刷新配置

	TaskKeyRefreshZxKc     TaskKey = "RefreshZxKc"     //刷新最新库存变动
	TaskKeyRefreshMdYyInfo TaskKey = "RefreshMdYyInfo" //刷新门店营业信息

	TaskKeyRestoreZxKc     TaskKey = "RestoreZxKc"
	TaskKeyRestoreMdYyInfo TaskKey = "RestoreMdYyInfo"
)

//运行模式
type RunMode string

const (
	RunModeMdCollect RunMode = "MdCollect" //门店采集
	RunModeBbRestore RunMode = "BbRestore" //报表恢复
)
