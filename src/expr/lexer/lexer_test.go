package lexer

import (
	"bytes"
	"testing"
)

func TestLexer(t *testing.T) {
	input := `
int := 10
float := 10.1
float1 := 10.
float2 := .10
    
float = 20.1

if int == 10 {
	return 20
} else {
	if int == 20 {
		return 30
	}
}

t := table{}

t.function = func[](x,y,z){
	return x+y+z
}

t.function(1,2,3)
t:method(1,2,3)

+ - * / == != > < >= <= ()[]{} . , : :=

hello
hello_23_world
"hello"
`
	expect := `T_IDENT line: 2; int;
T_DECLARATION line: 2; ;
T_INT line: 2; 10;
T_IDENT line: 3; float;
T_DECLARATION line: 3; ;
T_FLOAT line: 3; 10.1;
T_IDENT line: 4; float1;
T_DECLARATION line: 4; ;
T_FLOAT line: 4; 10.0;
T_IDENT line: 5; float2;
T_DECLARATION line: 5; ;
T_FLOAT line: 5; 0.10;
T_IDENT line: 7; float;
T_ASSIGN line: 7; ;
T_FLOAT line: 7; 20.1;
T_IF line: 9; ;
T_IDENT line: 9; int;
T_EQ line: 9; ;
T_INT line: 9; 10;
T_LBRACE line: 9; ;
T_RETURN line: 10; ;
T_INT line: 10; 20;
T_RBRACE line: 11; ;
T_ELSE line: 11; ;
T_LBRACE line: 11; ;
T_IF line: 12; ;
T_IDENT line: 12; int;
T_EQ line: 12; ;
T_INT line: 12; 20;
T_LBRACE line: 12; ;
T_RETURN line: 13; ;
T_INT line: 13; 30;
T_RBRACE line: 14; ;
T_RBRACE line: 15; ;
T_IDENT line: 17; t;
T_DECLARATION line: 17; ;
T_TABLE line: 17; ;
T_LBRACE line: 17; ;
T_RBRACE line: 17; ;
T_IDENT line: 19; t;
T_DOT line: 19; ;
T_IDENT line: 19; function;
T_ASSIGN line: 19; ;
T_FUNCTION line: 19; ;
T_LBRACKET line: 19; ;
T_RBRACKET line: 19; ;
T_LPAREN line: 19; ;
T_IDENT line: 19; x;
T_COMMA line: 19; ;
T_IDENT line: 19; y;
T_COMMA line: 19; ;
T_IDENT line: 19; z;
T_RPAREN line: 19; ;
T_LBRACE line: 19; ;
T_RETURN line: 20; ;
T_IDENT line: 20; x;
T_PLUS line: 20; ;
T_IDENT line: 20; y;
T_PLUS line: 20; ;
T_IDENT line: 20; z;
T_RBRACE line: 21; ;
T_IDENT line: 23; t;
T_DOT line: 23; ;
T_IDENT line: 23; function;
T_LPAREN line: 23; ;
T_INT line: 23; 1;
T_COMMA line: 23; ;
T_INT line: 23; 2;
T_COMMA line: 23; ;
T_INT line: 23; 3;
T_RPAREN line: 23; ;
T_IDENT line: 24; t;
T_COLON line: 24; ;
T_IDENT line: 24; method;
T_LPAREN line: 24; ;
T_INT line: 24; 1;
T_COMMA line: 24; ;
T_INT line: 24; 2;
T_COMMA line: 24; ;
T_INT line: 24; 3;
T_RPAREN line: 24; ;
T_PLUS line: 26; ;
T_MINUS line: 26; ;
T_ASTERISK line: 26; ;
T_SLASH line: 26; ;
T_EQ line: 26; ;
T_NEQ line: 26; ;
T_GT line: 26; ;
T_LT line: 26; ;
T_GE line: 26; ;
T_LE line: 26; ;
T_LPAREN line: 26; ;
T_RPAREN line: 26; ;
T_LBRACKET line: 26; ;
T_RBRACKET line: 26; ;
T_LBRACE line: 26; ;
T_RBRACE line: 26; ;
T_DOT line: 26; ;
T_COMMA line: 26; ;
T_COLON line: 26; ;
T_DECLARATION line: 26; ;
T_IDENT line: 28; hello;
T_IDENT line: 29; hello_23_world;
T_STRING line: 30; hello;
T_EOF line: 31; EOF;
`
	buf := bytes.NewBufferString(input)
	output := bytes.Buffer{}
	l := New(buf)
	for {
		t := l.NextToken()
		output.WriteString(t.String() + "\n")
		if t.Type == T_EOF {
			break
		}
	}
	if output.String() != expect {
		t.Errorf("Lexer ouput error")
		t.Error(output.String())
	}
}
