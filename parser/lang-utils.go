package parser

// This is only for binary operators
var operatorPrecedence = map[string]int {
	"+": 13,
	"-": 13,
	"*": 14,
	"/": 14,
	".": 18,
}
var assignmentOperators = []string {
	"=",
	"+=",
	"-=",
	"*=",
	"/=",
	"%=",
	"**=",
	"<<=",
	">>=",
	">>>=",
	"&=",
	"^=",
	"|=",
}

// Called once from tokeniser-utils
func getOperators () []string {
	var keys = make([]string, len(operatorPrecedence) + len(assignmentOperators))
	var i = 0
	for k := range operatorPrecedence {
		keys[i] = k
		i++
	}
	for j, op := range assignmentOperators {
		keys[j + len(operatorPrecedence)] = op
	}
	return keys
}