package expr

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

const identifier_nil = "nil"

var opMapping = map[string]string{
	">":  "$gt",
	">=": "$gte",
	"<":  "$lt",
	"<=": "$lte",
	"==": "$eq", // Or just use the value directly without an operator
	"!=": "$ne",
	"&&": "$and",
	"||": "$or",
}

type Analyzer struct {
	Op    string
	Left  interface{}
	Right interface{}
}

type ConstAnalyzer struct {
	Value interface{}
}

type FuncAnalyzer struct {
	Name string
	Args []interface{}
}
type FieldAnalyzer struct {
	Name string
}
type MongoExpression map[string]interface{}

func (m MongoExpression) String() string {
	if m == nil {
		return "null"
	}
	jsonData, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return fmt.Sprintf("error marshaling JSON: %v", err)
	}
	return string(jsonData)
}
func selectorExprToString(sel *ast.SelectorExpr) string {
	var result string

	// Recursively build the string
	x := sel.X
	for {
		switch v := x.(type) {
		case *ast.Ident:
			result = v.Name + result // Prepend the identifier
			return result
		case *ast.SelectorExpr:
			result = "." + sel.Sel.Name + result
			sel = v
			x = sel.X
		default:
			// Handle other expression types if needed (e.g., IndexExpr, CallExpr)
			return "<unsupported expression type>"
		}
	}
}
func analyzeNode(fset *token.FileSet, node ast.Node) interface{} {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		op := n.Op.String()
		left := analyzeNode(fset, n.X)
		right := analyzeNode(fset, n.Y)
		if _, ok := right.(FuncAnalyzer); ok {
			// right is of type FuncAnalyzer
			// You can now use right as a FuncAnalyzer
			// Example:
			analyzer := right.(FuncAnalyzer)
			// Use analyzer methods
			if analyzer.Name == "IsNull" {
				return Analyzer{Op: op, Left: left, Right: nil} // Use equality operator for IsNull function
			}
		}

		return Analyzer{Op: op, Left: left, Right: right}

	case *ast.UnaryExpr:
		op := n.Op.String()
		operand := analyzeNode(fset, n.X)
		return Analyzer{Op: op, Left: operand, Right: nil}

	case *ast.SelectorExpr:
		x, ok := n.X.(*ast.Ident)
		if !ok {
			var buf strings.Builder
			printer.Fprint(&buf, fset, n.X)
			fName := buf.String() + "." + n.Sel.Name
			//return fmt.Sprintf("Unsupported selector base: %s", buf.String())
			// x.Name = selectorExprToString(n) // Use the string representation of the selector
			return FieldAnalyzer{Name: fName}
		}
		sel := n.Sel.Name
		if strings.ToLower(x.Name) == "id" {
			x.Name = "_id"
		}
		return FieldAnalyzer{Name: x.Name + "." + sel} // Now a FieldAnalyzer

	case *ast.Ident:
		if strings.ToLower(n.Name) == "id" {
			return FieldAnalyzer{Name: "_id"} // Now a FieldAnalyzer
		}
		if n.Name == identifier_nil {
			return FuncAnalyzer{Name: "IsNull"} // Return nil for nil identifierreturn FieldAnalyzer{Name: n.Name}
		}
		return FieldAnalyzer{Name: n.Name} // Now a FieldAnalyzer

	case *ast.BasicLit:
		switch n.Kind {
		case token.STRING:
			val, _ := strconv.Unquote(n.Value)
			return ConstAnalyzer{Value: val}
		case token.INT:
			val, _ := strconv.Atoi(n.Value)
			return ConstAnalyzer{Value: val}
		case token.FLOAT:
			val, _ := strconv.ParseFloat(n.Value, 64)
			return ConstAnalyzer{Value: val}
		default:
			return ConstAnalyzer{Value: n.Value}
		}
	case *ast.ParenExpr:
		return analyzeNode(fset, n.X)
	case *ast.CallExpr:
		fun := analyzeNode(fset, n.Fun)
		args := make([]interface{}, len(n.Args))
		for i, arg := range n.Args {
			args[i] = analyzeNode(fset, arg)
		}
		return FuncAnalyzer{Name: fun.(FieldAnalyzer).Name, Args: args} //Use FieldAnalyzer.Name
	default:
		var buf strings.Builder
		printer.Fprint(&buf, fset, n)
		return fmt.Sprintf("Unsupported node type: %s", buf.String())
	}
}

