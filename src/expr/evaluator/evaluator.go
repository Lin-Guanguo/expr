package evaluator

import (
	"expr/lexer"
	"expr/parser"
	"fmt"
)

func Eval(e parser.Expression, env *Environment) (Object, *EvalError) {
	switch e := e.(type) {
	//基本
	case *parser.IntegerExpr:
		return IntegerObj{Value: e.Value}, nil
	case *parser.FloatExpr:
		return FloatObj{Value: e.Value}, nil
	case *parser.BooleanExpr:
		return BooleanObj{Value: e.Value}, nil
	case *parser.StringExpr:
		return StringObj{Value: e.Value}, nil
	case *parser.NilExpr:
		return NilObj, nil
	case *parser.TableExpr:
		return evalTableExpr(e, env)
	case *parser.PackExpr:
		return evalPackExpr(e, env)

	//算数
	case *parser.ArithPrefixExpr:
		return evalArithPrefixExpr(e, env)
	case *parser.ArithInfixExpr:
		return evalArithInfixExpr(e, env)

	//控制
	case *parser.BlockExpr:
		inner := NewInnerEnv(env)
		return evalBlockExpr(e, inner)
	case *parser.IndexExpr:
		return evalIndexExpr(e, env)
	case *parser.IfExpr:
		return evalIfExpr(e, env)
	case *parser.FuncExpr:
		return evalFuncExpr(e, env)
	case *parser.CallExpr:
		return evalCallExpr(e, env)

	//变量
	case *parser.DeclarationExpr:
		return evalDeclarationExpr(e, env)
	case *parser.AssignExpr:
		return evalAssignExpr(e, env)
	case *parser.Identifier:
		return evalIdentifierExpr(e, env)

	//包装
	case *parser.ReturnExpr:
		obj, err := Eval(e.ReturnValue, env)
		if err != nil {
			return nil, err
		}
		return ReturnObj{Value: obj}, nil
	case *parser.BreakExpr:
		obj, err := Eval(e.BreakValue, env)
		if err != nil {
			return nil, err
		}
		return BreakObj{Value: obj}, nil
	}
	panic(fmt.Sprintf("Eval unhand expression: %s", e.String(0)))
}

func evalTableExpr(expr *parser.TableExpr, env *Environment) (Object, *EvalError) {
	table := &TableValue{Store: make(map[Object]Object)}
	for _, pair := range expr.InitValue {
		key, err := Eval(pair.Key, env)
		if err != nil {
			return nil, err
		}
		value, err := Eval(pair.Value, env)
		if err != nil {
			return nil, err
		}
		table.Store[key] = value
	}
	return TableObj{Table: table}, nil
}

func evalPackExpr(expr *parser.PackExpr, env *Environment) (Object, *EvalError) {
	pack := &PackValue{Objs: make([]Object, 0)}
	for _, e := range expr.Exprs {
		obj, err := Eval(e, env)
		if err != nil {
			return nil, err
		}
		pack.Objs = append(pack.Objs, obj)
	}
	return PackObj{Pack: pack}, nil
}

func toBooleanObj(obj Object) BooleanObj {
	switch obj := obj.(type) {
	case BooleanObj:
		return obj
	case NilValue:
		return BooleanObj{Value: false}
	default:
		return BooleanObj{Value: true}
	}
}

func toFloatObj(obj Object) FloatObj {
	switch obj := obj.(type) {
	case FloatObj:
		return obj
	case IntegerObj:
		return FloatObj{Value: float64(obj.Value)}
	default:
		panic("toFloatObj with wrong parameter")
	}
}

func evalArithPrefixExpr(expr *parser.ArithPrefixExpr, env *Environment) (Object, *EvalError) {
	right, err := Eval(expr.Right, env)
	if err != nil {
		return nil, err
	}
	if expr.Op == lexer.T_BANG {
		right := toBooleanObj(right)
		return BooleanObj{Value: !right.Value}, nil
	} else if expr.Op == lexer.T_MINUS {
		switch right := right.(type) {
		case IntegerObj:
			return IntegerObj{Value: -right.Value}, nil
		case FloatObj:
			return FloatObj{Value: -right.Value}, nil
		default:
			return nil, &EvalError{
				Message: fmt.Sprintf("minus prefix operator with wrong value type %s", right.Type().String()),
			}
		}
	} else {
		return nil, &EvalError{Message: "unknown Arith Prefix operator"}
	}
}

