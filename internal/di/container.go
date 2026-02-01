package di

import (
	"context"
	"log"
	"sync"
)

type Scope int

const (
	Singleton Scope = iota + 1
	Scoped
)

type contextKey int

const containerKey contextKey = 1

type DepFactoryFunc func(c Container) (any, error)

type depInfo struct {
	key     string
	scope   Scope
	factory DepFactoryFunc
}

type Container interface {
	AddSingleton(key string, fn DepFactoryFunc)
	AddScoped(key string, fn DepFactoryFunc)
	Scoped(ctx context.Context) context.Context
	Get(key string) any
}

type container struct {
	parent  *container
	deps    map[string]depInfo
	vals    map[string]any
	tracked tracked
	mu      sync.Mutex
}

var _ Container = (*container)(nil)

func New() Container {
	return &container{
		deps: make(map[string]depInfo),
		vals: make(map[string]any),
	}
}

func (c *container) AddSingleton(key string, fn DepFactoryFunc) {
	c.deps[key] = depInfo{
		key:     key,
		scope:   Singleton,
		factory: fn,
	}
}

func (c *container) AddScoped(key string, fn DepFactoryFunc) {
	c.deps[key] = depInfo{
		key:     key,
		scope:   Scoped,
		factory: fn,
	}
}

func (c *container) Scoped(ctx context.Context) context.Context {
	childContainer := &container{
		parent: c,
		deps:   c.deps,
		vals:   make(map[string]any),
	}

	return context.WithValue(ctx, containerKey, childContainer)
}

func (c *container) Get(key string) any {
	info, exists := c.deps[key]
	if !exists {
		log.Panicf("there is not dependency registered with '%s'", key)
	}

	if _, exists := c.tracked[info.key]; exists {
		log.Panicf("cyclic dependencies encountered while building '%s', tracked: %s", info.key, c.tracked)
	}

	if info.scope == Singleton {
		return c.getFromParent(info)
	}

	return c.get(info)
}

func (c *container) getFromParent(info depInfo) any {
	if c.parent != nil {
		return c.parent.getFromParent(info)
	}

	return c.get(info)
}

func (c *container) get(info depInfo) any {
	c.mu.Lock()

	val, exists := c.vals[info.key]
	if !exists {
		tempValue := make(chan struct{})
		c.vals[info.key] = tempValue

		c.mu.Unlock()

		return c.build(info, tempValue)
	}

	c.mu.Unlock()

	tempValue, isTemp := val.(chan struct{})
	if !isTemp {
		return val
	}

	<-tempValue

	return c.get(info)
}

func (c *container) build(info depInfo, tempValue chan struct{}) any {
	c.mu.Lock()

	val, err := info.factory(&container{
		parent:  c.parent,
		deps:    c.deps,
		vals:    c.vals,
		tracked: c.tracked.add(info),
	})
	if err != nil {
		delete(c.vals, info.key)

		c.mu.Unlock()
		close(tempValue)

		log.Panicf("error building dependency '%s': %s", info.key, err)
	}

	c.vals[info.key] = val

	c.mu.Unlock()
	close(tempValue)

	return val
}