func AnalyzeFunction(expr string) (interface{}, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseExprFrom(fset, "", expr, 0)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}

	return analyzeNode(fset, node), nil
}
func AnalyzeExpressionWithPlaceholders(expressionTemplate string, args ...interface{}) (interface{}, error) {
	// Replace placeholders
	var expr strings.Builder
	parts := strings.Split(expressionTemplate, "?")
	if len(parts)-1 != len(args) {
		return nil, fmt.Errorf("number of placeholders does not match number of arguments")
	}

	for i, part := range parts {
		expr.WriteString(part)
		if i < len(args) {
			switch v := args[i].(type) {
			case string:
				expr.WriteString("\"" + v + "\"")
			case int:
				expr.WriteString(strconv.Itoa(v))
			case float64:
				expr.WriteString(strconv.FormatFloat(v, 'G', 10, 64))
			case nil:
				expr.WriteString(identifier_nil)
			default:
				return nil, fmt.Errorf("unsupported argument type: %T", v)
			}
		}
	}

	// Analyze the resulting expression
	return AnalyzeFunction(expr.String())
}
func builEq(left interface{}, right interface{}, originalExpr string) (bson.D, error) {
	if left == nil {
		return nil, fmt.Errorf("left operand cannot be nil")
	}
	switch leftType := left.(type) {
	case FieldAnalyzer:
		if right == nil {
			checkNotExist := bson.D{{leftType.Name, bson.D{{"$exists", false}}}}
			checkNull := bson.D{{leftType.Name, "null"}}
			return bson.D{{"$or", bson.A{checkNotExist, checkNull}}}, nil

		}
		switch rightType := right.(type) {
		case ConstAnalyzer:
			return bson.D{{leftType.Name, rightType.Value}}, nil
		default:
			return nil, fmt.Errorf("unsupported right operand type: %T parse from %s", right, originalExpr)
		}

	default:
		// Handle unexpected types (e.g., return an error or default value)
		return nil, fmt.Errorf("unsupported expression type: %T parse from %s", leftType, originalExpr)
	}
}

func convertToMongoExpr(analyExpr interface{}, originalExpr string) (interface{}, error) {
	switch expr := analyExpr.(type) {
	case Analyzer:
		left, err := convertToMongoExpr(expr.Left, originalExpr)
		if err != nil {
			return nil, err
		}
		right, err := convertToMongoExpr(expr.Right, originalExpr)
		if err != nil {
			return nil, err
		}
		if expr.Op == "==" {
			return builEq(left, right, originalExpr)
		}
		// get the operator
		if op, ok := opMapping[expr.Op]; ok {
			return bson.D{{op, bson.A{left, right}}}, nil
		} else {
			return nil, fmt.Errorf("unsupported operator: %s parse from %s", expr.Op, originalExpr)
		}
	case ConstAnalyzer:
		return expr, nil
	case nil:
		return nil, nil
	case FieldAnalyzer:
		return expr, nil
	case FuncAnalyzer:
		return buildWithFunc(expr, originalExpr)
	default:
		// Handle unexpected types (e.g., return an error or default value)
		return nil, fmt.Errorf("unsupported expression type: %T, parse from %s", expr, originalExpr)
	}
}

func GetMongoQueryFromString(expr string, args ...interface{}) (bson.D, error) {
	analyExpr, err := AnalyzeExpressionWithPlaceholders(expr, args...)
	if err != nil {
		return nil, err
	}
	ret, e := convertToMongoExpr(analyExpr, expr)
	if e != nil {
		return nil, e
	}
	mongoExpr, ok := ret.(bson.D)
	if !ok {
		return nil, fmt.Errorf("type assertion failed: expected bson.D, got %T", ret)
	}
	return mongoExpr, nil
}

