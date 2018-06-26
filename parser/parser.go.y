%{
package parser
import (
    "strconv"
    "github.com/potix2/goscheme/ast"
)
%}

%union{
    exprs []ast.Expr
    expr  ast.Expr
    tok   Token
}

%type<expr> expr proc_call
%type<exprs> operands

%token<tok> IDENT UINT10

%%

expr :
        IDENT
        {
            $$ = ast.IdentExpr{Lit: $1.Lit}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }
        | UINT10
        {
            lit, _ := strconv.Atoi($1.Lit)
            $$ = ast.Uint10Expr{Lit: lit}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }
        | proc_call
        {
            $$ = $1
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }

proc_call :
        '(' expr ')'
        {
            $$ = ast.AppExpr{Operator: $2, Operands: []ast.Expr{}}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }
        | '(' expr operands ')'
        {
            $$ = ast.AppExpr{Operator: $2, Operands: $3}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }
operands :
         expr
         {
            $$ = append([]ast.Expr{$1})
         }
         |
         expr operands
         {
            $$ = append([]ast.Expr{$1}, $2...)
         }
%%
