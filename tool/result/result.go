package result

import "encoding/json"

// status and msg
const (
	Status = "status"
	Msg    = "msg"
)

// 约定状态码
// 或 通过GetMapData()自定义
const (
	CodeSuccess    = 200 // 请求成功
	CodeCreate     = 201 // 创建成功
	CodeNoAuth     = 203 // 请求非法
	CodeNoResult   = 204 // 暂无数据
	CodeUpdate     = 206 // 修改成功
	CodeDelete     = 209 // 删除成功
	CodeValidator  = 210 // 字段验证
	CodeCount      = 211 // 账号相关
	CodeValSuccess = 217 // 验证成功
	CodeValError   = 218 // 验证失败
	CodeExistOrNo  = 220 // 数据无变化
	CodeSQL        = 222 // 数据库相关
	CodeText       = 271 // 全局文字提示
	CodeError      = 500 // 系统繁忙
)

// 约定提示信息
const (
	MsgSuccess    = "请求成功"
	MsgCreate     = "创建成功"
	MsgNoAuth     = "请求非法"
	MsgNoResult   = "暂无数据"
	MsgDelete     = "删除成功"
	MsgUpdate     = "修改成功"
	MsgError      = "未知错误"
	MsgExistOrNo  = "数据无变化"
	MsgCountErr   = "用户账号或密码错误"
	MsgNoCount    = "用户账号不存在"
	MsgValSuccess = "验证成功"
	MsgValError   = "验证失败"
)

// 约定提示信息
var (
	MapSuccess    = GetMapData(CodeSuccess, MsgSuccess)       // 请求成功
	MapError      = GetMapData(CodeError, MsgError)           // 通用失败
	MapUpdate     = GetMapData(CodeUpdate, MsgUpdate)         // 修改成功
	MapDelete     = GetMapData(CodeDelete, MsgDelete)         // 删除成功
	MapCreate     = GetMapData(CodeCreate, MsgCreate)         // 创建成功
	MapNoResult   = GetMapData(CodeNoResult, MsgNoResult)     // 暂无数据
	MapNoAuth     = GetMapData(CodeNoAuth, MsgNoAuth)         // 请求非法
	MapExistOrNo  = GetMapData(CodeExistOrNo, MsgExistOrNo)   // 指数据修改没有变化 或者 给的条件值不存在
	MapCountErr   = GetMapData(CodeCount, MsgCountErr)        // 用户账号密码错误
	MapNoCount    = GetMapData(CodeCount, MsgNoCount)         // 用户账号不存在
	MapValSuccess = GetMapData(CodeValSuccess, MsgValSuccess) // 验证成功
	MapValError   = GetMapData(CodeValError, MsgValError)     // 验证失败
)

// 分页数据信息
type GetInfoPager struct {
	GetInfo
	Pager Pager `json:"pager"`
}

// pager info
type Pager struct {
	ClientPage int64 `json:"client_page"` // 当前页码
	EveryPage  int64 `json:"every_page"`  // 每一页显示的数量
	TotalNum   int64 `json:"total_num"`   // 数据总数量
}

// 无分页数据信息
// 分页数据信息
type GetInfo struct {
	MapData
	Data interface{} `json:"data"` // 数据存储
}

// 信息,通用
type MapData struct {
	Status int64       `json:"status"`
	Msg    interface{} `json:"msg"`
}

// 信息通用,状态码及信息提示
func GetMapData(status int64, msg interface{}) MapData {

	return MapData{
		Status: status,
		Msg:    msg,
	}
}

// 后端提示
func GetText(Msg interface{}) MapData {

	return GetMapData(CodeText, Msg)
}

// 信息成功通用(成功通用, 无分页)
func GetSuccess(data interface{}) GetInfo {

	return GetInfo{
		MapData: MapSuccess,
		Data:    data,
	}
}

// 信息分页通用(成功通用, 分页)
func GetSuccessPager(data interface{}, pager Pager) GetInfoPager {

	return GetInfoPager{
		GetInfo: GetSuccess(data),
		Pager:   pager,
	}
}

// 信息失败通用
func GetError(data interface{}) GetInfo {

	return GetInfo{
		MapData: MapError,
		Data:    data,
	}
}

// 无分页通用
func GetData(data interface{}, mapData MapData) GetInfo {

	return GetInfo{
		MapData: mapData,
		Data:    data,
	}
}

// 分页通用
func GetDataPager(data interface{}, mapData MapData, pager Pager) GetInfoPager {

	return GetInfoPager{
		GetInfo: GetData(data, mapData),
		Pager:   pager,
	}
}

// string
func (m *MapData) String() (string, error) {

	return structToString(m)
}

func (m *GetInfo) String() (string, error) {

	return structToString(m)
}

func (m *GetInfoPager) String() (string, error) {

	return structToString(m)
}

func structToString(st interface{}) (string, error) {
	s, err := json.Marshal(st)
	if err != nil {
		return "", err
	}
	return string(s), nil
}
