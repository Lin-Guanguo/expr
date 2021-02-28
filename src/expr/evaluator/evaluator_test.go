package evaluator

import (
	"bytes"
	"expr/lexer"
	"expr/parser"
	"testing"
)

func TestEnv(t *testing.T) {
	env1 := NewEnv()
	env2 := NewInnerEnv(env1)
	env3 := NewInnerEnv(env2)

	env1.SetNewObj("hello", StringObj{Value: "world"})
	env2.SetNewObj("hello", StringObj{Value: "shijie"})
	env2.SetNewObj("k2", StringObj{Value: "v2"})

	if o := *(env1.Get("hello")); o.(StringObj).Value != "world" {
		t.Error("error")
	}
	if o := *(env2.Get("hello")); o.(StringObj).Value != "shijie" {
		t.Error("error")
	}
	if o := *(env3.Get("hello")); o.(StringObj).Value != "shijie" {
		t.Error("error")
	}

	if env1.Get("k2") != nil {
		t.Error("error")
	}
	if o := *(env2.Get("k2")); o.(StringObj).Value != "v2" {
		t.Error("error")
	}
	if o := *(env3.Get("k2")); o.(StringObj).Value != "v2" {
		t.Error("error")
	}
}

func TestBasicValueEval(t *testing.T) {
	testProgram(t, `return 10`, IntegerObj{Value: 10})
	testProgram(t, `return true`, BooleanObj{Value: true})
	testProgram(t, `return false`, BooleanObj{Value: false})
	testProgram(t, `return 10.2`, FloatObj{Value: 10.2})
	testProgram(t, `return "hello"`, StringObj{Value: "hello"})

}

func TestTableEval(t *testing.T) {
	obj, _ := testProgram(t, `return table{ hello = "world", [false] = 10, [10.2] = true}`, nil)
	table := obj.(TableObj).Table.Store
	if table[StringObj{"hello"}].(StringObj).Value != "world" {
		t.Errorf("TestTableEval error")
	}
	if table[BooleanObj{false}].(IntegerObj).Value != 10 {
		t.Errorf("TestTableEval error")
	}
	if table[FloatObj{10.2}].(BooleanObj).Value != true {
		t.Errorf("TestTableEval error")
	}
}

func TestDeclareAssign(t *testing.T) {
	testProgram(t, `x := 10; return x`, IntegerObj{Value: 10})
	testProgram(t, `x := 10.2; return x`, FloatObj{Value: 10.2})
	testProgram(t, `x := 10; x = false; return x`, BooleanObj{Value: false})

	testProgram(t, `[x,y,z] := [10,20,30]; return x`, IntegerObj{Value: 10})
	testProgram(t, `[x,y,z] := [10,20,30]; return y`, IntegerObj{Value: 20})
	testProgram(t, `[x,y,z] := [10,20,30]; return z`, IntegerObj{Value: 30})

	testProgram(t, `pack := [7,20,30]; [x] := pack; return x`, IntegerObj{Value: 7})
}

func TestArith(t *testing.T) {
	testProgram(t, `return -1`, IntegerObj{Value: -1})
	testProgram(t, `return --10`, IntegerObj{Value: 10})
	testProgram(t, `return !true`, BooleanObj{Value: false})
	testProgram(t, `return !!10086`, BooleanObj{Value: true})

	testProgram(t, `return 1 + 2`, IntegerObj{Value: 3})
	testProgram(t, `return 1 + 2 * 3`, IntegerObj{Value: 7})
	testProgram(t, `return 1 + 2 * ( 3 + 1)`, IntegerObj{Value: 9})

	testProgram(t, `return 1.0 + 2.5`, FloatObj{Value: 3.5})
	testProgram(t, `return 2 * 2 + 2.5`, FloatObj{Value: 6.5})
	testProgram(t, `return 2 * ( 2 + 2.5)`, FloatObj{Value: 9.0})

	testProgram(t, `x := 20; return x + 2.5`, FloatObj{Value: 22.5})

	testProgram(t, `x := true; return x and 10.2 or 20`, FloatObj{Value: 10.2})
	testProgram(t, `x := false; return x and 10.2 or 20`, IntegerObj{Value: 20})

	testProgram(t, `return 10 < 20`, BooleanObj{Value: true})
	testProgram(t, `return 10 > 20`, BooleanObj{Value: false})
	testProgram(t, `return 10.2 < 20.5`, BooleanObj{Value: true})
	testProgram(t, `return 10 < 20.5`, BooleanObj{Value: true})
}

