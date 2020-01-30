package parser

type Parser struct {
  tokens tokenStream
  ast []astNode
}

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

func (p *Parser) parseDelimited (opening string, delim string, closing string) []astNode {
  var args []astNode
  p.expectPunctuation(opening)

  for {
    args = append(args, p.parseComponent(false))

    if !p.isNextPunctuation(delim) {
      break
    }
    // Read in the comma
    p.tokens.read()
  }

  p.expectPunctuation(closing)
  return args
}

func (p *Parser) mightBeCall (node astNode) astNode {
  // TODO: Panic if we're calling something stupid eg. 3.14("Hello")
  if p.isNextPunctuation("(") {
    return astNodeFunctionCall{
      funcName: node,
      args: p.parseDelimited("(", ",", ")"),
    }
  }
  return node
}

func (p *Parser) mightBeBinary (me astNode, myPrecedence int) astNode {
  if p.isNextOperator() {
    var op = p.tokens.peek().(operatorToken).operator

    var theirPrecedence = operatorPrecedence[op]
    var theirAssociativity = operatorAssociativity[op]
    if theirAssociativity != rightAssociative {
      panic("Wasn't expecting an with right associativity here!")
    }

    if theirPrecedence > myPrecedence {
      p.tokens.read()
      var them = p.parseComponent(false)

      var node = astNodeBinary{
        operator: op,
        left: me,
        right: them,
      }

      return p.mightBeBinary(node, myPrecedence)
    }
  }
  return me
}

func (p *Parser) mightBePropertyAccess (me astNode) astNode {
  if p.isNextOperator(".") {
    p.tokens.read()
    return astNodePropertyAccess{
      object: me,
      property: p.mightBePropertyAccess(p.parseAtom(false)),
    }
  }

  return me
}

func (p *Parser) parseAssigment (isConst bool) astNode {
  // Read in the vars
  var vars []string
  for {
    var t = p.tokens.read()
    if t.getTokenType() != tkIdentifier {
      panic("Attempted assignment to something that's not a variable. ie. let 3 = 4")
    }
    vars = append(vars, t.(identifierToken).value)

    // TODO: Implement check for ',' after punctuation is in
    break
  }

  // Expect the =
  var eq = p.tokens.read()
  if eq.getTokenType() != tkOperator || eq.(operatorToken).operator != "=" {
    panic("After let x, you must have an = operator.")
  }

  var value = p.parseComponent(false)

  return astNodeAssignment{
    value: value,
    vars: vars,
  }
}

func (p *Parser) expectToken (typ tokenType) token {
  var tk = p.tokens.read()
  var actual = tk.getTokenType()
  if actual != typ {
    panic("Expected token type " + string(int(typ)) + ", but got " + string(int(actual)))
  }
  return tk
}

func (p *Parser) parseComponent (acceptStatements bool) astNode {
  return p.mightBeBinary(
    p.mightBeCall(
      p.mightBePropertyAccess(
      p.parseAtom(acceptStatements))), 0)
}

func (p *Parser) parseAtom (acceptStatements bool) astNode {
  var t = p.tokens.read()

  // Attempt to parse an expression
  // TODO: Bracketed expressions

  switch t.getTokenType() {
  case tkNumber:
    return astNodeNumber{value:t.(numberToken).value}
  case tkString:
    return astNodeString{value:t.(stringToken).value}
  case tkIdentifier:
    return astNodeIdentifier{name:t.(identifierToken).value}
  }

  if !acceptStatements {
    panic("Attempted to use a statement somewhere where they are not allowed")
  }

  return p.parseStatement(t)
}

func (p *Parser) parseStatement (t token) astNode {
  if t.getTokenType() == tkKeyword {
    var keyword = t.(keywordToken).value
    switch keyword {
    case "let":
      return p.parseAssigment(false)
    case "const":
      return p.parseAssigment(true)
    }
  }

  panic("Unhandled AST node type in parseStatement!")
}

func (p *Parser) ParseAST () []astNode {
  var ast []astNode
  for !p.tokens.endOfStream {
    ast = append(ast, p.parseComponent(true))
    p.expectToken(tkLineTerminator)
  }
  p.ast = ast

  return p.ast
}

func New (code string) Parser {
  var codeStream = stringStream { code: code }
  var tkStream = tokenStream { code: codeStream }
  tkStream.read()
  var p = Parser { tokens: tkStream }

  return p
}
