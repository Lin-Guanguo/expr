package parser

import (
	"bytes"
	"expr/lexer"
	"fmt"
)

const indentationBlank = 4

func printIndentation(deep int) string {
	if deep != 0 {
		return fmt.Sprintf("%*s", deep*indentationBlank, " ")
	}
	return ""
}

type Expression interface {
	String(deep int) string
}

type AssignableExpr interface {
	Expression
	IsAssignable() bool
}

type FuncCaptureExpr interface {
	Expression
	FuncCaptureExpr()
}

type ValueExpr interface {
	Expression
	ValueExpr()
}

type IntegerExpr struct {
	Token *lexer.Token
	Value int64
}

func (e *IntegerExpr) ValueExpr() { var _ ValueExpr = e }
func (e *IntegerExpr) String(deep int) string {
	return fmt.Sprintf("%s%d", printIndentation(deep), e.Value)
}

type FloatExpr struct {
	Token *lexer.Token
	Value float64
}

func (e *FloatExpr) ValueExpr() { var _ ValueExpr = e }
func (e *FloatExpr) String(deep int) string {
	return fmt.Sprintf("%s%gf", printIndentation(deep), e.Value)
}

type BooleanExpr struct {
	Token *lexer.Token
	Value bool
}

func (e *BooleanExpr) ValueExpr() { var _ ValueExpr = e }
func (e *BooleanExpr) String(deep int) string {
	return fmt.Sprintf("%s%t", printIndentation(deep), e.Value)
}

type StringExpr struct {
	Token *lexer.Token
	Value string
}

func (e *StringExpr) ValueExpr() { var _ ValueExpr = e }
func (e *StringExpr) String(deep int) string {
	return fmt.Sprintf("%s\"%s\"", printIndentation(deep), e.Value)
}

type TableExpr struct {
	Token     *lexer.Token
	InitValue []KeyValuePair
}

type KeyValuePair struct {
	Key   Expression
	Value Expression
}

func (e *TableExpr) ValueExpr() { var _ ValueExpr = e }
func (e *TableExpr) String(deep int) string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("%stable{ key:value", printIndentation(deep)))
	for _, pair := range e.InitValue {
		buf.WriteString(fmt.Sprintf("\n%s%s = %s", printIndentation(deep+1),
			pair.Key.String(0), pair.Value.String(0)))
	}
	buf.WriteString(fmt.Sprintf("\n%s}", printIndentation(deep)))
	return buf.String()
}

type NilExpr struct {
	Token *lexer.Token
}

func (e *NilExpr) ValueExpr() { var _ ValueExpr = e }
func (e *NilExpr) String(deep int) string {
	return fmt.Sprintf("%snil", printIndentation(deep))
}

type PackExpr struct {
	Token *lexer.Token // []
	Exprs []Expression
}

func (e *PackExpr) IsAssignable() bool {
	var _ AssignableExpr = e
	for _, expr := range e.Exprs {
		switch expr := expr.(type) {
		case *PackExpr:
			if expr.IsAssignable() == false {
				return false
			}
		case *Identifier:
		default:
			return false
		}
	}
	return true
}

func (e *PackExpr) String(deep int) string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("%spack[", printIndentation(deep)))
	for i, expr := range e.Exprs {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(expr.String(0))
	}
	buf.WriteString(fmt.Sprintf("]"))
	return buf.String()
}

type Identifier struct {
	Token *lexer.Token
	Ident string
}

func (e *Identifier) IsAssignable() bool { var _ AssignableExpr = e; return true }
func (e *Identifier) FuncCaptureExpr()   { var _ FuncCaptureExpr = e }
func (e *Identifier) String(deep int) string {
	return fmt.Sprintf("%s%s", printIndentation(deep), e.Ident)
}

type ArithPrefixExpr struct {
	Token *lexer.Token
	Op    lexer.TokenType
	Right Expression
}

func (e *ArithPrefixExpr) String(deep int) string {
	var op string
	switch e.Op {
	case lexer.T_MINUS:
		op = "-"
	case lexer.T_BANG:
		op = "!"
	default:
		panic("unknown ArithPrefixExpr operator")
	}
	return fmt.Sprintf("%s%s%s", printIndentation(deep), op, e.Right.String(0))
}

type ArithInfixExpr struct {
	Token *lexer.Token
	Op    lexer.TokenType
	Left  Expression
	Right Expression
}

