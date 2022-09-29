grammar Sql;


//java -jar antlr-4.10.1-complete.jar  -Dlanguage=Go -o parser Sql.g4 -no-listener -visitor

sqlQuery
  : selectQuery EOF;

selectQuery
  : K_SELECT resultColumns (whereStatement)? EOQ                                     #selectSimple
  | K_SELECT aggregationColumns windowTumbling (whereStatement)? EOQ                 #selectTumbling
  | K_SELECT groupByAggregationColumns windowTumbling (whereStatement)? groupBy EOQ  #selectTumblingGroupBy
  ;

windowTumbling
   : K_WINDOW_TUMBLING NUMERIC_LITERAL
   ;


resultColumns
 : column (COMMA column)* # selectColumns
 | STAR                   # selectStar
 ;


aggregationColumns
 :   aggregationColumn (COMMA aggregationColumn)*                       #selectAggregations
 ;

groupByAggregationColumns
 :  (column | aggregationColumn) (COMMA (aggregationColumn | column))*  #selectGroupByAggregations
 ;


aggregationColumn
  :  (K_MIN | K_MAX | K_COUNT | K_AVG | K_SUM) L_BRACKET ( IDENTIFIER | IDENTIFIER DOT IDENTIFIER) R_BRACKET alias?
  ;

column
  : (IDENTIFIER | IDENTIFIER DOT IDENTIFIER) alias?   #identifierColumn
  | function alias?                                   #functionColumn
  ;



alias
  : K_AS IDENTIFIER
  ;

function
  :  functionName L_BRACKET ( IDENTIFIER | IDENTIFIER DOT IDENTIFIER) (COMMA literalValue)* R_BRACKET         #simpleFunction
  |  functionName L_BRACKET ( IDENTIFIER | IDENTIFIER DOT IDENTIFIER) (COMMA literalValue)* R_BRACKET         #simpleFunction
  |  functionName L_BRACKET function  (COMMA literalValue)* R_BRACKET                                         #recursiveFunction
  ;

functionName
  : IDENTIFIER
  ;


whereStatement
  : K_WHERE expr                        #whereStmt
  ;

expr
 : simpleExpr                                                                    #simpleCondition
 | (compoundExpr | simpleExpr)   ( K_AND | K_OR )  (compoundExpr | simpleExpr)   #compoundRecursiveCondition
 // IDENTIFIER (K_IS_NULL | K_IS_NOT_NULL)           #nullCondition
 | expr (K_AND| K_OR) expr                                                       #simpleRecursiveCondition
 | compoundExpr                                                                  #simpleCompoundCondition
 ;

simpleExpr
 : IDENTIFIER comparisonOperator literalValue                #simpleExpression
 | IDENTIFIER DOT IDENTIFIER comparisonOperator literalValue #nestedExpression
 ;

compoundExpr
 : L_BRACKET expr R_BRACKET                                  #compoundExpression
 ;

comparisonOperator
  : K_EQUAL | K_GREATER | K_LESS | K_LESS_EQUAL | K_GREATER_EQUAL |  K_NOT_EQUAL  | K_LIKE |K_NOT_LIKE | K_IN | K_IS | K_NOT_IN
  ;


literalValue
 : NUMERIC_LITERAL
 | BOOLEAN_LITERAL
 | STRING_LITERAL
 ;


groupBy
 : K_GROUP_BY  column
 ;


SPACE
 : [ \u000B\t\r\n] -> channel(HIDDEN)
 ;

WS : [ \t]+ -> skip ;

COMMA : ',' ;
L_BRACKET : '(' ;
R_BRACKET : ')' ;
DOT : '.';

EOQ: ';';

K_SELECT : S E L E C T;
K_WHERE : W H E R E;
K_WINDOW_TUMBLING : W I N D O W SPACE T U M B L I N G;
K_GROUP_BY : G R O U P SPACE B Y;
K_AND : A N D;
K_AS : A S;
K_OR : O R;
K_IS : I S;
K_LIKE : L I K E;
K_NOT_LIKE : N O T SPACE L I K E;
K_EQUAL : '=';
K_GREATER : '>';
K_LESS : '<';
K_LESS_EQUAL : (K_LESS K_EQUAL);
K_GREATER_EQUAL : (K_GREATER K_EQUAL);
K_NOT_EQUAL : ('!' K_EQUAL);
K_NULL : N U L L;
K_IS_NULL :  (K_IS SPACE K_NULL);
K_IS_NOT_NULL : (K_IS SPACE K_NOT SPACE K_NULL);
K_NOT : N O T;
K_NOT_IN : (K_NOT SPACE I N);
K_IN : I N;
K_COUNT : C O U N T;
K_SUM : S U M;
K_MIN : M I N;
K_MAX : M A X;
K_AVG : A V G;
K_TRUE : T R U E;
K_FALSE : F A L S E;



IDENTIFIER
  : '"' (~'"' | '""')* '"'
  | '`' (~'`' | '``')* '`'
  | '[' ~']'* ']'
  | [a-zA-Z_] [a-zA-Z_0-9]*
  ;


NUMERIC_LITERAL
 : DIGIT+ ( '.' DIGIT* )? ( E [-+]? DIGIT+ )?
 | '.' DIGIT+ ( E [-+]? DIGIT+ )?
 ;

STRING_LITERAL
 : '\'' ( ~'\'' | '\'\'' )* '\''
 ;


BOOLEAN_LITERAL
 : (K_TRUE | K_FALSE)
 ;

STAR : '*';

fragment DIGIT : [0-9];
fragment A : [aA];
fragment B : [bB];
fragment C : [cC];
fragment D : [dD];
fragment E : [eE];
fragment F : [fF];
fragment G : [gG];
fragment H : [hH];
fragment I : [iI];
fragment J : [jJ];
fragment K : [kK];
fragment L : [lL];
fragment M : [mM];
fragment N : [nN];
fragment O : [oO];
fragment P : [pP];
fragment Q : [qQ];
fragment R : [rR];
fragment S : [sS];
fragment T : [tT];
fragment U : [uU];
fragment V : [vV];
fragment W : [wW];
fragment X : [xX];
fragment Y : [yY];
fragment Z : [zZ];
