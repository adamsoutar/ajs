package parser

func (p *Parser) isNextOperator (params ...string) bool {
	var opNext = p.tokens.peek()
	var isOp = opNext.getTokenType() == tkOperator
	// Optional arg specifies 'Is there an operator next' vs 'Is THIS operator next?'
	if len(params) == 0 {
		return isOp
	}

	if len(params) > 1 {
		panic("Too many args for isNextOperator")
	}

	return isOp && opNext.(operatorToken).operator == params[0]
}

func (p *Parser) expectPunctuation (punc string) {
	var t = p.tokens.read()
	if t.getTokenType() != tkPunctuation || t.(punctuationToken).punctuation != punc {
		panic("Expected punctuation \"" + punc + "\"")
	}
}

func (p *Parser) isNextPunctuation (punc string) bool {
	var t = p.tokens.peek()
	if t.getTokenType() != tkPunctuation || t.(punctuationToken).punctuation != punc {
		return false
	}
	return true
}

func (p *Parser) expectToken (typ tokenType) token {
	var tk = p.tokens.read()
	var actual = tk.getTokenType()
	if actual != typ {
		panic("Expected token type " + string(int(typ)) + ", but got " + string(int(actual)))
	}
	return tk
}

func (p *Parser) expectOperator (op string) {
	var t = p.tokens.read()
	if t.getTokenType() != tkOperator || t.(operatorToken).operator != op {
		panic("Expected punctuation \"" + op + "\"")
	}
}