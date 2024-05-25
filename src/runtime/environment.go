package runtime

import "fmt"

type environment struct {
	Parent     *environment
	Variables  map[string]Variable
	StructDefs map[string]StructDef
	Functions  map[string]Function
}

func NewEnv(parent *environment) *environment {
	env := environment{
		Parent:     parent,
		Variables:  make(map[string]Variable),
		StructDefs: make(map[string]StructDef),
		Functions:  make(map[string]Function),
	}

	declareNativeFunctions(env)
	declareNativeValues(env)

	return &env
}

func (e *environment) containsVar(varName string) bool {
	_, exists := e.Variables[varName]
	return exists
}

func (e *environment) declareVar(varName string, value RuntimeVal, expectedType ValueType, isConst bool) RuntimeVal {
	if e.containsVar(varName) {
		panic("Already has var")
	}

	if expectedType == "" {
		expectedType = AnyType
	}

	if !checkType(value.Type(), expectedType) {
		panic(fmt.Sprintf("expected %s got %s", expectedType, value.Type()))
	}

	v := Variable{Value: value, ExpectedType: expectedType, Constant: isConst}
	e.Variables[varName] = v

	return v

}

func (e *environment) resolveVar(varName string) *environment {
	if e.containsVar(varName) {
		return e
	}
	if e.Parent != nil {
		return e.Parent.resolveVar(varName)
	}

	panic(fmt.Sprintf("Variable %s not found.", varName))
}

func (e *environment) lookupVar(varName string) Variable {
	env := e.resolveVar(varName)
	return env.Variables[varName]
}

func (e *environment) declareStruct(structName string, properties map[string]ValueType) RuntimeVal {
	s := StructDef{Name: structName, Properties: properties, Methods: make(map[string]Function)}
	e.StructDefs[structName] = s
	return s
}

func (e *environment) containsStruct(structName string) bool {
	_, exists := e.StructDefs[structName]
	return exists
}

func (e *environment) implMethod(structName string, method Function) RuntimeVal {

	if !e.containsStruct(structName) {
		panic("Struct not defined")
	}

	e.StructDefs[structName].Methods[method.Name] = method
	return method
}

func (e *environment) containsFn(fnName string) bool {
	_, exists := e.Functions[fnName]
	return exists
}

func (e *environment) declareFn(fn Function) RuntimeVal {
	if e.containsFn(fn.Name) {
		panic("")
	}

	e.Functions[fn.Name] = fn

	return fn
}

func (e *environment) resolveStruct(structName string) *environment {
	if e.containsStruct(structName) {
		return e
	}
	if e.Parent != nil {
		return e.Parent.resolveStruct(structName)
	}

	panic(fmt.Sprintf("Variable %s not found.", structName))
}

func (e *environment) lookupStruct(structName string) RuntimeVal {
	env := e.resolveStruct(structName)
	return env.StructDefs[structName]
}

func (e *environment) resolveFn(fnName string) *environment {
	if e.containsFn(fnName) {
		return e
	}
	if e.Parent != nil {
		return e.Parent.resolveFn(fnName)
	}

	panic(fmt.Sprintf("Function %s not found.", fnName))
}

func (e *environment) lookupFn(fnName string) RuntimeVal {
	env := e.resolveFn(fnName)
	return env.Functions[fnName]
}

func (e *environment) assignVar(varName string, value RuntimeVal) Variable {

	env := e.resolveVar(varName)
	variable := env.Variables[varName]

	if variable.Constant {
		panic("")
	}

	if !checkType(value.Type(), variable.ExpectedType) {
		panic("")
	}

	env.Variables[varName] = Variable{Value: value, ExpectedType: variable.ExpectedType, Constant: variable.Constant}
	return variable
}

func (e *environment) assignStruct(varName string, memberName string, value RuntimeVal) RuntimeVal {
	env := e.resolveVar(varName)
	variable := env.Variables[varName]

	structDef := e.lookupStruct(string(variable.Value.Type())).(StructDef)
	structVal := variable.Value.(Struct)

	_, memberExists := structDef.Properties[memberName]
	if !memberExists {
		panic("")
	}

	structVal.Properties[memberName] = value

	env.Variables[varName] = Variable{
		Value:        structVal,
		ExpectedType: variable.ExpectedType,
		Constant:     false,
	}

	return structVal
}

func (e *environment) declareNativeFn(fnName string, call FunctionCall) {
	e.Functions[fnName] = Function{Name: fnName, NativeFn: NativeFunction{call}}
}
