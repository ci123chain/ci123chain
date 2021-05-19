package commands

import (
	"github.com/ci123chain/ci123chain/pkg/collactor/collactor"
	"sync"

	"golang.org/x/sync/errgroup"
	"time"
)

const StatusConnected = "Connected"
const StatusUnConnected = "UnConnected"

// Service represents a relayer listen service
// TODO: sync services to disk so that they can survive restart
type Service struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Src    string `json:"src"`
	//SrcKey string `json:"src-key"`
	Dst    string `json:"dst"`
	Status string `json:"status"`
	//DstKey string `json:"dst-key"`
	shouldRestart bool
	strategyType  collactor.Strategy
	srcChain      *collactor.Chain
	dstChain      *collactor.Chain
	doneFunc      func()
	terminal      chan struct{}
}

// NewService returns a new instance of Service
func NewService(name, path string, src, dst *collactor.Chain, strategyType collactor.Strategy, doneFunc func()) *Service {
	s := &Service{
		Name: name,
		Path: path,
		Src: src.ChainID,
		Dst: dst.ChainID,
		Status: StatusUnConnected,
		shouldRestart: false,
		strategyType: strategyType,
		doneFunc: doneFunc,
		srcChain: src,
		dstChain: dst,
	}
	s.terminal = make(chan struct{})
	return s
}

func (s *Service) ReStart() {
	collactor.DefaultChainLogger().Info("Service ReStart", "path", s.Path, "src", s.Src, "dst", s.Dst)
	doneFunc, err := collactor.RunStrategy(s.srcChain, s.dstChain, s.strategyType)
	if err == nil {
		s.doneFunc = doneFunc
	}
}

func (s *Service) Delete() {
	s.terminal <- struct{}{}
}
func (s *Service) HealthCheck() {
	tick := time.NewTicker(20 * time.Second)
	for {
		select {
		case <- s.terminal:
			tick.Stop()
			return
		case _ = <- tick.C:
			var eg = new(errgroup.Group)
			eg.Go(func() error {
				err := s.srcChain.HealthCheck()
				collactor.DefaultChainLogger().Error("Chain Endpoint ShutDown", "path", s.Path, "src", s.Src, "desc", err.Error())
				return err
			})
			eg.Go(func() error {
				err := s.srcChain.HealthCheck()
				collactor.DefaultChainLogger().Error("Chain Endpoint ShutDown", "path", s.Path, "dst", s.Dst, "desc", err.Error())
				return err
			})
			if err := eg.Wait(); err != nil {
				// 有节点健康检查失败
				s.doneFunc()
				s.shouldRestart = true
				s.Status = StatusUnConnected
				continue
			}
			s.Status = StatusConnected
			// 健康检查正常 && shouldRestart = true
			if s.shouldRestart {
				s.shouldRestart = false
				s.ReStart()
			}
		}
	}
}


// ServicesManager represents the manager of the various services the relayer is running
type ServicesManager struct {
	Services map[string]*Service

	sync.Mutex
}

// NewServicesManager returns a new instance of a services manager
func NewServicesManager() *ServicesManager {
	return &ServicesManager{Services: make(map[string]*Service)}
}