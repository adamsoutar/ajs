package parser

// TODO: Wordy operators like instanceof
// This is only for binary operators
var operatorPrecedence = map[string]int {
	".": 0, // Not a real binary operator,
			// Still here for the tokeniser
	"+": 13,
	"-": 13,
	"*": 14,
	"**": 14,
	"%": 14,
	"/": 14,
	"==": 10,
	"!=": 10,
	"===": 10,
	"!==": 10,
	"&": 9,
	"^": 8,
	"|": 7,
	"&&": 6,
	"||": 5,
	">=": 11,
	">": 11,
	"<=": 11,
	"<": 11,
	">>>": 12,
	">>": 12,
	"<<": 12,
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