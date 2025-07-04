package event

import (
	"fmt"

	"github.com/gocraft/dbr/v2"
	"github.com/mushanyux/MSChatServerLib/config"
	"github.com/mushanyux/MSChatServerLib/pkg/log"
	"github.com/mushanyux/MSChatServerLib/pkg/util"
	"github.com/mushanyux/MSIMServer/modules/file"
	"go.uber.org/zap"

	et "github.com/mushanyux/MSChatServerLib/pkg/msevent"
)

const (
	// GroupCreate 群创建
	GroupCreate string = "group.create"
	// GroupUnableAddDestroyAccount 无法添加注销账号到群聊
	GroupUnableAddDestroyAccount string = "group.unable.add.destroy.account"
	// GroupUpdate 群更新
	GroupUpdate string = "group.update"
	// GroupMemberAdd 群成员添加
	GroupMemberAdd string = "group.memberadd"
	// GroupMemberScanJoin 扫码加入群
	GroupMemberScanJoin string = "group.member.scan.join"
	// GroupMemberTransferGrouper 转让群主
	GroupMemberTransferGrouper string = "group.member.transfer.grouper"
	// GroupAvatarUpdate 群头像更新
	GroupAvatarUpdate string = "group.avatar.update"
	// GroupMemberRemove 群成员移除
	GroupMemberRemove string = "group.memberremove"
	// GroupDisband 群解散
	GroupDisband string = "group.disband"
	// FriendApply 好友申请
	FriendApply string = "friend.apply"
	// GroupMemberInviteRequest 群邀请请求
	GroupMemberInviteRequest string = "group.member.invite"
	// ConversationDelete 删除最近会话
	ConversationDelete string = "conversation.delete"
	// EventUserRegister 用户注册
	EventUserRegister string = "user.register"
	// EventUserPublishMoment 用户发布动态
	EventUserPublishMoment string = "moment.publish"
	// EventUserDeleteMoment 用户删除动态
	EventUserDeleteMoment string = "moment.delete"
	// FriendSure 好友确认
	FriendSure string = "friend.sure"
	// FriendDelete 好友删除
	FriendDelete string = "friend.delete"
	// OrgOrDeptCreate 组织货部门创建
	OrgOrDeptCreate string = "organization_department.create"
	// OrgOrDeptEmployeeUpdate 组织或部门员工更改
	OrgOrDeptEmployeeUpdate string = "organization_department.employee.update"
	// OrgEmployeeExit 组织成员退出
	OrgEmployeeExit string = "organization.employee.exit"
	// EventUpdateSearchMessage 修改搜索消息内容
	EventUpdateSearchMessage string = "message.update.search.data"
)

// Event 事件
type Event struct {
	db  *DB
	ctx *config.Context
	log.Log
	fileService file.IService
}

// New 创建一个事件
func New(ctx *config.Context) *Event {
	e := &Event{
		ctx:         ctx,
		db:          NewDB(ctx.DB()),
		Log:         log.NewTLog("Event"),
		fileService: file.NewService(ctx),
	}
	e.registerHandlers()
	return e
}

// Begin 开启事件
func (e *Event) Begin(data *et.Data, tx *dbr.Tx) (int64, error) {
	eventID, err := e.db.InsertTx(&Model{
		Event: data.Event,
		Type:  data.Type.Int(),
		Data:  util.ToJson(data.Data),
	}, tx)
	return eventID, err
}

// Commit 提交事件
func (e *Event) Commit(eventID int64) {
	eventModel, err := e.db.QueryWithID(eventID)
	if err != nil {
		e.Error("查询事件失败！", zap.Error(err), zap.Int64("eventID", eventID))
		return
	}
	e.handleEvent(eventModel)
}

// Support 是否支持的事件类型
func (e *Event) Support(typ int) bool {
	switch typ {
	case et.Message.Int():
		return true
	}
	return false
}

func (e *Event) updateEventStatus(err error, versionLock int64, eventID int64) {
	var reason string
	var status = et.Success
	if err != nil {
		e.Warn("执行事件失败！", zap.Error(err), zap.Int64("eventID", eventID))
		reason = fmt.Sprintf("执行事件失败！-> %v", err)
		status = et.Fail
	}
	err = e.db.UpdateStatus(reason, status.Int(), versionLock, eventID)
	if err != nil {
		e.Error("更新事件状态失败！", zap.Int64("eventID", eventID), zap.Error(err))
		return
	}
}

// EventTimerPush 定时发布事件
func (e *Event) EventTimerPush() {
	models, err := e.db.QueryAllWait(1000)
	if err != nil {
		e.Error("查询所有待发布的事件失败！", zap.Error(err))
		return
	}
	if len(models) > 0 {
		for _, model := range models {
			e.handleEvent(model)
		}
	}
}
