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

func (b *Builder) AddOnStartOperation(beforeScenario BeforeScenario) *Builder {
	b._scenario.onStartOperations = append(b._scenario.onStartOperations, beforeScenario)
	return b
}

func (b *Builder) AddOnEndOperation(afterScenario AfterScenario) *Builder {
	b._scenario.onEndOperations = append(b._scenario.onEndOperations, afterScenario)
	return b
}

func (b *Builder) AddOperation(o IOperation) *Builder {

	if len(b._scenario.operations) > 0 {
		if op := b._scenario.operations[len(b._scenario.operations)-1]; op != nil {
			op.SetNext(o)
		}
	}

	b._scenario.operations = append(b._scenario.operations, o)
	return b
}

func (b *Builder) SetLogger(log *zap.Logger) *Builder {
	b._scenario.SetLogger(log)
	return b
}

type BeforeScenario func(s *Scenario, ctx *Context) error
type AfterScenario func(s *Scenario, ctx *Context) error

type Scenario struct {
	name              string
	onStartOperations []BeforeScenario
	operations        []IOperation
	onEndOperations   []AfterScenario
	log               *zap.Logger
}

func (s *Scenario) SetLogger(logger *zap.Logger) {
	s.log = logger
}

func (s *Scenario) Run(ctx *Context) {
	//TODO: Сделать обработку ошибок

	for _, berOf := range s.onStartOperations {
		if err := berOf(s, ctx); err != nil {
			s.log.Error("Error on before operation", zap.Error(err))
		}
	}

	for _, op := range s.operations {
		op.Exec(ctx)
	}

	for _, afterScenarioOperation := range s.onEndOperations {
		if err := afterScenarioOperation(s, ctx); err != nil {
			s.log.Error("Error on After scenario operation", zap.Error(err))
		}
	}
}
