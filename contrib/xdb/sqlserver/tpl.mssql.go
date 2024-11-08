package sqlserver

import (
	"database/sql"
	"fmt"

	"github.com/zhiyunliu/glue/contrib/xdb/tpl"
	"github.com/zhiyunliu/glue/xdb"
)

// MssqlTemplate  模板
type MssqlTemplate struct {
	name    string
	prefix  string
	matcher xdb.TemplateMatcher
}

type mssqlPlaceHolder struct {
	ctx *MssqlTemplate
}

func (ph *mssqlPlaceHolder) Get(propName string) (argName, phName string) {
	argName = fmt.Sprint(ph.ctx.prefix, propName)
	phName = ph.NamedArg(argName)
	return
}

func (ph *mssqlPlaceHolder) NamedArg(argName string) (phName string) {
	phName = "@" + argName
	return
}

func (ph *mssqlPlaceHolder) BuildArgVal(argName string, val interface{}) interface{} {
	if arg, ok := val.(sql.NamedArg); ok {
		return arg
	}
	return sql.NamedArg{Name: argName, Value: val}

}

func (ph *mssqlPlaceHolder) Clone() xdb.Placeholder {
	return &mssqlPlaceHolder{
		ctx: ph.ctx,
	}
}

func New(name, prefix string, matcher xdb.TemplateMatcher) xdb.SQLTemplate {

	if matcher == nil {
		panic(fmt.Errorf("New ,TemplateMatcher Can't be nil"))
	}
	return &MssqlTemplate{
		name:    name,
		prefix:  prefix,
		matcher: matcher,
	}
}

func (ctx *MssqlTemplate) Name() string {
	return ctx.name
}

func (ctx *MssqlTemplate) Placeholder() xdb.Placeholder {
	return &mssqlPlaceHolder{ctx: ctx}
}

// GetSQLContext 获取查询串
func (template *MssqlTemplate) GetSQLContext(sqlTpl string, input map[string]any, opts ...xdb.TemplateOption) (query string, args []any, err error) {
	return tpl.AnalyzeTPLFromCache(template, sqlTpl, input, opts...)
}

func (template *MssqlTemplate) RegistExpressionMatcher(matchers ...xdb.ExpressionMatcher) {
	template.matcher.RegistMatcher(matchers...)
}

func (template *MssqlTemplate) HandleExpr(item xdb.SqlState, sqlTpl string, input xdb.DBParam) (sql string, err error) {
	return template.matcher.GenerateSQL(item, sqlTpl, input)
}

func (template *MssqlTemplate) GetSqlState(tplOpts *xdb.TemplateOptions) xdb.SqlState {
	return xdb.NewDefaultSqlState(template.Placeholder(), tplOpts)
}
