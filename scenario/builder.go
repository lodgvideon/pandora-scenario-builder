package scenario

import (
	"go.uber.org/zap"
)

func NewScenarioBuilder(name string) *Builder {
	return &Builder{name: name, _scenario: &Scenario{name: name}}
}

type Builder struct {
	_scenario *Scenario
	log       *zap.Logger
	name      string
}

func (b *Builder) Build() *Scenario {
	//Тут Инжект всех зависимостей
	return b._scenario
}

func (b *Builder) AddOperation(o *Operation) {

	if len(b._scenario.operations) > 0 {
		if op := b._scenario.operations[len(b._scenario.operations)]; op != nil {
			op.SetNext(o)
		}
	}

	b._scenario.operations = append(b._scenario.operations, o)
}

func (b *Builder) SetLogger(log *zap.Logger) {
	b._scenario.SetLogger(log)
}

type Scenario struct {
	name       string
	onStart    func() error
	operations []*Operation
	onEnd      func() error
	log        *zap.Logger
}

func (s *Scenario) SetLogger(logger *zap.Logger) {
	s.log = logger
}

func (s *Scenario) Run(ctx *Context) {
	//TODO: Сделать обработку ошибок
	if s.onStart != nil {
		_ = s.onStart()
	}

	for _, op := range s.operations {
		op.Exec(op, ctx)
	}
	//TODO: Сделать обработку ошибок
	if s.onEnd != nil {
		_ = s.onEnd()
	}

}
