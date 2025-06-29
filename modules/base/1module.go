package base

import (
	"embed"

	"github.com/mushanyux/MSChatServerLib/config"
	"github.com/mushanyux/MSChatServerLib/pkg/register"
	"github.com/mushanyux/MSIMServer/modules/base/app"
)

//go:embed sql
var sqlFS embed.FS

func init() {

	register.AddModule(func(ctx interface{}) register.Module {

		return register.Module{
			SetupAPI: func() register.APIRouter {
				return app.New(ctx.(*config.Context))
			},
			SQLDir: register.NewSQLFS(sqlFS),
		}
	})
}
