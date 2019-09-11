package code

const (
	none           = iota
	JobAddFail     = 10000
	JobMarshalFail = iota
	PageError
	SyncError
	GetListError
	JobEditFail
	JobUpdateFail
	JobDeleteFail
)

// 错误信息
const (
	InsertSuccess = "添加成功"
	SyncSuccess = "更新成功"
	DeleteSuccess = "删除成功"
)
