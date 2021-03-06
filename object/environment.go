package object

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, thoughts: []string{}, quotes: make([]Quote, 0)}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

type Environment struct {
	store    map[string]Object
	thoughts []string
	quotes   []Quote
	outer    *Environment
}

func (e *Environment) Thoughts() []string {
	return e.thoughts
}

func (e *Environment) AddThought(thought string) {
	e.thoughts = append(e.thoughts, thought)
}

func (e *Environment) AddQuote(q Quote) {
	e.quotes = append(e.quotes, q)
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Quotes() []Quote {
	return e.quotes
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) Store() map[string]Object {
	return e.store
}
