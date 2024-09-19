package assert

import (
	"context"
	"reflect"
	"regexp"
	"sync"

	"github.com/jmespath-community/go-jmespath/pkg/parsing"
	reflectutils "github.com/kyverno/kyverno-json/pkg/utils/reflect"
)

var (
	foreachRegex = regexp.MustCompile(`^~(\w+)?\.(.*)`)
	bindingRegex = regexp.MustCompile(`(.*)\s*->\s*(\w+)$`)
	escapeRegex  = regexp.MustCompile(`^\\(.+)\\$`)
	engineRegex  = regexp.MustCompile(`^\((?:(\w+);)?(.+)\)$`)
)

type expression struct {
	foreach     bool
	foreachName string
	statement   string
	binding     string
	engine      string
	ast         func() (parsing.ASTNode, error)
}

func parseExpressionRegex(_ context.Context, in string) *expression {
	expression := &expression{}
	// 1. match foreach
	if match := foreachRegex.FindStringSubmatch(in); match != nil {
		expression.foreach = true
		expression.foreachName = match[1]
		in = match[2]
	}
	// 2. match binding
	if match := bindingRegex.FindStringSubmatch(in); match != nil {
		expression.binding = match[2]
		in = match[1]
	}
	// 3. match escape, if there's no escaping then match engine
	if match := escapeRegex.FindStringSubmatch(in); match != nil {
		in = match[1]
	} else {
		if match := engineRegex.FindStringSubmatch(in); match != nil {
			expression.engine = match[1]
			// account for default engine
			if expression.engine == "" {
				expression.engine = "jp"
			}
			in = match[2]
		}
	}
	// parse statement
	expression.statement = in
	if expression.statement == "" {
		return nil
	}
	return expression
}

func parseExpression(ctx context.Context, value any) *expression {
	if reflectutils.GetKind(value) != reflect.String {
		return nil
	}
	exp := parseExpressionRegex(ctx, reflect.ValueOf(value).String())
	if exp != nil && exp.engine == "jp" {
		exp.ast = sync.OnceValues(func() (parsing.ASTNode, error) {
			parser := parsing.NewParser()
			return parser.Parse(exp.statement)
		})
	}
	return exp
}
