package app

import (
	"reflect"
)

type serverCreator[T any] func() *T

var container = newContainer()

func GetContainer() *Container {
	return container
}

func CreateContainer() *Container {
	return newContainer()
}

func AddService(descriptor ServiceDescriptor) *Container {
	container.Add(descriptor)
	return container
}

func TryAddService(descriptor ServiceDescriptor) *Container {
	container.TryAdd(descriptor)
	return container
}

func Configure[T any](f func(container *Container, v *T)) {
	container.Configure(getServiceType[T](), func(container *Container, instance any) {
		f(container, instance.(*T))
	})
}

func Add[T interface{}](scope ServiceScope, creator serverCreator[T]) {
	AddService(ServiceDescriptor{
		ServiceType: getServiceType[T](),
		Creator:     convertCreator(creator),
		Scope:       scope,
	})
}

func TryAdd[T interface{}](scope ServiceScope, creator serverCreator[T]) {
	TryAddService(ServiceDescriptor{
		ServiceType: getServiceType[T](),
		Creator:     convertCreator(creator),
		Scope:       scope,
	})
}

func AddTransient[T interface{}](creator serverCreator[T]) {
	Add(Transient, creator)
}

// AddTransientDefault 添加Transient无参构造对象
func AddTransientDefault[T interface{}]() {
	Add(Transient, func() *T {
		return new(T)
	})
}

func TryAddTransient[T interface{}](creator serverCreator[T]) {
	TryAdd(Transient, creator)
}

// TryAddTransientDefault 添加Transient无参构造对象
func TryAddTransientDefault[T interface{}]() {
	TryAdd(Transient, func() *T {
		return new(T)
	})
}

func AddSingleton[T interface{}](creator serverCreator[T]) {
	Add(Singleton, creator)
}

// AddSingletonDefault 添加 Singleton 无参构造对象
func AddSingletonDefault[T interface{}]() {
	Add(Singleton, func() *T {
		return new(T)
	})
}

func TryAddSingleton[T interface{}](creator serverCreator[T]) {
	TryAdd(Singleton, creator)
}

// TryAddSingletonDefault 添加 Singleton 无参构造对象
func TryAddSingletonDefault[T interface{}]() {
	TryAdd(Singleton, func() *T {
		return new(T)
	})
}

func AddValue(value interface{}) {
	AddService(ServiceDescriptor{
		ServiceType: getInterfaceType(value),
		Value:       value,
		Scope:       Singleton,
	})
}

func TryAddValue(value interface{}) {
	TryAddService(ServiceDescriptor{
		ServiceType: getInterfaceType(value),
		Value:       value,
		Scope:       Singleton,
	})
}

func AddByType[T interface{}](serviceType reflect.Type, scope ServiceScope, creator func(*Container) *T) {
	AddService(ServiceDescriptor{
		ServiceType: serviceType,
		Creator:     func(c *Container) interface{} { return creator(c) },
		Scope:       scope,
	})
}

func TryAddByType[T interface{}](serviceType reflect.Type, scope ServiceScope, creator func(*Container) *T) {
	TryAddService(ServiceDescriptor{
		ServiceType: serviceType,
		Creator:     func(c *Container) interface{} { return creator(c) },
		Scope:       scope,
	})
}

func Get[T interface{}]() *T {
	return container.Get(getServiceType[T]()).(*T)
}

func GetPtr[T interface{}]() T {
	return *container.Get(getServiceType[T]()).(*T)
}

func GetOptional[T interface{}]() (*T, bool) {
	v, ok := container.GetOptional(getServiceType[T]())
	if !ok {
		return nil, ok
	}
	return v.(*T), true
}

func GetPtrOptional[T interface{}]() (T, bool) {
	v, ok := container.GetOptional(getServiceType[T]())
	if !ok {
		var defaultValue T
		return defaultValue, ok
	}
	v, ok = v.(T)
	return v, ok
}

func GetByType[T interface{}](serviceType reflect.Type) *T {
	return container.Get(serviceType).(*T)
}

func GetArray[T interface{}]() []T {
	items := container.GetArray(getServiceType[T]())
	typeItems := make([]T, len(items))
	for i := 0; i < len(items); i++ {
		typeItems[i] = items[i].(T)
	}
	return typeItems
}

func getServiceType[T interface{}]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func getInterfaceType(value interface{}) reflect.Type {
	return reflect.TypeOf(value).Elem()
}

func convertCreator[T interface{}](creator serverCreator[T]) ServiceCreator {
	if creator == nil {
		return nil
	}
	return func(c *Container) interface{} { return creator() }
}
