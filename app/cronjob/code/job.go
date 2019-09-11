package code

const (
	JobAddFail = 10000 + iota
	JobMarshalFail
	PageError
	SyncError
	GetListError
	JobEditFail
	JobUpdateFail
	JobDeleteFail
)

// 错误信息
const (
	InsertSuccess  = "添加成功"
	SyncSuccess    = "更新成功"
	DeleteSuccess  = "删除成功"
	GetListSuccess = "获取列表成功"
)
