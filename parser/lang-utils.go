package parser

var operatorPrecedence = map[string]int {
	"=": 3,
	"+": 13,
	"-": 13,
	"*": 14,
	"/": 14,
	".": 18,
}
var assignmentOperators = []string {
	"=",
}

type associativity bool
const (
	rightAssociative associativity = true
	leftAssociative = false
)
var operatorAssociativity = map[string]associativity {
	"=": leftAssociative,
	"+": rightAssociative,
	"-": rightAssociative,
	"*": rightAssociative,
	"/": rightAssociative,
	".": rightAssociative,
}

// Called once from tokeniser-utils
func getOperators () []string {
	var keys = make([]string, len(operatorPrecedence))
	var i = 0
	for k := range operatorPrecedence {
		keys[i] = k
		i++
	}
	return keys
}