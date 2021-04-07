package model

import (
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
)

// Directive ...
func (o *ObjectField) Directive(name string) *ast.Directive {
	for _, d := range o.Def.Directives {
		if d.Name.Value == name {
			return d
		}
	}
	return nil
}

// HasDirective ...
func (o *ObjectField) HasDirective(name string) bool {
	return o.Directive(name) != nil
}

// ColumnType ...
func (o *ObjectField) ColumnType() (value string) {
	directive := o.Directive("column")
	if directive == nil {
		return
	}
	for _, arg := range directive.Arguments {
		if arg.Name.Value == "type" {
			val := arg.Value.GetValue()
			value, _ = val.(string)
			break
		}
	}
	return
}

// ModelTags ...
func (o *ObjectField) ModelTags() string {
	_gorm := fmt.Sprintf("column:%s", o.Name())
	if o.Name() == "id" {
		_gorm += ";primary_key"
	}

	if o.IsEmbeddedColumn() {
		_gorm += ";type:text"
	} else {

		columnDirective := o.Directive("column")
		for _, arg := range columnDirective.Arguments {
			if arg.Name.Value == "type" {
				_gorm += fmt.Sprintf(";type:%v", arg.Value.GetValue())
			}
			if arg.Name.Value == "unique" {
				val, ok := arg.Value.GetValue().(bool)
				if ok && val {
					_gorm += fmt.Sprintf(";unique")
				}
			}
			if arg.Name.Value == "index" {
				_gorm += fmt.Sprintf(";index:%v", arg.Value.GetValue())
			}
			if arg.Name.Value == "default" {
				_gorm += fmt.Sprintf(";default:%v", arg.Value.GetValue())
			}
		}
	}

	return fmt.Sprintf(`json:"%s" gorm:"%s"`, o.Name(), _gorm)
}
