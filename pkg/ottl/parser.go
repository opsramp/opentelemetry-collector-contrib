// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ottl // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"

import (
	"github.com/alecthomas/participle/v2"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/multierr"
)

type Parser[K any] struct {
	functions         map[string]interface{}
	pathParser        PathExpressionParser[K]
	enumParser        EnumParser
	telemetrySettings component.TelemetrySettings
}

// Statement holds a top level statement for processing telemetry data.
type Statement[K any] struct {
	function  ExprFunc[K]
	condition boolExpressionEvaluator[K]
}

// Execute is a function that will execute the statement's function if the statement's condition is met.
// Returns true if the function was run, returns false otherwise.
// If the statement contains no condition, the function will run and true will be returned.
// In addition, the functions return value is always returned.
func (s *Statement[K]) Execute(ctx K) (any, bool, error) {
	condition, err := s.condition(ctx)
	if err != nil {
		return nil, false, err
	}
	var result any
	if condition {
		result, err = s.function(ctx)
		if err != nil {
			return nil, true, err
		}
	}
	return result, condition, nil
}

func NewParser[K any](functions map[string]interface{}, pathParser PathExpressionParser[K], enumParser EnumParser, telemetrySettings component.TelemetrySettings) Parser[K] {
	return Parser[K]{
		functions:         functions,
		pathParser:        pathParser,
		enumParser:        enumParser,
		telemetrySettings: telemetrySettings,
	}
}

func (p *Parser[K]) ParseStatements(statements []string) ([]*Statement[K], error) {
	var parsedStatements []*Statement[K]
	var errors error

	for _, statement := range statements {
		parsed, err := parseStatement(statement)
		if err != nil {
			errors = multierr.Append(errors, err)
			continue
		}
		function, err := p.newFunctionCall(parsed.Invocation)
		if err != nil {
			errors = multierr.Append(errors, err)
			continue
		}
		expression, err := p.newBooleanExpressionEvaluator(parsed.WhereClause)
		if err != nil {
			errors = multierr.Append(errors, err)
			continue
		}
		parsedStatements = append(parsedStatements, &Statement[K]{
			function:  function,
			condition: expression,
		})
	}

	if errors != nil {
		return nil, errors
	}
	return parsedStatements, nil
}

var parser = newParser()

func parseStatement(raw string) (*parsedStatement, error) {
	parsed, err := parser.ParseString("", raw)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

// newParser returns a parser that can be used to read a string into a parsedStatement. An error will be returned if the string
// is not formatted for the DSL.
func newParser() *participle.Parser[parsedStatement] {
	lex := buildLexer()
	parser, err := participle.Build[parsedStatement](
		participle.Lexer(lex),
		participle.Unquote("String"),
		participle.Elide("whitespace"),
	)
	if err != nil {
		panic("Unable to initialize parser; this is a programming error in the transformprocessor:" + err.Error())
	}
	return parser
}
