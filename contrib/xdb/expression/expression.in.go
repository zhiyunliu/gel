package expression

import (
	"regexp"
	"sync"

	"github.com/zhiyunliu/glue/xdb"
)

var _ xdb.ExpressionMatcher = &inExpressionMatcher{}

func NewInExpressionMatcher(symbolMap xdb.SymbolMap) xdb.ExpressionMatcher {
	//in aaa
	//in t.aaa
	//t.aaa in aaa
	//bbb in aaa
	const pattern = `[&|\|](({in\s+(\w+(\.\w+)?)\s*})|({(\w+(\.\w+)?)\s+in\s+(\w+)\s*}))`
	return &inExpressionMatcher{
		regexp:          regexp.MustCompile(pattern),
		expressionCache: &sync.Map{},
		symbolMap:       symbolMap,
	}
}

type inExpressionMatcher struct {
	symbolMap       xdb.SymbolMap
	regexp          *regexp.Regexp
	expressionCache *sync.Map
}

func (m *inExpressionMatcher) Name() string {
	return "in"
}

func (m *inExpressionMatcher) Pattern() string {
	return m.regexp.String()
}

func (m *inExpressionMatcher) LoadSymbol(symbol string) (xdb.Symbol, bool) {
	return m.symbolMap.Load(symbol)
}

func (m *inExpressionMatcher) MatchString(expression string) (valuer xdb.ExpressionValuer, ok bool) {
	tmp, ok := m.expressionCache.Load(expression)
	if ok {
		valuer = tmp.(xdb.ExpressionValuer)
		return
	}

	parties := m.regexp.FindStringSubmatch(expression)
	if len(parties) <= 0 {
		return
	}
	ok = true

	var (
		item = &xdb.ExpressionItem{
			Symbol: getExpressionSymbol(expression),
			Oper:   m.Name(),
		}
		fullField string
		propName  string
	)
	// fullfield,oper,oper
	//&{in tbl.field} => 3,in,prop(3)
	//&{tt.field  in    property} => 6,in, 8

	if parties[3] != "" {
		fullField = parties[3]
		propName = getExpressionPropertyName(fullField)
	} else {
		fullField = parties[6]
		propName = parties[8]
	}

	item.FullField = fullField
	item.PropName = propName

	item.ExpressionBuildCallback = m.buildCallback()
	m.expressionCache.Store(expression, item)

	return item, ok
}

func (m *inExpressionMatcher) buildCallback() xdb.ExpressionBuildCallback {
	return func(item *xdb.ExpressionItem, param xdb.DBParam, argName string) (expression string, err xdb.MissError) {

		return
	}
}
