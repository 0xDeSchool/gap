package app

import (
	"fmt"
	"reflect"
	"sync"
)

func newContainer() *Container {
	return &Container{
		initors:  sync.Map{},
		Services: make([]ServiceDescriptor, 0),
		values:   sync.Map{},
	}
}

type InitFunc func(container *Container, instance any)

type Container struct {
	initors  sync.Map
	Services []ServiceDescriptor
	values   sync.Map

	lock sync.Mutex
}

func (c *Container) Get(serviceType reflect.Type) interface{} {
	v, ok := c.GetOptional(serviceType)
	if !ok {
		panic(fmt.Sprintf("service %s not found", serviceType))
	}
	return v
}

// 获取对象，如果没有，返回nil
func (c *Container) GetOptional(serviceType reflect.Type) (interface{}, bool) {
	descriptor := c.firstOrDefault(serviceType)
	if descriptor == nil {
		return nil, false
	}
	return c.create(serviceType, descriptor), true
}

func (c *Container) GetArray(baseType reflect.Type) []interface{} {
	instances := make([]interface{}, 0)
	for i := 0; i < len(c.Services); i++ {
		v := &c.Services[i]
		implType := reflect.PtrTo(v.ServiceType)
		if v.ServiceType == baseType || implType == baseType || implType.AssignableTo(baseType) {
			instances = append(instances, c.create(v.ServiceType, v))
		}
	}
	return instances
}

func (c *Container) Add(descriptor ServiceDescriptor) {
	c.Services = append(c.Services, descriptor)
}

func (c *Container) TryAdd(descriptor ServiceDescriptor) {
	index := c.findIndex(&descriptor.ServiceType)
	if index < 0 {
		c.Services = append(c.Services, descriptor)
	}
}

func (c *Container) Configure(t reflect.Type, initFunc InitFunc) {
	v, ok := c.initors.Load(t)
	if ok {
		c.initors.Store(t, append(v.([]InitFunc), initFunc))
	} else {
		inits := []InitFunc{initFunc}
		c.initors.Store(t, inits)
	}
}

func (c *Container) firstOrDefault(serviceType reflect.Type) *ServiceDescriptor {
	for i := 0; i < len(c.Services); i++ {
		v := &c.Services[i]
		if v.ServiceType == serviceType {
			return v
		}
	}
	return nil
}

func (c *Container) findIndex(serviceType *reflect.Type) int {
	for i := 0; i < len(c.Services); i++ {
		v := &c.Services[i]
		if v.ServiceType == *serviceType {
			return i
		}
	}
	return -1
}

func (c *Container) create(serviceType reflect.Type, descriptor *ServiceDescriptor) interface{} {
	var instance any
	if descriptor.Scope == Singleton {
		v, ok := c.values.Load(serviceType)
		if !ok {
			v, _ = c.values.LoadOrStore(serviceType, c.createInstance(descriptor))
		}
		instance = v
		c.instanceInit(serviceType, instance, true)
	} else if descriptor.Scope == Transient {
		instance = c.createInstance(descriptor)
		c.instanceInit(serviceType, instance, false)
	} else {
		panic("invalid scope")
	}
	return instance
}

func (c *Container) createInstance(descriptor *ServiceDescriptor) interface{} {
	if descriptor.Value != nil {
		return descriptor.Value
	}
	if descriptor.Creator != nil {
		return descriptor.Creator(c)
	}
	return reflect.New(descriptor.ServiceType).Interface()
}

func (c *Container) instanceInit(serviceType reflect.Type, v any, delete bool) {
	if delete {
		c.lock.Lock() // lock for init, Warn: 当一个单例依赖另一个单例并且都需要配置时，可能会死锁
		defer c.lock.Unlock()
		if initor, ok := c.initors.LoadAndDelete(serviceType); ok {
			if init, ok := initor.([]InitFunc); ok {
				for i := 0; i < len(init); i++ {
					init[i](c, v)
				}
			}
		}
	} else {
		if initor, ok := c.initors.Load(serviceType); ok {
			if init, ok := initor.([]InitFunc); ok {
				for i := 0; i < len(init); i++ {
					init[i](c, v)
				}
			}
		}
	}

}