func ToPrettyJSON(data interface{}) string {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Sprintf("error marshaling JSON: %v", err)
	}
	return string(jsonData)
}
func ToPrettyJSONOfBSON(d bson.D) string {
	var m map[string]interface{}
	data, err := bson.Marshal(d)
	if err != nil {
		return "error marshaling BSON"
	}

	err = bson.Unmarshal(data, &m)
	if err != nil {
		return "error unmarshaling BSON to map"
	}
	r := ToPrettyJSON(m)
	return r
}
func buildIsNull(f FuncAnalyzer, originalExpr string) (bson.D, error) {
	if len(f.Args) != 1 {
		return nil, fmt.Errorf("function 'IsNull' expects exactly one argument, parse from %s", originalExpr)
	}
	arg, err := convertToMongoExpr(f.Args[0], originalExpr)
	if err != nil {
		return nil, err
	}
	switch argType := arg.(type) {
	case FieldAnalyzer:
		checjNoyExist := bson.D{{argType.Name, bson.D{{"$exists", false}}}}
		checkNull := bson.D{{argType.Name, "null"}}
		return bson.D{{"$or", bson.A{checjNoyExist, checkNull}}}, nil
	default:
		return nil, fmt.Errorf("unsupported argument type: %T, parse from %s", argType, originalExpr)
	}
}

func buildContains(f FuncAnalyzer, originalExpr string) (bson.D, error) {
	if len(f.Args) != 2 {
		return nil, fmt.Errorf("function 'Like' expects exactly two arguments, parse from %s", originalExpr)
	}
	left, err := convertToMongoExpr(f.Args[0], originalExpr)
	if err != nil {
		return nil, err
	}
	right, err := convertToMongoExpr(f.Args[1], originalExpr)
	if err != nil {
		return nil, err
	}
	switch leftType := left.(type) {
	case FieldAnalyzer:
		switch rightType := right.(type) {
		case ConstAnalyzer:
			return bson.D{{leftType.Name, bson.D{{"$regex", right.(ConstAnalyzer).Value.(string)}}}}, nil
		default:
			return nil, fmt.Errorf("unsupported right operand type: %T, parse from %s", rightType, originalExpr)

		}
	default:
		return nil, fmt.Errorf("unsupported left operand type: %T, parse from %s", leftType, originalExpr)
	}
}
func buildStartsWith(f FuncAnalyzer, originalExpr string) (bson.D, error) {
	if len(f.Args) != 2 {
		return nil, fmt.Errorf("function 'StartsWith' expects exactly two arguments, parse from %s", originalExpr)
	}
	left, err := convertToMongoExpr(f.Args[0], originalExpr)
	if err != nil {
		return nil, err
	}
	right, err := convertToMongoExpr(f.Args[1], originalExpr)
	if err != nil {
		return nil, err
	}
	switch leftType := left.(type) {
	case FieldAnalyzer:
		switch rightType := right.(type) {
		case ConstAnalyzer:
			return bson.D{{leftType.Name, bson.D{{"$regex", "^" + right.(ConstAnalyzer).Value.(string)}}}}, nil
		default:
			return nil, fmt.Errorf("unsupported right operand type: %T, parse from %s", rightType, originalExpr)

		}
	default:
		return nil, fmt.Errorf("unsupported left operand type: %T, parse from %s", leftType, originalExpr)
	}
}
func buildEndsWith(f FuncAnalyzer, originalExpr string) (bson.D, error) {
	if len(f.Args) != 2 {
		return nil, fmt.Errorf("function 'EndsWith' expects exactly two arguments, parse from %s", originalExpr)
	}
	left, err := convertToMongoExpr(f.Args[0], originalExpr)
	if err != nil {
		return nil, err
	}
	right, err := convertToMongoExpr(f.Args[1], originalExpr)
	if err != nil {
		return nil, err
	}
	switch leftType := left.(type) {
	case FieldAnalyzer:
		switch rightType := right.(type) {
		case ConstAnalyzer:
			return bson.D{{leftType.Name, bson.D{{"$regex", right.(ConstAnalyzer).Value.(string) + "$"}}}}, nil
		default:
			return nil, fmt.Errorf("unsupported right operand type: %T, parse from %s", rightType, originalExpr)

		}
	default:
		return nil, fmt.Errorf("unsupported left operand type: %T, parse from %s", leftType, originalExpr)
	}
}
func buildWithFunc(f FuncAnalyzer, originalExpr string) (bson.D, error) {

	switch strings.ToLower(f.Name) {
	case "isnull":
		return buildIsNull(f, originalExpr)
	case "contains":
		return buildContains(f, originalExpr)
	case "startswith":
		return buildStartsWith(f, originalExpr)
	case "endswith":
		return buildEndsWith(f, originalExpr)
	default:
		return nil, fmt.Errorf("unsupported function: %s, parse from %s", f.Name, originalExpr)

	}
}
