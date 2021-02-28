package evaluator

type Environment struct {
	LocalVars map[string]*Object
	Outer     *Environment
}

func NewEnv() *Environment {
	return &Environment{
		LocalVars: make(map[string]*Object),
		Outer:     nil,
	}
}

func NewInnerEnv(outer *Environment) *Environment {
	return &Environment{
		LocalVars: make(map[string]*Object),
		Outer:     outer,
	}
}

func (e *Environment) Get(varName string) *Object {
	searchEnv := e
	for searchEnv != nil {
		obj := searchEnv.LocalVars[varName]
		if obj != nil {
			return obj
		}
		searchEnv = searchEnv.Outer
	}
	return nil
}

// 可共享底层obj
func (e *Environment) Set(varName string, object *Object) {
	e.LocalVars[varName] = object
}

func (e *Environment) SetNewObj(varName string, object Object) {
	e.Set(varName, &object)
}

func (e *Environment) Close() {
	e.Outer = nil
}
