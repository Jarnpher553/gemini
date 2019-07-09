package service

import (
	"fmt"
	"github.com/Janrpher553/micro-core/log"
	consul "github.com/hashicorp/consul/api"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Registry struct {
	sync.Mutex
	*consul.Client
	Services  []*NodeInfo
}

func NewRegistry(addr string) *Registry {
	config := consul.DefaultConfig()
	config.Address = addr
	cli, err := consul.NewClient(config)
	if err != nil {
		log.Logger.Mark("Registry").Fatalln(err)
	}

	return &Registry{
		Client:   cli,
		Services: make([]*NodeInfo, 0),
	}
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
		log.Logger.Mark("Registry").Fatalln(err)

		errChan <- err
		return
	}
	log.Logger.Mark("Registry").Infof(`register service {"id":"%s", "name":"%s"} ok`, node.Id, asr.Name)
}

func (r *Registry) Deregister(node *NodeInfo, group *sync.WaitGroup, errChan chan error) {
	defer group.Done()

	if err := r.Agent().ServiceDeregister(node.Id); err != nil {
		log.Logger.Mark("Registry").Fatalln(err)

		errChan <- err
		return
	}
	log.Logger.Mark("Registry").Infof(`deregister service {"id":"%s", "name":"%s.%s"} ok`, node.Id, node.ServerName, node.Name)
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
	node := selector(nodes)

	log.Logger.Mark("Registry").Infof(`get service {"id":"%s", "name":"%s"} ok`, node.Id, node.Name)
	return node, nil
}

func selector(nodes []*NodeInfo) *NodeInfo {
	l := len(nodes)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	index := r.Intn(l)

	return nodes[index]
}
