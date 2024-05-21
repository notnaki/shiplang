package interpreter

import (
	"fmt"
	"strings"
	"time"
)

type environment struct {
	Parent     *environment
	Variables  map[string]Variable
	Constants  map[string]struct{}
	StructDefs map[string]Struct
}

func showFN(args []RuntimeVal) RuntimeVal {
	for _, arg := range args {
		inspected := arg.Inspect()
		cleaned := strings.Trim(inspected, "{}")
		fmt.Print(cleaned)
	}
	fmt.Println()
	return MKNULL()
}

func rangeFN(args []RuntimeVal) RuntimeVal {
	if len(args) < 1 || len(args) > 2 {
		panic("range function expects one or two arguments")
	}

	// Ensure the first argument is a number
	if args[0].Type() != NumberType {
		panic("range function expects a number as the first argument")
	}

	start := int(args[0].(NumberVal).Value)
	var end int

	// If two arguments are provided, ensure the second argument is also a number
	if len(args) == 2 {
		if args[1].Type() != NumberType {
			panic("range function expects a number as the second argument")
		}
		end = int(args[1].(NumberVal).Value)
	} else {
		// If only one argument is provided, default end value to start
		end = start
		start = 0
	}

	// Create a range of numbers between start and end (inclusive)
	result := make([]ArrayElement, end-start+1)
	for i := 0; i <= end-start; i++ {
		result[i] = ArrayElement{Value: MKNUM(float64(start + i)), ElementType: NumberType}
	}

	return Array{Elements: result}
}

func timeFN(_ []RuntimeVal) RuntimeVal {

	return MKNUM(float64(time.Now().UnixMilli()))
}

func dateFN(_ []RuntimeVal) RuntimeVal {
	currentTime := time.Now()
	formattedDateTime := currentTime.Format("15:04:05.000 02-01-2006")

	return MKSTR(formattedDateTime)
}

func sleepFN(_ []RuntimeVal) RuntimeVal {
	time.Sleep(2 * time.Second)
	return MKNULL()
}

func CreateGlobalEnv() *environment {
	env := environment{
		Parent:     nil,
		Variables:  make(map[string]Variable),
		Constants:  make(map[string]struct{}),
		StructDefs: make(map[string]Struct),
	}

	env.declare_var("true", MKBOOL(true), BooleanType, true)
	env.declare_var("false", MKBOOL(false), BooleanType, true)
	env.declare_var("null", MKNULL(), NullType, true)

	env.declare_var("show", MKNATIVEFN(showFN), NativeFnType, true)
	env.declare_var("time", MKNATIVEFN(timeFN), NativeFnType, true)
	env.declare_var("date", MKNATIVEFN(dateFN), NativeFnType, true)
	env.declare_var("sleep", MKNATIVEFN(sleepFN), NativeFnType, true)
	env.declare_var("range", MKNATIVEFN(rangeFN), NativeFnType, true)

	return &env
}

func (e *environment) declare_var(varName string, value RuntimeVal, varType ValueType, isConst bool) RuntimeVal {
	if e.contains_var(varName) {
		panic(fmt.Sprintf("Cannot declare variable %s: it is already defined.", varName))
	}

	if varType == "" {
		varType = AnyType
	}

	valueType := value.Type()
	if valueType == StructType && varType != AnyType {
		if value.(StructInstance).StructName != string(varType) {
			panic(fmt.Sprintf("Variable %s is of type %s, got %s", varName, varType, value.(StructInstance).StructName))
		}
	} else if value.Type() == ArrayType && varType != AnyType {
		if ValueType(value.Inspect()) != varType {
			panic(fmt.Sprintf("Variable %s is of type %s, got %s", varName, varType, value.Inspect()))
		}
	} else if valueType != varType && varType != AnyType && valueType != NullType && valueType != StructType && value.Type() != ArrayType {
		panic(fmt.Sprintf("Variable %s is of type %s, got %s", varName, varType, valueType))
	}

	v := Variable{Value: value, VarType: varType}
	e.Variables[varName] = v

	if isConst {
		e.Constants[varName] = struct{}{}
	}

	return v
}

func (e *environment) resolve(varName string) *environment {
	if e.contains_var(varName) {
		return e
	}
	if e.Parent != nil {
		return e.Parent.resolve(varName)
	}

	panic(fmt.Sprintf("Variable %s not found.", varName))
}

func (e *environment) lookup_var(varName string) RuntimeVal {
	env := e.resolve(varName)
	return env.Variables[varName]
}

func (e *environment) assign_var(varName string, value RuntimeVal) RuntimeVal {
	env := e.resolve(varName)

	if _, exists := env.Constants[varName]; exists {
		panic(fmt.Sprintf("Cannot reassign to variable %s as it was declared constant.", varName))
	}

	variable := env.Variables[varName]
	if variable.VarType != value.Type() && variable.VarType != AnyType {
		panic(fmt.Sprintf("Variable %s is of type %s, got %s", varName, variable.VarType, value.Type()))
	}

	env.Variables[varName] = Variable{VarType: variable.VarType, Value: value}
	return variable
}

func (e *environment) contains_var(varName string) bool {
	_, exists := e.Variables[varName]
	return exists
}

func (e *environment) resolve_struct(structName string) (Struct, bool) {
	s, ok := e.StructDefs[structName]
	if ok {
		return s, true
	}
	if e.Parent != nil {
		return e.Parent.resolve_struct(structName)
	}
	return Struct{}, false
}

func (e *environment) declare_struct(structName string, properties []Property) RuntimeVal {
	s := Struct{Name: structName, Properties: properties}
	e.StructDefs[structName] = s
	return s
}

func (e *environment) lookup_struct(structName string) (RuntimeVal, bool) {
	f, ok := e.resolve_struct(structName)
	return f, ok
}

func (e *environment) assign_struct_member(varName string, memberName string, value RuntimeVal) RuntimeVal {
	env := e.resolve(varName)

	// Ensure the variable exists and is a struct
	varVariable, exists := env.Variables[varName]
	if !exists {
		panic(fmt.Sprintf("Variable %s not found.", varName))
	}
	if varVariable.Value.Type() != StructType {
		panic(fmt.Sprintf("Variable %s is not a struct.", varName))
	}

	// Get the struct instance from the variable value
	structInstance := varVariable.Value.(StructInstance)

	// Get the struct definition from the environment
	structDef, ok := env.lookup_struct(structInstance.StructName)
	if !ok {
		panic(fmt.Sprintf("Struct %s not defined.", structInstance.StructName))
	}

	// Check if the member exists in the struct definition
	memberExists := false
	for _, prop := range structDef.(Struct).Properties {
		if prop.PropName == memberName {
			memberExists = true
			break
		}
	}
	if !memberExists {
		panic(fmt.Sprintf("Member %s not found in struct %s.", memberName, structInstance.StructName))
	}

	structInstance.Properties[memberName] = value

	env.Variables[varName] = Variable{
		VarType: varVariable.VarType,
		Value:   structInstance,
	}

	return value
}

func (e *environment) delete_variable(varName string) {
	if !e.contains_var(varName) {
		panic(fmt.Sprintf("Variable %s does not exist.", varName))
	}

	delete(e.Variables, varName)
	delete(e.Constants, varName)
}

func isDefaultVariable(varName string) bool {
	defaultVariables := map[string]struct{}{
		"true":  {},
		"false": {},
		"null":  {},
		"show":  {},
		"time":  {},
		"date":  {},
		"sleep": {},
		"range": {},
	}
	_, isDefault := defaultVariables[varName]
	return isDefault
}
