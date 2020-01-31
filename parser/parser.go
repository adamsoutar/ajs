package parser

type Parser struct {
  tokens tokenStream
  ast []astNode
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
  if p.isNextPunctuation("(") {
    // NOTE: This doesn't stop you calling something stupid using property access,
    //       that's for the interpreter. But it handles what you can from the parser side
    if node.getNodeType() != astIdentifier && node.getNodeType() != astPropertyAccess {
      panic("Called something unreasonable (eg. 3.14())")
    }

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
  // TODO: Support const x = 3, y = 2, z = 1 syntax

  // Read in the vars
  var t = p.tokens.read()
  if t.getTokenType() != tkIdentifier {
    panic("Attempted assignment to something that's not a variable. ie. let 3 = 4")
  }

  p.expectOperator("=")

  var value = p.parseComponent(false)

  return astNodeAssignment{
    value: value,
    varNm: t.(identifierToken).value,
  }
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
  if t.getTokenType() == tkPunctuation {
    // Bracketed expressions
    if t.(punctuationToken).punctuation == "(" {
      var exp = p.parseComponent(false)
      p.expectPunctuation(")")
      return exp
    }
  }

  if t.getTokenType() == tkLineTerminator {
    return astNodeEmptyStatement{}
  }

  switch t.getTokenType() {
  case tkNumber:
    return astNodeNumber{value:t.(numberToken).value}
  case tkString:
    return astNodeString{value:t.(stringToken).value}
  case tkIdentifier:
    return astNodeIdentifier{name:t.(identifierToken).value}
  }

  // Non-statement keyword expressions
  // Like function
  if t.getTokenType() == tkKeyword {
    var keyword = t.(keywordToken).value
    if keyword == "function" {
      return p.parseFunctionDefinition()
    }
  }

  if !acceptStatements {
    // TODO: If expecting an expression and you see {, parse an object here
  }

  if !acceptStatements {
    panic("Attempted to use a statement somewhere where they are not allowed")
  }

  return p.parseStatement(t)
}

func (p *Parser) parseFunctionDefinition () astNode {
  // Function keyword has already been consumed by parseAtom

  var nameTkn = p.tokens.read()
  if nameTkn.getTokenType() != tkIdentifier {
    panic("Function name must be an identifier")
  }
  var name = nameTkn.(identifierToken).value

  p.expectPunctuation("(")
  var params []string
  if !p.isNextPunctuation(")") {
    // This function takes parameters, let's get them
    // TODO: Default parameters
    for {
      var nextParam = p.tokens.read()
      if nextParam.getTokenType() != tkIdentifier {
        panic("Function parameters must be identifiers")
      }

      params = append(params, nextParam.(identifierToken).value)

      if !p.isNextPunctuation(",") {
        break
      }
      // Consume the comma
      p.tokens.read()
    }
  }
  p.expectPunctuation(")")

  var body = p.parseBlockStatement(true)

  return astNodeFunctionDefinition{
    name: name,
    params: params,
    body: body,
  }
}

func (p *Parser) parseBlockStatement (expectBraces bool) astNodeBlock {
  if expectBraces {
    p.expectPunctuation("{")
  }

  // Parse statements
  var stmts []astNode
  for !p.tokens.endOfStream {
    stmts = append(stmts, p.parseComponent(true))
    // TODO: Find out how we need to enforce line seperation.
    //p.expectToken(tkLineTerminator)

    if expectBraces && p.isNextPunctuation("}") {
      break
    }
  }

  if expectBraces {
    p.expectPunctuation("}")
  }

  return astNodeBlock{
    nodes: stmts,
  }
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

  if t.getTokenType() == tkPunctuation {
    if t.(punctuationToken).punctuation == "{" {
      return p.parseBlockStatement(true)
    }
  }

  panic("Unhandled AST node type in parseStatement!")
}

func (p *Parser) ParseAST () []astNode {
  p.ast = p.parseBlockStatement(false).nodes

  return p.ast
}

func New (code string) Parser {
  var codeStream = stringStream { code: code }
  var tkStream = tokenStream { code: codeStream }
  tkStream.read()
  var p = Parser { tokens: tkStream }

  return p
}
