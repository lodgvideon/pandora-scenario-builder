package scenario

import (
	"go.uber.org/zap"
	"pandoraScript/vars"
)

type Context struct {
	vars *vars.Vars
	logs *zap.Logger
	prev *OperationResult
}

func (c *Context) Vars() *vars.Vars {
	return c.vars
}

func (c *Context) V() *vars.Vars {
	return c.vars
}

func (c *Context) SetPrev(prev *OperationResult) {
	c.prev = prev
}

func (c *Context) Prev() *OperationResult {
	return c.prev
}
func (c *Context) L() *zap.Logger {
	return c.logs
}

func (c *Context) Log() *zap.Logger {
	return c.logs
}

func NewContext(l *zap.Logger) *Context {
	return &Context{vars: vars.NewVars(), logs: l}
}
