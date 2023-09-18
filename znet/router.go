package znet

import "myzinx/ziface"

// 实现router时，先嵌入这个BaseRouter基类， 然后根据需要对这个基类的方法进行重写
type BaseRouter struct {
}

/*
	这里之所以BaseRouter的方法都是空的
	是因为由的Router不希望有PreHandle、PostHandle这两个业务
	所以Router全部继承BaseRouter的好处是，不需要实现PreHandle、PostHandle
*/

func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

func (br *BaseRouter) Handle(request ziface.IRequest) {}

func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
