grammar sql;



sqlQuery
  : selectQuery EOF;

selectQuery
  : K_SELECT (resultColumns)  ( whereStatement )?  EOQ
  ;


resultColumns
 : column (COMMA column)+
 | avg (COMMA avg)+
 | STAR
 ;


column
  : IDENTIFIER
  ;

whereStatement
  : ( K_WHERE expr )?
  | K_WINDOW_TUMBLING NUMERIC_LITERAL( K_WHERE expr )? (groupBy)
  ;

groupBy
 : K_GROUP_BY column
 ;


avg
  : (K_MIN | K_MAX | K_COUNT | K_AVG) L_BRACKET column R_BRACKET;

SPACE
 : [ \u000B\t\r\n] -> channel(HIDDEN)
 ;

COMMA : ',' ;
L_BRACKET : '(' ;
R_BRACKET : ')' ;

EOQ: ';';

K_SELECT : S E L E C T;
K_WHERE : W H E R E;
K_WINDOW_TUMBLING : W I N D O W SPACE T U M B L I N G;
K_GROUP_BY : G R O U P SPACE B Y;
K_AND : A N D;
K_OR : O R;
K_IS : I S;
K_LIKE : L I K E;
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
K_MIN : M I N;
K_MAX : M A X;
K_AVG : A V G;

IDENTIFIER
  : '"' (~'"' | '""')* '"'
  | '`' (~'`' | '``')* '`'
  | '[' ~']'* ']'
  | [a-zA-Z_] [a-zA-Z_0-9]*
  ;

expr
 : IDENTIFIER comparisonOperator literalValue
 | expr ( K_AND | K_OR ) expr
 | IDENTIFIER (K_IS_NULL | K_IS_NOT_NULL)
 ;



comparisonOperator
  : K_EQUAL | K_GREATER | K_LESS | K_LESS_EQUAL | K_GREATER_EQUAL |  K_NOT_EQUAL  | K_LIKE | K_IN | K_IS | K_NOT_IN
  ;


literalValue
 : NUMERIC_LITERAL
 | STRING_LITERAL
 ;


NUMERIC_LITERAL
 : DIGIT+ ( '.' DIGIT* )? ( E [-+]? DIGIT+ )?
 | '.' DIGIT+ ( E [-+]? DIGIT+ )?
 ;

STRING_LITERAL
 : '\'' ( ~'\'' | '\'\'' )* '\''
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
