package elastic

import (
	"github.com/gocraft/dbr/v2"
	"github.com/mushanyux/MSChatServerLib/pkg/db"
	"github.com/mushanyux/MSChatServerLib/pkg/util"
)

type DB struct {
	session *dbr.Session
}

func NewDB(session *dbr.Session) *DB {
	return &DB{
		session: session,
	}
}

func (d *DB) Insert(model *IndexerErrorModel) error {
	_, err := d.session.InsertInto("indexer_error").Columns(util.AttrToUnderscore(model)...).Record(model).Exec()
	return err
}

type IndexerErrorModel struct {
	Index      string
	Action     string
	DocumentID string
	Body       string
	Error      string
	db.BaseModel
}
