package parser

import (
	"bytes"
	"expr/lexer"
	"testing"
)

func TestParseValue(t *testing.T) {
	//block :=
	simpleTestParse(t, "10 10.12 .12 12. true false nil \"hello\"")
	//fmt.Println(block.String(0))
}

func TestArithExpr(t *testing.T) {
	//block :=
	simpleTestParse(t, `
-1
!false
1 + 2 + 3 + 4 + 5 + 6
1+1 2*1  1 - 2  4 / 2;;;
1+2/3*4+2
4 * (1 + -2 * (2 + 1))
true and 21 or 12 and -ident or 1 == 2 and 3 < 4
`)
	//fmt.Println(block.String(0))
}

func TestBlockExpr(t *testing.T) {
	//block :=
	simpleTestParse(t, `
10
{
	20
}
{{ false }}
`)
	//fmt.Println(block.String(0))
}

func TestIfExpr(t *testing.T) {
	//block :=
	simpleTestParse(t, `
if 2 > 1 { 10 } else 20
if false { 20 }
if 10 false else true
`)
	//fmt.Println(block.String(0))
}

func TestReturnBreakExpr(t *testing.T) {
	//block :=
	simpleTestParse(t, `
return 10
return 1+ 2*3
return nil
break nil
`)
	//fmt.Println(block.String(0))
}

func TestDeclarationAssignExpr(t *testing.T) {
	//block :=
	simpleTestParse(t, `
ident := 20
x = y = z = 1+1+1
[x,y,z] = [1,2,3]
[x,y,z,[x,y]] := [1,2,3,[1,2]]
`)
	//fmt.Println(block.String(0))
}

func TestFuncExpr(t *testing.T) {
	//block :=
	simpleTestParse(t, `
func(){ return 10 }
func()[] 10086 
outer := 20
outer2 := 30
f := func(x, y)[outer, copy := outer2, x] { return x + y + outer + copy}
`)
	//fmt.Println(block.String(0))
}

func TestTableExpr(t *testing.T) {
	//block :=
	simpleTestParse(t, `
table{
	[1] = 2,
	hello = 2,
	[false] = true,
	[2.3] = 1.23
}
`)
	//fmt.Println(block.String(0))
}

func TestCallExpr(t *testing.T) {
	//block :=
	simpleTestParse(t, `
f(1,2,3)
x := f(1+2,y,3)
func(x,y)[]{ return x + y }(1,2)
func(x,y)[] return x; (1,2)
`)
	//fmt.Println(block.String(0))
}

func TestIndexExpr(t *testing.T) {
	//block :=
	simpleTestParse(t, `
t = table{}

t.a = 10
t.hello()
t.[10]
t.[t]
`)
	//fmt.Println(block.String(0))
}

func TestForExpr(t *testing.T) {
	//block :=
	simpleTestParse(t, `
for x := 1; x < 100; x = x + 1 y:=1

for i := 1; i < 10; i = 1 {
	t.[i] = 10
	break 10
}
`)
	//fmt.Println(block.String(0))
}

func TestPackExpr(t *testing.T) {
	//block :=
	simpleTestParse(t, `
x := [1,2,3]
[1+2, 3, 4,]
`)
	//fmt.Println(block.String(0))
}

func simpleTestParse(t *testing.T, input string) *BlockExpr {
	buf := bytes.NewBufferString(input)
	l := lexer.New(buf)
	p := New(l)

	block := p.ParseProgram()
	if len(p.Errors) != 0 {
		t.Errorf("%s\n", block.String(0))
		for i, e := range p.Errors {
			t.Errorf("[%d] %s\n", i, e.Error())
		}
	}
	return block
}
