package register

import "github.com/mushanyux/MSIMServer/pkg/mshttp"

// APIRouter api路由者
type APIRouter interface {
	Route(r *mshttp.MSHttp)
}

var apiRoutes = make([]APIRouter, 0)

// Add 添加api
func Add(r APIRouter) {
	apiRoutes = append(apiRoutes, r)
}

var taskRoutes = make([]TaskRouter, 0)

// GetRoutes 获取所有路由者
func GetRoutes() []APIRouter {
	return apiRoutes
}

// TaskRouter task路由者
type TaskRouter interface {
	RegisterTasks()
}

// AddTask 添加任务
func AddTask(task TaskRouter) {
	taskRoutes = append(taskRoutes, task)
}

// GetTasks 获取所有任务
func GetTasks() []TaskRouter {
	return taskRoutes
}
