%{
package parser
import (
    "strconv"
    "github.com/potix2/goscheme/scm"
)
%}

%union{
    exprs []scm.Expr
    expr  scm.Expr
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
            $$ = append([]scm.Expr{$1})
        }
        |
        command commands
        {
            $$ = append([]scm.Expr{$1}, $2...)
        }

command: expr

expr : quotation | self_evaluating | identifier | proc_call

self_evaluating : bool | number | string

quotation :
        '\'' datum
        {
            $$ = scm.QuoteExpr{$2}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }

datum : simple_datum | compound_datum

simple_datum : bool | number | symbol

compound_datum :
               '(' ')'
               {
                    $$ = scm.MakeListFromSlice([]scm.Expr{})
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
        $$ = append([]scm.Expr{$1})
     }
     | datum data
     {
        $$ = append([]scm.Expr{$1}, $2...)
     }

symbol : identifier

identifier :
        IDENT
        {
            $$ = scm.IdentExpr{Lit: $1.Lit}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }

bool :
        BOOLEAN
        {
            lit, _ := strconv.ParseBool($1.Lit)
            $$ = scm.BooleanExpr{Lit: lit}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }

number : NUMBER
       {
            $$ = scm.StringToNumber($1.Lit)
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
       }

string : STRING
       {
            $$ = scm.StringExpr($1.Lit)
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
       }

proc_call :
        '(' expr ')'
        {
            $$ = scm.AppExpr{Exprs: []scm.Expr{$2}}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }
        | '(' expr operands ')'
        {
            $$ = scm.AppExpr{Exprs: append([]scm.Expr{$2}, $3...)}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }
operands :
         expr
         {
            $$ = append([]scm.Expr{$1})
         }
         |
         expr operands
         {
            $$ = append([]scm.Expr{$1}, $2...)
         }
%%
