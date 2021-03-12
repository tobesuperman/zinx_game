package znet

import "zinx/ziface"

// BaseRouter 实现IRouter时，先嵌入这个基类，再根据需要对这个基类的方法进行重写
type BaseRouter struct {
}

/*
这里之所以BaseRouter的方法都为空
是因为有的Router不希望由PreHandle，PostHandle这两个方法（业务）
所有Router全部继承BaseRouter的好处就是，不需要实现PreHandler和PostHandler
*/

func (br *BaseRouter) PreHandle(request ziface.IRequest) {

}
func (br *BaseRouter) Handle(request ziface.IRequest) {

}
func (br *BaseRouter) PostHandle(request ziface.IRequest) {

}