func evalArithInfixExpr(expr *parser.ArithInfixExpr, env *Environment) (Object, *EvalError) {
	left, err := Eval(expr.Left, env)
	if err != nil {
		return nil, err
	}
	right, err := Eval(expr.Right, env)
	if err != nil {
		return nil, err
	}
	switch expr.Op {
	case lexer.T_AND:
		if toBooleanObj(left).Value {
			return right, nil
		} else {
			return left, nil
		}
	case lexer.T_OR:
		if toBooleanObj(left).Value {
			return left, nil
		} else {
			return right, nil
		}
	case lexer.T_PLUS, lexer.T_MINUS, lexer.T_ASTERISK, lexer.T_SLASH,
		lexer.T_LT, lexer.T_LE, lexer.T_GT, lexer.T_GE, lexer.T_EQ, lexer.T_NEQ:
		if left.Type() == right.Type() {
			switch expr.Op {
			case lexer.T_PLUS:
				if left.Type() == TIntegerObj {
					return IntegerObj{Value: left.(IntegerObj).Value + right.(IntegerObj).Value}, nil
				} else if left.Type() == TFloatObj {
					return FloatObj{Value: left.(FloatObj).Value + right.(FloatObj).Value}, nil
				}
			case lexer.T_MINUS:
				if left.Type() == TIntegerObj {
					return IntegerObj{Value: left.(IntegerObj).Value - right.(IntegerObj).Value}, nil
				} else if left.Type() == TFloatObj {
					return FloatObj{Value: left.(FloatObj).Value - right.(FloatObj).Value}, nil
				}
			case lexer.T_ASTERISK:
				if left.Type() == TIntegerObj {
					return IntegerObj{Value: left.(IntegerObj).Value * right.(IntegerObj).Value}, nil
				} else if left.Type() == TFloatObj {
					return FloatObj{Value: left.(FloatObj).Value * right.(FloatObj).Value}, nil
				}
			case lexer.T_SLASH:
				if left.Type() == TIntegerObj {
					return IntegerObj{Value: left.(IntegerObj).Value / right.(IntegerObj).Value}, nil
				} else if left.Type() == TFloatObj {
					return FloatObj{Value: left.(FloatObj).Value / right.(FloatObj).Value}, nil
				}
			case lexer.T_LT:
				if left.Type() == TIntegerObj {
					return BooleanObj{Value: left.(IntegerObj).Value < right.(IntegerObj).Value}, nil
				} else if left.Type() == TFloatObj {
					return BooleanObj{Value: left.(FloatObj).Value < right.(FloatObj).Value}, nil
				}
			case lexer.T_LE:
				if left.Type() == TIntegerObj {
					return BooleanObj{Value: left.(IntegerObj).Value <= right.(IntegerObj).Value}, nil
				} else if left.Type() == TFloatObj {
					return BooleanObj{Value: left.(FloatObj).Value <= right.(FloatObj).Value}, nil
				}
			case lexer.T_GT:
				if left.Type() == TIntegerObj {
					return BooleanObj{Value: left.(IntegerObj).Value > right.(IntegerObj).Value}, nil
				} else if left.Type() == TFloatObj {
					return BooleanObj{Value: left.(FloatObj).Value > right.(FloatObj).Value}, nil
				}
			case lexer.T_GE:
				if left.Type() == TIntegerObj {
					return BooleanObj{Value: left.(IntegerObj).Value >= right.(IntegerObj).Value}, nil
				} else if left.Type() == TFloatObj {
					return BooleanObj{Value: left.(FloatObj).Value >= right.(FloatObj).Value}, nil
				}
			case lexer.T_EQ:
				switch left.Type() {
				case TIntegerObj:
					return BooleanObj{Value: left.(IntegerObj).Value == right.(IntegerObj).Value}, nil
				case TFloatObj:
					return BooleanObj{Value: left.(FloatObj).Value == right.(FloatObj).Value}, nil
				case TBooleanObj:
					return BooleanObj{Value: left.(BooleanObj).Value == right.(BooleanObj).Value}, nil
				case TStringObj:
					return BooleanObj{Value: left.(StringObj).Value == right.(StringObj).Value}, nil
				case TTableObj:
					return BooleanObj{Value: left.(TableObj).Table == right.(TableObj).Table}, nil
				case TPackObj:
					return BooleanObj{Value: left.(PackObj).Pack == right.(PackObj).Pack}, nil
				case TNilObj:
					return BooleanObj{Value: true}, nil
				}
			case lexer.T_NEQ:
				switch left.Type() {
				case TIntegerObj:
					return BooleanObj{Value: left.(IntegerObj).Value != right.(IntegerObj).Value}, nil
				case TFloatObj:
					return BooleanObj{Value: left.(FloatObj).Value != right.(FloatObj).Value}, nil
				case TBooleanObj:
					return BooleanObj{Value: left.(BooleanObj).Value != right.(BooleanObj).Value}, nil
				case TStringObj:
					return BooleanObj{Value: left.(StringObj).Value != right.(StringObj).Value}, nil
				case TTableObj:
					return BooleanObj{Value: left.(TableObj).Table != right.(TableObj).Table}, nil
				case TPackObj:
					return BooleanObj{Value: left.(PackObj).Pack != right.(PackObj).Pack}, nil
				case TNilObj:
					return BooleanObj{Value: false}, nil
				}
			}
		} else if (left.Type() == TIntegerObj || left.Type() == TFloatObj) &&
			(right.Type() == TIntegerObj || right.Type() == TFloatObj) {
			switch expr.Op {
			case lexer.T_PLUS:
				return FloatObj{Value: toFloatObj(left).Value + toFloatObj(right).Value}, nil
			case lexer.T_MINUS:
				return FloatObj{Value: toFloatObj(left).Value - toFloatObj(right).Value}, nil
			case lexer.T_ASTERISK:
				return FloatObj{Value: toFloatObj(left).Value * toFloatObj(right).Value}, nil
			case lexer.T_SLASH:
				return FloatObj{Value: toFloatObj(left).Value / toFloatObj(right).Value}, nil
			case lexer.T_LT:
				return BooleanObj{Value: toFloatObj(left).Value < toFloatObj(right).Value}, nil
			case lexer.T_GT:
				return BooleanObj{Value: toFloatObj(left).Value > toFloatObj(right).Value}, nil
			case lexer.T_LE:
				return BooleanObj{Value: toFloatObj(left).Value <= toFloatObj(right).Value}, nil
			case lexer.T_GE:
				return BooleanObj{Value: toFloatObj(left).Value >= toFloatObj(right).Value}, nil
			case lexer.T_EQ:
				return BooleanObj{Value: toFloatObj(left).Value == toFloatObj(right).Value}, nil
			case lexer.T_NEQ:
				return BooleanObj{Value: toFloatObj(left).Value != toFloatObj(right).Value}, nil
			}
		}
	}
	return nil, &EvalError{Message: fmt.Sprintf("Arith Infix Expr operator and operand: %s", expr.String(0))}
}

