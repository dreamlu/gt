package result

import (
	"encoding/json"
	"errors"
	"fmt"
	errors2 "github.com/dreamlu/gt/tool/type/errors"
)

// status and msg
const (
	Status = "status"
	Msg    = "msg"
)

// 约定状态码
// 或 通过GetMapData()自定义
const (
	CodeSuccess   = 200 // 请求成功
	CodeNoAuth    = 203 // 请求非法
	CodeNoResult  = 204 // 暂无数据
	CodeValidator = 210 // 字段验证
	CodeText      = 271 // 全局文字提示
	CodeError     = 500 // 系统繁忙
)

// 约定提示信息
const (
	MsgSuccess  = "请求成功"
	MsgNoAuth   = "请求非法"
	MsgNoResult = "暂无数据"
)

// 约定提示信息
var (
	MapSuccess  = GetMapData(CodeSuccess, MsgSuccess)   // 请求成功
	MapNoResult = GetMapData(CodeNoResult, MsgNoResult) // 暂无数据
	MapNoAuth   = GetMapData(CodeNoAuth, MsgNoAuth)     // 请求非法
)

// 分页数据信息
type GetInfoPager struct {
	*GetInfo
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
	*MapData
	Data interface{} `json:"data"` // 数据存储
}

// 信息,通用
type MapData struct {
	Status int64       `json:"status"`
	Msg    interface{} `json:"msg"`
}

// 信息通用,状态码及信息提示
func GetMapData(status int64, msg interface{}) *MapData {

	return &MapData{
		Status: status,
		Msg:    msg,
	}
}

// text
func GetText(Msg interface{}) *MapData {

	return GetMapData(CodeText, Msg)
}

// 信息成功通用(成功通用, 无分页)
func GetSuccess(data interface{}) *GetInfo {

	return &GetInfo{
		MapData: MapSuccess,
		Data:    data,
	}
}

// 信息分页通用(成功通用, 分页)
func GetSuccessPager(data interface{}, pager Pager) *GetInfoPager {

	return &GetInfoPager{
		GetInfo: GetSuccess(data),
		Pager:   pager,
	}
}

// 信息失败通用
func GetError(msg interface{}) *MapData {

	return &MapData{
		Status: CodeError,
		Msg:    msg,
	}
}

// 无分页通用
func GetData(data interface{}, mapData *MapData) *GetInfo {

	return &GetInfo{
		MapData: mapData,
		Data:    data,
	}
}

// 分页通用
func GetDataPager(data interface{}, mapData *MapData, pager Pager) *GetInfoPager {

	return &GetInfoPager{
		GetInfo: GetData(data, mapData),
		Pager:   pager,
	}
}

// 转化
func (m *GetInfo) Parent() *MapData {

	return m.MapData
}

func (m *GetInfoPager) Parent() *GetInfo {

	return &GetInfo{
		MapData: m.MapData,
		Data:    m.Data,
	}
}

// string
func (m MapData) String() string {

	return StructToString(m)
}

func (m GetInfo) String() string {

	return StructToString(m)
}

func (m GetInfoPager) String() string {

	return StructToString(m)
}

func (m *GetInfoPager) Add(key string, value interface{}) (rmp ResultMap) {

	return m.Parent().Add("pager", m.Pager).Add(key, value)
}

func (m *GetInfo) Add(key string, value interface{}) (rmp ResultMap) {

	return m.Parent().Add("data", m.Data).Add(key, value)
}

func (m *MapData) Add(key string, value interface{}) (rmp ResultMap) {

	return rmp.Add(Status, m.Status).Add(Msg, m.Msg).Add(key, value)
}

func (m *GetInfoPager) AddStruct(value interface{}) (rmp ResultMap) {

	return m.Parent().Add("pager", m.Pager).AddStruct(value)
}

func (m *GetInfo) AddStruct(value interface{}) (rmp ResultMap) {

	return m.Parent().Add("data", m.Data).AddStruct(value)
}

func (m *MapData) AddStruct(value interface{}) (rmp ResultMap) {

	return rmp.Add(Status, m.Status).Add(Msg, m.Msg).AddStruct(value)
}

func StructToString(st interface{}) string {
	s, err := json.Marshal(st)
	if err != nil {
		return ""
	}
	return string(s)
}

func StringToStruct(str string, st interface{}) error {
	return json.Unmarshal([]byte(str), st)
}

// error sugar
// result sugar
func CError(err error) *MapData {
	if errors.As(err, &errors2.TextErr) {
		if err.Error() == MsgNoResult {
			return MapNoResult
		}
		return GetText(err.Error())
	}
	return GetError(err.Error())
}

func TextError(msg string) error {
	return fmt.Errorf("%w", &errors2.TextError{Msg: msg})
}
