package user

import (
	"github.com/gocraft/dbr/v2"
	"github.com/mushanyux/MSChatServerLib/config"
	"github.com/mushanyux/MSChatServerLib/pkg/db"
)

type deviceFlagDB struct {
	session *dbr.Session
	ctx     *config.Context
}

func newDeviceFlagDB(ctx *config.Context) *deviceFlagDB {
	return &deviceFlagDB{
		session: ctx.DB(),
		ctx:     ctx,
	}
}

func (d *deviceFlagDB) queryAll() ([]*deviceFlagModel, error) {
	var deviceFlags []*deviceFlagModel
	_, err := d.session.Select("*").From("device_flag").Load(&deviceFlags)
	return deviceFlags, err
}

type deviceFlagModel struct {
	DeviceFlag uint8
	Weight     int
	Remark     string
	db.BaseModel
}
