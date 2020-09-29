package acl

import (
	"errors"
	cmap "github.com/orcaman/concurrent-map"
)

type Enforcer struct {
	adapter  IAdapter
	policies cmap.ConcurrentMap
}

type PolicyKeyValuePair struct {
	Subject string
	Polices []Policy
}

type Policy struct {
	Subject string
	Domain  string
	Object  string
	Action  string
}

type Request struct {
	Subject string
	Domain  string
	Object  string
	Action  string
}

type IAdapter interface {
	LoadPolicy(a ...interface{}) (*PolicyKeyValuePair, error)
}

func NewEnforcer(adapter IAdapter) *Enforcer {
	return &Enforcer{adapter: adapter, policies: cmap.New()}
}

func (e *Enforcer) LoadPolicy(a ...interface{}) error {
	policyKeyValuePair, err := e.adapter.LoadPolicy()
	if err != nil {
		return err
	}
	e.policies.Set(policyKeyValuePair.Subject, policyKeyValuePair.Polices)
	return nil
}

func (e *Enforcer) Enforce(request *Request) error {
	val, ok := e.policies.Get(request.Subject)
	if !ok {
		return errors.New("can't access")
	}
	policies := val.([]Policy)
	for _, policy := range policies {
		if policy.Subject == request.Subject &&
			policy.Domain == request.Domain &&
			policy.Object == request.Object &&
			policy.Action == request.Action {
			return nil
		}
	}
	return errors.New("can't access")
}