func TestIndex(t *testing.T) {
	testProgram(t, `t := table{ str = 10 }; return t.str`, IntegerObj{Value: 10})
	testProgram(t, `t := table{ str = 10 }; return t.world`, NilObj)
	testProgram(t, `t := table{ [10086] = 10.2 }; return t.[10086]`, FloatObj{Value: 10.2})

	testProgram(t, `t := table{ message = 10.2 }; y := t.message; return y`, FloatObj{Value: 10.2})
	testProgram(t, `t := table{ message = [10,20,30] }; [x] := t.message; return x`, IntegerObj{Value: 10})
	testProgram(t, `t := table{ message = [10,20,30] }; [x,y] := t.message; return y`, IntegerObj{Value: 20})

	testProgram(t, `t := table{ }; t.str = "hello";  return t.str`, StringObj{Value: "hello"})
	testProgram(t, `t := table{ }; t.[false] = "hello";  return t.[false]`, StringObj{Value: "hello"})
}

func TestIfEval(t *testing.T) {
	testProgram(t, `if x := true return 10 else return 20`, IntegerObj{Value: 10})
	testProgram(t, `if x := false { return 10 } else { return 20 }`, IntegerObj{Value: 20})
	testProgram(t, `if x := true { return x } return 20 `, BooleanObj{Value: true})

	testProgram(t, `return if 10 { break 10 } else { break 20 } `, IntegerObj{Value: 10})
	testProgram(t, `return if 10 break 10 else break 20 `, IntegerObj{Value: 10})
}

func TestFuncCall(t *testing.T) {
	testProgram(t, `return func(){return 10}()`, IntegerObj{Value: 10})
	testProgram(t, `f := func(){return 10}; return f()`, IntegerObj{Value: 10})
	testProgram(t, `f := func(x,y,z,w){ return if x > y break z else break w}; return f(1, 2, 0.5, 1.0)`,
		FloatObj{Value: 1.0})
	testProgram(t, `f := func(x,y,z,w){ return if x > y break z else break w}; return f(1, 2, 0.5)`,
		NilObj)

	testProgram(t, `x := 10; y := 20; f := func()[x]{return x}; return f()`, IntegerObj{Value: 10})
	testProgram(t, `
f1 := func(){ 
	x := 0; 
	return [
		func()[x]{ x = x + 1 }, 
		func()[x]{ return x },
		func() [ x := x ] return x + 10
	]
}

[f2, f3, f4] := f1()
f2() f2() f2()
return f3() + f4()
`, IntegerObj{Value: 13})
}

func testProgram(t *testing.T, input string, expect Object) (Object, *Environment) {
	inputReader := bytes.NewBufferString(input)
	l := lexer.New(inputReader)
	p := parser.New(l)
	block := p.ParseProgram()
	if p.Errors != nil {
		t.Errorf("parse Error input: %s \n%s\n", input, block.String(0))
		for i, e := range p.Errors {
			t.Errorf("[%d] %s\n", i, e.Error())
		}
	}

	env := NewEnv()
	obj, err := evalFuncBlockExpr(block, env)
	if err != nil {
		t.Errorf("Eval Error input: %s\n%s", input, err.Error())
	}
	if expect != nil && obj != expect {
		t.Errorf("expect Object mismatch Error input: %s", input)
	}
	return obj, env
}
