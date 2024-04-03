package session

import (
	"7days-golang-learn/ORM/day6-transaction/log"
	"reflect"
)

// Hook，翻译为钩子，其主要思想是提前在可能增加功能的地方埋好(预设)一个钩子.
// 当我们需要重新修改或者增加这个地方的逻辑的时候，把扩展的类或者方法挂载到这个点即可。

// 如前端常用的 hot reload 机制，前端代码发生变更时，自动编译打包，通知浏览器自动刷新页面，实现所写即所得。

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

type Account struct {
	ID       int `geeorm:"PRIMARY KEY"`
	Password string
}

func (account *Account) BeforeInsert(s *Session) error {
	log.Info("before inert", account)
	account.ID += 1000
	return nil
}

func (account *Account) AfterQuery(s *Session) error {
	log.Info("after query", account)
	account.Password = "******"
	return nil
}

func (s *Session) CallMethod(method string, value interface{}) {
	// hook方法注册在类的方法中，这里获取当前处理的对象：model/values的方法
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	// 获取对象
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		// 反射的方法调用，传入参数param为s的反射值：类型，指针，flag
		// 这里观察account类的方法定义，应该是⭐针对反射Call⭐的标准用法，不然⭐直接传入s⭐便可
		if v := fm.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return
}

// 使用实例
//func (s *Session) Find(values interface{}) error {
//  ⭐⭐⭐
//	s.CallMethod(BeforeQuery, nil)
//	// ...
//	for rows.Next() {
//		dest := reflect.New(destType).Elem()
//		// ...
//		s.CallMethod(AfterQuery, dest.Addr().Interface())
//		// ...
//	}
//	return rows.Close()
//}