func evalFuncBlockExpr(block *parser.BlockExpr, env *Environment) (Object, *EvalError) {
	obj, err := evalBlockExpr(block, env)
	if err != nil {
		return nil, err
	}
	switch obj := obj.(type) {
	case ReturnObj:
		return obj.Value, nil
	default:
		return obj, nil
	}
}

func evalBlockExpr(block *parser.BlockExpr, env *Environment) (Object, *EvalError) {
	for _, expr := range block.Exprs {
		obj, err := Eval(expr, env)
		if err != nil {
			return nil, err
		}
		switch obj := obj.(type) {
		case ReturnObj:
			return obj, nil
		case BreakObj:
			return obj.Value, nil
		}
	}
	return NilObj, nil
}

func evalIndexExpr(expr *parser.IndexExpr, env *Environment) (Object, *EvalError) {
	table, err := Eval(expr.Table, env)
	if err != nil {
		return nil, err
	}
	if table, ok := table.(TableObj); ok {
		index, err := Eval(expr.Index, env)
		if err != nil {
			return nil, err
		}
		if res, ok := table.Table.Store[index]; ok {
			return res, nil
		} else {
			return NilObj, nil
		}
	} else {
		return nil, &EvalError{Message: "eval an index of non table"}
	}
}

func evalIfExpr(expr *parser.IfExpr, env *Environment) (Object, *EvalError) {
	ifEnv := NewInnerEnv(env)
	condition, err := Eval(expr.Condition, ifEnv)
	if err != nil {
		return nil, err
	}
	if toBooleanObj(condition).Value {
		return Eval(expr.Consequence, ifEnv)
	} else if expr.Alternative != nil {
		return Eval(expr.Alternative, ifEnv)
	} else {
		return NilObj, nil
	}
}

