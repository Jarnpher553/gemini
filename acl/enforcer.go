package acl

import (
	"errors"
	"sync"
)

type Enforcer struct {
	mutex    sync.RWMutex
	adapter  IAdapter
	policies []Policy
}

type Policy struct {
	Subject interface{}
	Domain  interface{}
	Object  interface{}
	Action  interface{}
}

type Request struct {
	Subject interface{}
	Domain  interface{}
	Object  interface{}
	Action  interface{}
}

type IAdapter interface {
	LoadPolicy() ([]Policy, error)
	Enforce() func(*Request) error
}

func NewEnforcer(adapter IAdapter) *Enforcer {
	return &Enforcer{adapter: adapter}
}

func (e *Enforcer) LoadPolicy() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	policies, err := e.adapter.LoadPolicy()
	if err != nil {
		return err
	}
	e.policies = policies
	return nil
}

func (e *Enforcer) Enforce(request *Request) error {
	f := e.adapter.Enforce()
	if f != nil {
		return f(request)
	} else {
		e.mutex.RLock()
		defer e.mutex.RUnlock()
		for _, policy := range e.policies {
			if policy.Subject == request.Subject &&
				policy.Domain == request.Domain &&
				policy.Object == request.Object &&
				policy.Action == request.Action {
				return nil
			}
		}
		return errors.New("can't access")
	}
}
