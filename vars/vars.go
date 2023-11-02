package vars

import "sync"

type Vars struct {
	lock  *sync.RWMutex
	_vars map[string]interface{}
}

func NewVars() *Vars {
	return &Vars{
		_vars: make(map[string]interface{}),
	}
}

func (v *Vars) Put(key string, value interface{}) {
	v.lock.RLock()
	defer v.lock.RUnlock()
	v._vars[key] = value
}

func (v *Vars) Get(key string) interface{} {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v._vars[key]
}
