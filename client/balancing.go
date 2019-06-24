package client

import (
	"github.com/youpenglai/mfwgo/registry"
	"fmt"
	"google.golang.org/grpc/resolver"
)

const (
	mfwScheme = "mfw"
)

type mfwResolverBuilder struct {}

func (*mfwResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	r := &mfwResolver{
		target:target,
		cc: cc,
		rn: make(chan struct{}, 1),
	}

	if err := r.update(); err != nil {
		return nil, err
	}
	go r.watcher()

	return r, nil
}

func (*mfwResolverBuilder) Scheme() string {return mfwScheme}

type mfwResolver struct {
	target resolver.Target
	cc resolver.ClientConn
	rn chan struct{}
}

func (r *mfwResolver) update() error {
	services, err := registry.NewConsulService().GetServices(r.target.Endpoint)
	if err != nil {
		return err
	}

	addrs := make([]resolver.Address, len(services))
	for i, svc := range services {
		addrs[i] = resolver.Address{Addr: fmt.Sprintf("%s:%d", svc.Address, svc.Port)}
	}

	r.cc.UpdateState(resolver.State{Addresses:addrs})

	return nil
}

func (r *mfwResolver) watcher() {
	for {
		select {
		case <-r.rn:
			r.update()
		}
	}
}

func (r *mfwResolver) ResolveNow(o resolver.ResolveNowOption) {
	select {
	case r.rn <- struct{}{}:
	default:
	}
}
func (*mfwResolver) Close() {}


func init() {
	resolver.Register(&mfwResolverBuilder{})
}