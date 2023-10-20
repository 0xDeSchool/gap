package app

import (
	"fmt"
	"reflect"
	"sync"
)

func newContainer() *Container {
	return &Container{
		Services: make([]*serviceItem, 0),
		inits:    make(map[reflect.Type][]InitFunc),
	}
}

type serviceItem struct {
	*ServiceDescriptor
	Inits []InitFunc
	lock  sync.Mutex
}

func (item *serviceItem) GetValue(c *Container) any {
	if item.Value != nil {
		return item.Value
	}
	var v any
	if item.Scope == Transient {
		v = item.createInstance(c)
		for i := 0; i < len(item.Inits); i++ {
			item.Inits[i](c, v)
		}
	} else if item.Scope == Singleton {
		item.lock.Lock()
		defer item.lock.Unlock()
		v = item.createInstance(c)
		for i := 0; i < len(item.Inits); i++ {
			item.Inits[i](c, v)
		}
		item.Value = v
	} else {
		panic("not supported scope")
	}
	return v
}

func (item *serviceItem) createInstance(c *Container) interface{} {
	if item.Value != nil {
		return item.Value
	}
	if item.Creator != nil {
		return item.Creator(c)
	}
	return reflect.New(item.ServiceType).Interface()
}

type InitFunc func(container *Container, instance any)

type Container struct {
	Services []*serviceItem

	inits map[reflect.Type][]InitFunc
}

func (c *Container) Get(serviceType reflect.Type) interface{} {
	v, ok := c.GetOptional(serviceType)
	if !ok {
		panic(fmt.Sprintf("service %s not found", serviceType))
	}
	return v
}

// GetOptional 获取对象，如果没有，返回nil
func (c *Container) GetOptional(serviceType reflect.Type) (interface{}, bool) {
	descriptor := c.firstOrDefault(serviceType)
	if descriptor == nil {
		return nil, false
	}
	return descriptor.GetValue(c), true
}

// GetArray 有问题，只能transient有效，可使用 options 模式来替代，通过 ConfigureOption 来配置
func (c *Container) GetArray(baseType reflect.Type) []interface{} {
	instances := make([]interface{}, 0)
	for i := 0; i < len(c.Services); i++ {
		v := c.Services[i]
		implType := reflect.PtrTo(v.ServiceType)
		if v.ServiceType == baseType || implType == baseType || implType.AssignableTo(baseType) {
			instances = append(instances, v.GetValue(c))
		}
	}
	return instances
}

func (c *Container) Add(descriptor *ServiceDescriptor) {
	item := &serviceItem{
		ServiceDescriptor: descriptor,
	}
	if its, has := c.inits[descriptor.ServiceType]; has {
		item.Inits = its
		delete(c.inits, descriptor.ServiceType)
	}
	c.Services = append(c.Services, item)
}

func (c *Container) TryAdd(descriptor *ServiceDescriptor) {
	index := c.findIndex(descriptor.ServiceType)
	if index < 0 {
		c.Add(descriptor)
	}
}

// Configure 配置服务，不支持两个类型相互依赖
func (c *Container) Configure(t reflect.Type, initFunc InitFunc) {
	item := c.firstOrDefault(t)
	if item == nil {
		c.inits[t] = append(c.inits[t], initFunc)
	} else {
		item.Inits = append(item.Inits, initFunc)
	}
}

func (c *Container) firstOrDefault(serviceType reflect.Type) *serviceItem {
	for i := 0; i < len(c.Services); i++ {
		v := c.Services[i]
		if v.ServiceType == serviceType {
			return v
		}
	}
	return nil
}

func (c *Container) findIndex(serviceType reflect.Type) int {
	for i := 0; i < len(c.Services); i++ {
		if c.Services[i].ServiceType == serviceType {
			return i
		}
	}
	return -1
}
