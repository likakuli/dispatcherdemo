package main

// 设计一个分发器，当有数据进入分发器后，需要将数据分发到多个处理器处理，每个处理器可以想象为一个协程，处理器在没有数据的时候要阻塞
// ----------        /--> processor
// chan --> dispatch ---> processor
// ----------        \--> processor
type dispatcher struct {
	processors []*processor
}

func newDispatcher() *dispatcher {
	return &dispatcher{}
}

func (d *dispatcher) addProcessor(handler func(obj interface{})) {
	p := newProcessor(handler)
	d.processors = append(d.processors, p)
}

func (d *dispatcher) dispatch(obj interface{}) {
	for i := 0; i < len(d.processors); i++ {
		d.processors[i].add(obj)
	}
}

func (d *dispatcher) run(stopCh chan struct{}) {
	for i := 0; i < len(d.processors); i++ {
		go d.processors[i].run()
		go d.processors[i].pop()
	}
	<-stopCh
	for i := 0; i < len(d.processors); i++ {
		close(d.processors[i].addCh)
	}
}

type processor struct {
	nextCh chan interface{}
	addCh  chan interface{}

	buffer *RingGrowing

	handler func(obj interface{})
}

func newProcessor(handler func(obj interface{})) *processor {
	return &processor{
		handler: handler,
		nextCh:  make(chan interface{}),
		addCh:   make(chan interface{}),
		buffer:  NewRingGrowing(1024 * 1024),
	}
}

func (p *processor) run() {
	for {
		select {
		case obj, ok := <-p.nextCh:
			if !ok {
				return
			}
			p.handler(obj)
		}
	}
}

func (p *processor) add(obj interface{}) {
	p.addCh <- obj
}

func (p *processor) pop() {
	defer close(p.nextCh)

	var (
		nextCh     chan<- interface{}
		pendingObj interface{}
		ok         bool
	)

	for {
		select {
		case nextCh <- pendingObj:
			pendingObj, ok = p.buffer.ReadOne()
			if !ok {
				nextCh = nil
			}
		case obj, ok := <-p.addCh:
			if !ok {
				return
			}
			if pendingObj == nil {
				nextCh = p.nextCh
				pendingObj = obj
			} else {
				p.buffer.WriteOne(obj)
			}
		}
	}
}