func evalFuncExpr(expr *parser.FuncExpr, env *Environment) (Object, *EvalError) {
	funcCaptureEnv := NewInnerEnv(env)
	for _, capture := range expr.Capture {
		switch capture := capture.(type) {
		case *parser.DeclarationExpr:
			_, err := evalDeclarationExpr(capture, funcCaptureEnv)
			if err != nil {
				return nil, err
			}
		case *parser.Identifier:
			outerObj := env.Get(capture.Ident)
			if outerObj == nil {
				return nil, &EvalError{Message: "FuncExpr capture a wrong identifier"}
			}
			funcCaptureEnv.LocalVars[capture.Ident] = outerObj
		default:
			panic("evalFuncExpr unhand capture expr type")
		}
	}
	funcCaptureEnv.Close()
	return FuncObj{Func: &FuncValue{
		FuncEnv:    funcCaptureEnv,
		Parameters: expr.Parameters,
		Body:       expr.Body,
	}}, nil
}

func evalCallExpr(expr *parser.CallExpr, env *Environment) (Object, *EvalError) {
	fn, err := Eval(expr.Function, env)
	if err != nil {
		return nil, err
	}
	if fnObj, ok := fn.(FuncObj); !ok {
		return nil, &EvalError{Message: fmt.Sprintf("call to a non function object %s", fn.Type())}
	} else {
		funcCallEnv := NewInnerEnv(fnObj.Func.FuncEnv)
		needNParam := len(fnObj.Func.Parameters)
		giveNParam := len(expr.Parameters)
		for i := 0; ; i++ {
			var err *EvalError
			if i < needNParam && i < giveNParam {
				var obj Object
				obj, err = Eval(expr.Parameters[i], env)
				funcCallEnv.LocalVars[fnObj.Func.Parameters[i].Ident] = &obj
			} else if i < needNParam && i >= giveNParam {
				funcCallEnv.LocalVars[fnObj.Func.Parameters[i].Ident] = &NilObj
			} else if i >= needNParam && i < giveNParam {
				_, err = Eval(expr.Parameters[i], env)
			} else {
				break
			}
			if err != nil {
				return nil, err
			}
		}
		return evalFuncBlockExpr(fnObj.Func.Body, funcCallEnv)
	}
}

func evalDeclarationExpr(expr *parser.DeclarationExpr, env *Environment) (Object, *EvalError) {
	value, err := Eval(expr.Value, env)
	if err != nil {
		return nil, err
	}
	err = declareAssignHelper(expr.Left, value, env,
		func(ident string, value *Object) *EvalError {
			env.Set(ident, value)
			return nil
		})
	if err != nil {
		return nil, err
	}
	return value, nil
}

func evalAssignExpr(expr *parser.AssignExpr, env *Environment) (Object, *EvalError) {
	value, err := Eval(expr.Value, env)
	if err != nil {
		return nil, err
	}
	err = declareAssignHelper(expr.Left, value, env,
		func(ident string, value *Object) *EvalError {
			if o := env.Get(ident); o == nil {
				return &EvalError{Message: fmt.Sprintf("Assign value haven't declare %s", ident)}
			} else {
				*o = *value
				return nil
			}
		})
	if err != nil {
		return nil, err
	}
	return value, nil
}

func declareAssignHelper(
	left parser.AssignableExpr,
	value Object,
	env *Environment,
	fn func(ident string, value *Object) *EvalError) *EvalError {

	switch left := left.(type) {
	case *parser.Identifier:
		err := fn(left.Ident, &value)
		if err != nil {
			return err
		}
	case *parser.IndexExpr:
		table, err := Eval(left.Table, env)
		if err != nil {
			return err
		}
		if table, ok := table.(TableObj); ok {
			index, err := Eval(left.Index, env)
			if err != nil {
				return err
			}
			table.Table.Store[index] = value
		} else {
			return &EvalError{Message: "assign to an index of non table"}
		}
	case *parser.PackExpr:
		switch value := value.(type) {
		case PackObj:
			for i, left := range left.Exprs {
				if len(value.Pack.Objs) > i {
					err := declareAssignHelper(left.(parser.AssignableExpr), value.Pack.Objs[i], env, fn)
					if err != nil {
						return err
					}
				} else {
					err := declareAssignHelper(left.(parser.AssignableExpr), NilObj, env, fn)
					if err != nil {
						return err
					}
				}
			}
		default:
			if len(left.Exprs) > 0 {
				err := declareAssignHelper(left.Exprs[0].(parser.AssignableExpr), value, env, fn)
				if err != nil {
					return err
				}
			}
		}
	default:
		panic("declareAssignHelper unhand left type")
	}
	return nil
}

func evalIdentifierExpr(expr *parser.Identifier, env *Environment) (Object, *EvalError) {
	o := env.Get(expr.Ident)
	if o == nil {
		return nil, &EvalError{Message: fmt.Sprintf("evalIdentifierExpr: identifier haven't declare %s", expr.Ident)}
	} else {
		return *o, nil
	}
}
