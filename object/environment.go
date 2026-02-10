package object

type Environment struct {
	store      map[string]Object
	constStore map[string]bool
	outer      *Environment
}

func InitEnv() *Environment {
	return &Environment{
		store:      make(map[string]Object),
		constStore: make(map[string]bool),
		outer:      nil,
	}
}

func InitEnclosedEnv(outer *Environment) *Environment {
	env := InitEnv()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	if e.constStore[name] {
		return &Error{
			Message: "error const cannot be mutated",
		}
	}

	if e.isOuterHas(name) {
		if e.IsConst(name) {
			return nil
		}
	}

	e.store[name] = val
	return val
}

func (e *Environment) SetConst(name string, val Object) Object {
	e.store[name] = val
	e.constStore[name] = true
	return val
}

func (e *Environment) isOuterHas(name string) bool {
	if e.outer == nil {
		return false
	}

	if _, ok := e.outer.Get(name); !ok {
		return false
	}

	if !e.outer.IsConst(name) {
		return false
	}

	return true
}

func (e *Environment) IsConst(name string) bool {
	if e.constStore[name] {
		return true
	}

	if e.outer != nil {
		return e.outer.IsConst(name)
	}

	return false
}
