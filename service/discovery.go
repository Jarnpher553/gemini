package service

import (
	"fmt"
	"github.com/Jarnpher553/gemini/service/selector"
	"strconv"
	"sync"
	"time"

	"github.com/Jarnpher553/gemini/log"
	consul "github.com/hashicorp/consul/api"
)

type Registry struct {
	sync.Mutex
	*consul.Client
	Services []*NodeInfo
	selector selector.Selector
	logger   *log.ZapLogger
}

func NewRegistry(addr string, s ...selector.Selector) *Registry {
	logger := log.Zap.Mark("registry")
	config := consul.DefaultConfig()
	config.Address = addr
	cli, err := consul.NewClient(config)
	if err != nil {
		logger.Fatal(log.Message(err))
	}

	r := &Registry{
		Client:   cli,
		Services: make([]*NodeInfo, 0),
		logger:   logger,
	}
	if len(s) == 0 {
		r.selector = selector.RoundRobin()
	} else {
		r.selector = s[0]
	}
	return r
}

func (r *Registry) InjectSlice(services ...IBaseService) {
	for _, v := range services {
		r.inject(v)
	}
}

func (r *Registry) Inject(service IBaseService) {
	r.inject(service)
}

func (r *Registry) inject(service IBaseService) {
	r.Lock()
	defer r.Unlock()
	r.Services = append(r.Services, service.Node())
	service.SetReg(r)
}

func (r *Registry) Register(node *NodeInfo, group *sync.WaitGroup, errChan chan error) {
	defer group.Done()

	services, _, err := r.Health().Checks(node.ServerName+"."+node.Name, nil)
	if err == nil {
		for _, v := range services {
			if v.ServiceID == node.Id {
				return
			}
		}
	}

	check := &consul.AgentServiceCheck{
		TCP:                            fmt.Sprintf("%s:%s", node.Address, node.Port),
		Interval:                       fmt.Sprintf("%v", time.Second*5),
		DeregisterCriticalServiceAfter: fmt.Sprintf("%v", time.Minute),
	}

	port, _ := strconv.Atoi(node.Port)
	asr := &consul.AgentServiceRegistration{
		ID:      node.Id,
		Name:    fmt.Sprintf("%s.%s", node.ServerName, node.Name),
		Port:    port,
		Address: node.Address,
		Check:   check,
	}

	asr.Connect = &consul.AgentServiceConnect{
		Native: true,
	}

	if err := r.Agent().ServiceRegister(asr); err != nil {
		log.Zap.Mark("registry").Fatal(log.Message(err))

		errChan <- err
		return
	}
	r.logger.Info(log.Messagef(`register service {"id":"%s", "name":"%s"} ok`, node.Id, asr.Name))
}

func (r *Registry) Deregister(node *NodeInfo, group *sync.WaitGroup, errChan chan error) {
	defer group.Done()

	if err := r.Agent().ServiceDeregister(node.Id); err != nil {
		r.logger.Fatal(log.Message(err))

		errChan <- err
		return
	}
	r.logger.Info(log.Messagef(`deregister service {"id":"%s", "name":"%s.%s"} ok`, node.Id, node.ServerName, node.Name))
}

func (r *Registry) GetService(name string) (*NodeInfo, error) {
	rsp, _, err := r.Health().Connect(name, "", false, nil)

	if err != nil {
		return nil, err
	}

	nodes := make([]*NodeInfo, 0)
	for _, s := range rsp {
		if s.Service.Service != name {
			continue
		}

		nodes = append(nodes, &NodeInfo{Id: s.Service.ID, Name: s.Service.Service, Port: strconv.Itoa(s.Service.Port), Address: s.Service.Address})
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("%s service not found", name)
	}
	node := nodes[r.selector(len(nodes))]

	r.logger.Info(log.Messagef(`get service {"id":"%s", "name":"%s"} ok`, node.Id, node.Name))
	return node, nil
}
