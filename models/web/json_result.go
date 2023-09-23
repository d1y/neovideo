package web

import "github.com/kataras/iris/v12"

// copy by https://github.com/mlogclub/simple/blob/master/web/json_result.go

const (
	ActionSuccess = "获取成功"
	ActionFail    = "获取失败"
)

type JsonResult[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
	Success bool   `json:"success"`
}

func (jr *JsonResult[T]) SetMessage(msg string) *JsonResult[T] {
	jr.Message = msg
	return jr
}

func (jr *JsonResult[T]) SetData(data T) *JsonResult[T] {
	jr.Message = ActionFail
	jr.Data = data
	return jr
}

func (jr *JsonResult[T]) SetActionSuccess() *JsonResult[T] {
	jr.Message = ActionSuccess
	jr.Success = true
	return jr
}

func (jr *JsonResult[T]) SetActionFail() *JsonResult[T] {
	jr.Message = ActionFail
	jr.Success = false
	return jr
}

func (jr *JsonResult[T]) SetSuccessWithBool(success bool) *JsonResult[T] {
	jr.Success = success
	return jr
}

func (jr *JsonResult[T]) Build(ctx iris.Context) {
	ctx.JSON(jr)
}

func NewJSONResultWithMessage(msg string) *JsonResult[any] {
	return &JsonResult[any]{
		Message: msg,
	}
}

func NewJSONResultWithSuccess[T any](data T) *JsonResult[T] {
	return &JsonResult[T]{
		Success: true,
		Message: ActionSuccess,
		Data:    data,
	}
}

func NewJSONResultWithError(err error) *JsonResult[any] {
	return &JsonResult[any]{
		Message: err.Error(),
		Success: false,
	}
}
