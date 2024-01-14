package scenario

import (
	"go.uber.org/zap"
	"sync"
)

type IOperation interface {
	SetPrev(o IOperation)
	SetNext(o IOperation)
	Named
	Executor
}

type Named interface {
	SetName(name string)
	Name() string
}

type Executor interface {
	Exec(ctx *Context)
}

type Operation struct {
	name             string
	condition        func(context *Context) bool
	beforeOperations []BeforeOperation
	operation        *OperationFn
	afterOperations  []AfterOperation
	ctx              *Context
	next             IOperation
	prev             IOperation
	result           *OperationResult
}

func (o *Operation) Name() string {
	return o.name
}
func (o *Operation) SetName(name string) {
	o.name = name
}

type OperationBuilder struct {
	_operation *Operation
}

func (o *Operation) Exec(ctx *Context) {
	o.ctx = ctx
	//TODO проверка - будет ли выполнена операция
	if !o.condition(ctx) {
		return
	}
	for _, befOp := range o.beforeOperations {
		//TODO: Сделать обработку ошибок

		_ = befOp(o, ctx)
		//Выполняем список операций до появлении операции с ошибкой

	}
	operation := *o.operation
	if operationResult := operation(o, ctx); operationResult != nil {
		//Вот это место вызывает сомнения. Оно нужно - для доступа к результату предыдущей операии - из контекста
		//следующей. (Например для того - чтобы вытащить Бади ответа от сервера
		ctx.SetPrev(operationResult)
		o.SetResult(operationResult)
	} else {
		ctx.L().Error("Operation not return Result",
			zap.String("name", o.Name()),
		)
	}

	for _, aftOp := range o.afterOperations {
		//TODO: Сделать обработку ошибок
		_ = aftOp(o, ctx)
		//Выполняем список операций до появлении операции с ошибкой
		//	Если ошибка - то фейлить результат операции

	}
}

func (o *Operation) AddBeforeOperation(f BeforeOperation) {
	o.beforeOperations = append(o.beforeOperations, f)
}

func (o *Operation) AddAfterOperation(f AfterOperation) {
	o.afterOperations = append(o.afterOperations, f)
}

func (o *Operation) SetNext(op IOperation) {
	op.SetPrev(o)
	o.next = op
}

func (o *Operation) SetPrev(op IOperation) {
	o.prev = op
}

func (o *Operation) Next() IOperation {
	return o.next
}

func (o *Operation) Prev() IOperation {
	return o.prev
}

func (o *Operation) Result() *OperationResult {
	return o.result
}

func (o *Operation) ReleaseResult() {
	result := o.result
	o.result = nil
	ReleaseResult(result)
	o.ctx.prev = nil
}

func (o *Operation) SetResult(result *OperationResult) {
	o.result = result
}

func (b *OperationBuilder) Build() *Operation {
	operation := *b._operation
	return &operation
}

func (b *OperationBuilder) BeforeOperation(f BeforeOperation) *OperationBuilder {
	b._operation.AddBeforeOperation(f)
	return b
}

func (b *OperationBuilder) SetOperation(o OperationFn) *OperationBuilder {
	b._operation.operation = &o
	return b
}

func (b *OperationBuilder) After(o AfterOperation) *OperationBuilder {
	b._operation.AddAfterOperation(o)
	return b
}

func (b *OperationBuilder) SetName(name string) *OperationBuilder {
	b._operation.SetName(name)
	return b
}

type OperationFn = func(op *Operation, ctx *Context) *OperationResult

// NewAbstractOperationBuilder Init buuilder for abstract Operation.
// In abstract operation, any code can be executed and metered
func NewAbstractOperationBuilder(name string) *OperationBuilder {
	o := &Operation{name: name,
		condition: func(context *Context) bool {
			return true
		},
	}
	return &OperationBuilder{_operation: o}
}

type BeforeOperation func(operation *Operation, ctx *Context) error
type AfterOperation func(operation *Operation, ctx *Context) error

var resultPool = sync.Pool{New: func() interface{} {
	return &OperationResult{
		Request:  nil,
		Response: nil,
	}
}}

func AcquireResult() *OperationResult {
	return resultPool.Get().(*OperationResult)

}

func ReleaseResult(res *OperationResult) {
	res.Clear()
	resultPool.Put(res)
}
func (r *OperationResult) Clear() {
	r.Request = nil
	r.Response = nil
}

type OperationResult struct {
	Request  interface{}
	Response interface{}
	//sample netsample.Sample
}

func (o *OperationResult) Req() interface{} {
	return o.Request
}

func (o *OperationResult) Resp() interface{} {
	return o.Response
}
