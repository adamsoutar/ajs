package parser

var operatorPrecedence = map[string]int {
	"=": 3,
	""
}

type associativity bool
const (
	rightAssociative associativity = true
	leftAssociative = false
)
var operatorAssociativity = map[string]associativity {
	"=": leftAssociative,
}