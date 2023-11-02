package scenario

import "go.uber.org/zap"

type Operation struct {
	name             string
	condition        func(context *Context) bool
	beforeOperations []BeforeOperation
	operation        *OperationFn
	afterOperations  []AfterOperation
	ctx              Context
	next             *Operation
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

func (o *Operation) Exec(op *Operation, ctx *Context) {
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
		ctx.SetPrev(operationResult)
	} else {
		ctx.L().Error("Operation not return Result",
			zap.String("name", o.Name()),
		)
	}

	for _, aftOp := range o.afterOperations {
		//TODO: Сделать обработку ошибок
		_ = aftOp(o, ctx)
		//Выполняем список операций до появлении операции с ошибкой

	}
}

func (o *Operation) AddBeforeOperation(f BeforeOperation) {
	o.beforeOperations = append(o.beforeOperations, f)
}

func (o *Operation) AddAfterOperation(f AfterOperation) {
	o.afterOperations = append(o.afterOperations, f)
}

func (o *Operation) SetNext(op *Operation) {
	o.next = op
}

func (o *Operation) Next() *Operation {
	return o.next
}

func (b *OperationBuilder) Build() *Operation {
	return b._operation
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

type OperationFn = func(op *Operation, ctx *Context) *OperationResult

func NewOperationBuilder(name string) *OperationBuilder {
	o := &Operation{name: name,
		condition: func(context *Context) bool {
			return true
		},
	}
	return &OperationBuilder{_operation: o}
}

type BeforeOperation func(operation *Operation, ctx *Context) error
type AfterOperation func(operation *Operation, ctx *Context) error

type OperationResult struct {
	Request  interface{}
	Response interface{}
	//sample netsample.Sample
}

func (o *OperationResult) Req() interface{} {
	return o.Request
}

func (o *OperationResult) Res() interface{} {
	return o.Response
}
