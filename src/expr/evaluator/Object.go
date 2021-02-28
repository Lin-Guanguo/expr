package evaluator

import "expr/parser"

type ObjType int

type Object interface {
	Type() ObjType
}

const (
	TIdentObj ObjType = iota
	TIntegerObj
	TFloatObj
	TBooleanObj
	TStringObj
	TTableObj
	TPackObj
	TFuncObj
	TNilObj
	TReturnObj
	TBreakObj
)

func (t ObjType) String() string {
	switch t {
	case TIdentObj:
		return "TIdentObj"
	case TIntegerObj:
		return "TIntegerObj"
	case TFloatObj:
		return "TFloatObj"
	case TBooleanObj:
		return "TBooleanObj"
	case TStringObj:
		return "TStringObj"
	case TTableObj:
		return "TTableObj"
	case TPackObj:
		return "TPackObj"
	case TNilObj:
		return "TNilObj"
	case TReturnObj:
		return "TReturnObj"
	case TBreakObj:
		return "TBreakObj"
	default:
		panic("func (t ObjType) String() string")
	}
}

type IdentObj struct {
	Ident string
}

func (o IdentObj) Type() ObjType {
	return TIdentObj
}

type IntegerObj struct {
	Value int64
}

func (o IntegerObj) Type() ObjType {
	return TIntegerObj
}

type FloatObj struct {
	Value float64
}

func (o FloatObj) Type() ObjType {
	return TFloatObj
}

type BooleanObj struct {
	Value bool
}

func (o BooleanObj) Type() ObjType {
	return TBooleanObj
}

type StringObj struct {
	Value string
}

func (o StringObj) Type() ObjType {
	return TStringObj
}

//引用类型

type TableValue struct {
	Store map[Object]Object
}
type TableObj struct {
	Table *TableValue
}

func (t TableObj) Type() ObjType {
	return TTableObj
}

type PackValue struct {
	Objs []Object
}
type PackObj struct {
	Pack *PackValue
}

func (t PackObj) Type() ObjType {
	return TPackObj
}

type FuncValue struct {
	FuncEnv    *Environment
	Parameters []*parser.Identifier
	Body       *parser.BlockExpr
}
type FuncObj struct {
	Func *FuncValue
}

func (f FuncObj) Type() ObjType {
	return TFuncObj
}

type NilValue struct{}

var NilObj Object = NilValue{}

func (o NilValue) Type() ObjType {
	return TNilObj
}

//包装类型

type ReturnObj struct {
	Value Object
}

func (o ReturnObj) Type() ObjType {
	return TReturnObj
}

type BreakObj struct {
	Value Object
}

func (o BreakObj) Type() ObjType {
	return TBreakObj
}
