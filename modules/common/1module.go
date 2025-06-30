package common

import (
	"embed"

	"github.com/mushanyux/MSChatServerLib/config"
	"github.com/mushanyux/MSChatServerLib/pkg/register"
)

//go:embed sql
var sqlFS embed.FS

//go:embed swagger/api.yaml
var swaggerContent string

func init() {
	register.AddModule(func(ctx interface{}) register.Module {
		return register.Module{
			Name: "common",
			SetupAPI: func() register.APIRouter {
				return New(ctx.(*config.Context))
			},
			SQLDir:  register.NewSQLFS(sqlFS),
			Swagger: swaggerContent,
		}
	})

	register.AddModule(func(ctx interface{}) register.Module {
		return register.Module{
			SetupAPI: func() register.APIRouter {
				return NewManager(ctx.(*config.Context))
			},
		}
	})
}