func (e *ArithInfixExpr) String(deep int) string {
	var op string
	switch e.Op {
	case lexer.T_MINUS:
		op = "-"
	case lexer.T_PLUS:
		op = "+"
	case lexer.T_ASTERISK:
		op = "*"
	case lexer.T_SLASH:
		op = "/"
	case lexer.T_EQ:
		op = "=="
	case lexer.T_NEQ:
		op = "!="
	case lexer.T_LE:
		op = "<="
	case lexer.T_GE:
		op = ">="
	case lexer.T_LT:
		op = "<"
	case lexer.T_GT:
		op = ">"
	case lexer.T_AND:
		op = "and"
	case lexer.T_OR:
		op = "or"
	default:
		panic("unknown ArithInfixExpr operator")
	}
	return fmt.Sprintf("%s(%s %s %s)", printIndentation(deep), e.Left.String(0), op, e.Right.String(0))
}

type DeclarationExpr struct {
	Token *lexer.Token
	Left  AssignableExpr
	Value Expression
}

func (e *DeclarationExpr) FuncCaptureExpr() { var _ FuncCaptureExpr = e }
func (e *DeclarationExpr) String(deep int) string {
	return fmt.Sprintf("%s%s := %s", printIndentation(deep), e.Left.String(0), e.Value.String(0))
}

type AssignExpr struct {
	Token *lexer.Token
	Left  AssignableExpr
	Value Expression
}

func (e *AssignExpr) String(deep int) string {
	return fmt.Sprintf("(%s%s = %s)", printIndentation(deep), e.Left.String(0), e.Value.String(0))
}

type IfExpr struct {
	Token       *lexer.Token
	Condition   Expression
	Consequence *BlockExpr
	Alternative *BlockExpr
}

func (e *IfExpr) String(deep int) string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("%sif %s", printIndentation(deep), e.Condition.String(0)))
	buf.WriteString(fmt.Sprintf("\n%s", e.Consequence.String(deep)))
	if e.Alternative != nil {
		buf.WriteString(fmt.Sprintf("\n%selse", printIndentation(deep)))
		buf.WriteString(fmt.Sprintf("\n%s", e.Alternative.String(deep)))
	}
	return buf.String()
}

type FuncExpr struct {
	Token      *lexer.Token
	Parameters []*Identifier
	Capture    []FuncCaptureExpr
	Body       *BlockExpr
}

func (f *FuncExpr) String(deep int) string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("%sfunc(", printIndentation(deep)))
	for i, ident := range f.Parameters {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(ident.Ident)
	}
	buf.WriteString(")[")
	for i, capture := range f.Capture {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(capture.String(0))
	}
	buf.WriteString("]\n")
	buf.WriteString(f.Body.String(deep + 1))
	return buf.String()
}

type CallExpr struct {
	Token      *lexer.Token
	Function   Expression
	Parameters []Expression
}

func (f *CallExpr) String(deep int) string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("%s%s", printIndentation(deep), f.Function.String(0)))
	buf.WriteString("(")
	for i, param := range f.Parameters {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.String(0))
	}
	buf.WriteString(")")
	return buf.String()
}

type IndexExpr struct {
	Token *lexer.Token
	Table Expression
	Index Expression
}

func (i *IndexExpr) IsAssignable() bool { var _ AssignableExpr = i; return true }
func (i *IndexExpr) String(deep int) string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("%s%s", printIndentation(deep), i.Table.String(0)))
	buf.WriteString(".[")
	buf.WriteString(i.Index.String(0))
	buf.WriteString("]")
	return buf.String()
}

type ReturnExpr struct {
	Token       *lexer.Token
	ReturnValue Expression
}

func (e *ReturnExpr) String(deep int) string {
	return fmt.Sprintf("%sreturn %s", printIndentation(deep), e.ReturnValue.String(0))
}

type BlockExpr struct {
	Token *lexer.Token
	Exprs []Expression
}

func (e *BlockExpr) String(deep int) string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("%s{\n", printIndentation(deep)))
	for _, expr := range e.Exprs {
		buf.WriteString(fmt.Sprintf("%s\n", expr.String(deep+1)))
	}
	buf.WriteString(fmt.Sprintf("%s}", printIndentation(deep)))
	return buf.String()
}

//----------loop----------

type ForExpr struct {
	Token    *lexer.Token
	InitExpr Expression
	EdgeExpr Expression
	StepExpr Expression
	Body     *BlockExpr
}

func (f *ForExpr) String(deep int) string {
	forHead := fmt.Sprintf("%sfor %s; %s; %s\n", printIndentation(deep), f.InitExpr.String(0), f.EdgeExpr.String(0), f.StepExpr.String(0))
	return forHead + f.Body.String(deep)
}

type BreakExpr struct {
	Token      *lexer.Token
	BreakValue Expression
}

func (e *BreakExpr) String(deep int) string {
	return fmt.Sprintf("%sbreak %s", printIndentation(deep), e.BreakValue.String(0))
}

/*type Continue struct {
	Token *lexer.Token
}

func (e *Continue) String(deep int) string {
	return fmt.Sprintf("%scontinue", printIndentation(deep))
}*/
