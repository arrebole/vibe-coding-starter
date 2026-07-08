package request

type ListTodosRequest struct {
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

type CreateTodoRequest struct {
	Title string `json:"title" binding:"required,min=1,max=200"`
}

type UpdateTodoRequest struct {
	Title     *string `json:"title" binding:"omitempty,min=1,max=200"`
	Completed *bool   `json:"completed"`
}
