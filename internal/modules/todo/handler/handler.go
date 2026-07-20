package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/dto"
	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/request"
	"github.com/arrebole/vibe-coding-starter/internal/modules/todo/service"
	"github.com/arrebole/vibe-coding-starter/internal/response"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{service: service}
}

// Register 注册 todo 一级路由下的所有 HTTP 接口。
func (h *Handler) Register(group *gin.RouterGroup) {
	group.GET("/todos", h.List)
	group.POST("/todos", h.Create)
	group.PATCH("/todos/:id", h.Update)
	group.DELETE("/todos/:id", h.Delete)
}

func (h *Handler) List(c *gin.Context) {
	var req request.ListTodosRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	todos, err := h.service.List(c.Request.Context(), req.Limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "查询 Todo 列表失败")
		return
	}
	response.OK(c, todos)
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	todo, err := h.service.Create(c.Request.Context(), dto.CreateTodoInput{Title: req.Title})
	if errors.Is(err, service.ErrInvalidTitle) {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建 Todo 失败")
		return
	}
	response.Created(c, todo)
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req request.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	todo, err := h.service.Update(c.Request.Context(), dto.UpdateTodoInput{
		ID:        id,
		Title:     req.Title,
		Completed: req.Completed,
	})
	if errors.Is(err, service.ErrNotFound) {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	if errors.Is(err, service.ErrInvalidTitle) {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "更新 Todo 失败")
		return
	}
	response.OK(c, todo)
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	err := h.service.Delete(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "删除 Todo 失败")
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, "id 必须是正整数")
		return 0, false
	}
	return uint(id), true
}
