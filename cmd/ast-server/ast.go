package main

import (
	"strconv"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
)

func translateAST(node promql.Expr) interface{} {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *promql.AggregateExpr:
		return map[string]interface{}{
			"type":     "aggregation",
			"op":       n.Op.String(),
			"expr":     translateAST(n.Expr),
			"param":    translateAST(n.Param),
			"grouping": sanitizeList(n.Grouping),
			"without":  n.Without,
		}
	case *promql.BinaryExpr:
		var matching interface{}
		if m := n.VectorMatching; m != nil {
			matching = map[string]interface{}{
				"card":    m.Card.String(),
				"labels":  sanitizeList(m.MatchingLabels),
				"on":      m.On,
				"include": sanitizeList(m.Include),
			}
		}

		return map[string]interface{}{
			"type":     "binaryExpr",
			"op":       n.Op.String(),
			"lhs":      translateAST(n.LHS),
			"rhs":      translateAST(n.RHS),
			"matching": matching,
			"bool":     n.ReturnBool,
		}
	case *promql.Call:
		args := []interface{}{}
		for _, arg := range n.Args {
			args = append(args, translateAST(arg))
		}

		return map[string]interface{}{
			"type": "call",
			"func": map[string]interface{}{
				"name":       n.Func.Name,
				"argTypes":   n.Func.ArgTypes,
				"variadic":   n.Func.Variadic,
				"returnType": n.Func.ReturnType,
			},
			"args": args,
		}
	case *promql.MatrixSelector:
		return map[string]interface{}{
			"type":     "matrixSelector",
			"name":     n.Name,
			"range":    n.Range.Milliseconds(),
			"offset":   n.Offset.Milliseconds(),
			"matchers": translateMatchers(n.LabelMatchers),
		}
	case *promql.SubqueryExpr:
		return map[string]interface{}{
			"type":   "subquery",
			"expr":   translateAST(n.Expr),
			"range":  n.Range.Milliseconds(),
			"offset": n.Offset.Milliseconds(),
			"step":   n.Step.Milliseconds(),
		}
	case *promql.NumberLiteral:
		return map[string]string{
			"type": "numberLiteral",
			"val":  strconv.FormatFloat(n.Val, 'f', -1, 64),
		}
	case *promql.ParenExpr:
		return map[string]interface{}{
			"type": "parenExpr",
			"expr": translateAST(n.Expr),
		}
	case *promql.StringLiteral:
		return map[string]interface{}{
			"type": "stringLiteral",
			"val":  n.Val,
		}
	case *promql.UnaryExpr:
		return map[string]interface{}{
			"type": "unaryExpr",
			"op":   n.Op.String(),
			"expr": translateAST(n.Expr),
		}
	case *promql.VectorSelector:
		return map[string]interface{}{
			"type":     "vectorSelector",
			"name":     n.Name,
			"offset":   n.Offset.Milliseconds(),
			"matchers": translateMatchers(n.LabelMatchers),
		}
	}
	panic("unsupported node type")
}

func sanitizeList(l []string) []string {
	if l == nil {
		return []string{}
	}
	return l
}

func translateMatchers(in []*labels.Matcher) interface{} {
	out := []map[string]interface{}{}
	for _, m := range in {
		out = append(out, map[string]interface{}{
			"name":  m.Name,
			"value": m.Value,
			"type":  m.Type.String(),
		})
	}
	return out
}
