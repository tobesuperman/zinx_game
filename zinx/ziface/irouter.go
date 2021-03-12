package ziface

// IRouter 路由抽象接口，路由里的数据都是IRequest
type IRouter interface {
	PreHandle(request IRequest)  // 在处理conn业务之前的钩子方法
	Handle(request IRequest)     // 处理conn业务的钩子方法
	PostHandle(request IRequest) // 在处理conn业务之后的钩子方法
}
