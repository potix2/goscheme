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

%type<expr> expr proc_call quotation self_evaluating datum compound_datum identifier simple_datum bool number symbol command
%type<exprs> operands data commands program

%token<tok> IDENT UINT10 BOOLEAN

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
            $$ = append([]ast.Expr{$1})
        }
        |
        command commands
        {
            $$ = append([]ast.Expr{$1}, $2...)
        }

command: expr

expr : identifier | quotation | self_evaluating | proc_call

self_evaluating : bool | number

quotation :
        '\'' datum
        {
            $$ = ast.QuoteExpr{$2}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }

datum : simple_datum | compound_datum

simple_datum : bool | number | symbol

compound_datum :
               '(' ')'
               {
                    $$ = ast.MakeListFromSlice([]ast.Expr{})
                    if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
               }
               | '(' data ')'
               {
                    $$ = ast.MakeListFromSlice($2)
                    if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
               }
               | '(' data '.' datum ')'
               {
                    p := ast.MakeListFromSlice($2)
                    p.Cdr = $4
                    $$ = p
                    if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
               }

data :
     datum
     {
        $$ = append([]ast.Expr{$1})
     }
     | datum data
     {
        $$ = append([]ast.Expr{$1}, $2...)
     }

symbol : identifier

identifier :
        IDENT
        {
            $$ = ast.IdentExpr{Lit: $1.Lit}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }

bool :
        BOOLEAN
        {
            lit, _ := strconv.ParseBool($1.Lit)
            $$ = ast.BooleanExpr{Lit: lit}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }

number :
        UINT10
        {
            lit, _ := strconv.Atoi($1.Lit)
            $$ = ast.Uint10Expr{Lit: lit}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }


proc_call :
        '(' expr ')'
        {
            $$ = ast.AppExpr{Exprs: []ast.Expr{$2}}
            if l, ok := yylex.(*Lexer); ok { l.expr = $$ }
        }
        | '(' expr operands ')'
        {
            $$ = ast.AppExpr{Exprs: append([]ast.Expr{$2}, $3...)}
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
