%{
package parser
import (
    "strconv"
    "github.com/potix2/goscheme/scm"
)
%}

%union{
    exprs []scm.Object
    expr  scm.Object
    tok   Token
}

%type<expr> expr proc_call quotation self_evaluating datum compound_datum identifier simple_datum bool symbol command number string
%type<exprs> operands data commands program
%token<tok> IDENT NUMBER BOOLEAN STRING

%%

program: 
       commands
       {
            $$ = $1
            if l, ok := yylex.(*Lexer); ok { l.exprs = $$ }
       }

commands:
        command
        {
            $$ = append([]scm.Object{$1})
        }
        |
        command commands
        {
            $$ = append([]scm.Object{$1}, $2...)
        }

command: expr

expr : quotation | self_evaluating | identifier | proc_call

self_evaluating : bool | number | string

quotation :
        '\'' datum
        {
            $$ = scm.Quote{$2}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }

datum : simple_datum | compound_datum

simple_datum : bool | number | symbol

compound_datum :
               '(' ')'
               {
                    $$ = scm.MakeListFromSlice([]scm.Object{})
                    if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
               }
               | '(' data ')'
               {
                    $$ = scm.MakeListFromSlice($2)
                    if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
               }
               | '(' data '.' datum ')'
               {
                    p := scm.MakeListFromSlice($2)
                    p.Cdr = $4
                    $$ = p
                    if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
               }

data :
     datum
     {
        $$ = append([]scm.Object{$1})
     }
     | datum data
     {
        $$ = append([]scm.Object{$1}, $2...)
     }

symbol : identifier

identifier :
        IDENT
        {
            $$ = scm.Symbol{Lit: $1.Lit}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }

bool :
        BOOLEAN
        {
            lit, _ := strconv.ParseBool($1.Lit)
            $$ = scm.Boolean{Lit: lit}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }

number : NUMBER
       {
            $$ = scm.StringToNumber($1.Lit)
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
       }

string : STRING
       {
            $$ = scm.String($1.Lit)
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
       }

proc_call :
        '(' expr ')'
        {
            $$ = scm.Subp{Objs: []scm.Object{$2}}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }
        | '(' expr operands ')'
        {
            $$ = scm.Subp{Objs: append([]scm.Object{$2}, $3...)}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }
operands :
         expr
         {
            $$ = append([]scm.Object{$1})
         }
         |
         expr operands
         {
            $$ = append([]scm.Object{$1}, $2...)
         }
%%
