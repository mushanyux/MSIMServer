package app

import (
	"errors"
	"net/http"

	"github.com/mushanyux/MSChatServerLib/config"
	"github.com/mushanyux/MSChatServerLib/pkg/mshttp"
)

type App struct {
	service IService
}

func New(ctx *config.Context) *App {
	return &App{
		service: NewService(ctx),
	}
}

func (a *App) Route(r *mshttp.MSHttp) {
	r.GET("/v1/apps/:app_id", a.get)
}

func (a *App) get(c *mshttp.Context) {
	appID := c.Param("app_id")
	resp, err := a.service.GetApp(appID)
	if err != nil {
		c.ResponseError(err)
		return
	}
	if resp.Status == StatusDisable {
		c.ResponseError(errors.New("app is disable"))
		return
	}
	c.JSON(http.StatusOK, &appResp{
		AppID:   resp.AppID,
		AppName: resp.AppName,
		AppLogo: resp.AppLogo,
	})
}

type appResp struct {
	AppID   string `json:"app_id,omitempty"`
	AppName string `json:"app_name,omitempty"`
	AppLogo string `json:"app_logo,omitempty"`
}
