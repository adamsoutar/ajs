package parser

type Parser struct {
  tokens tokenStream
  ast []astNode
}

func (p *Parser) parseDelimited (opening string, delim string, closing string) []astNode {
  var args []astNode
  p.expectPunctuation(opening)

  // Functions etc. can have no args
  if p.isNextPunctuation(closing) {
    p.tokens.read()
    return []astNode {}
  }

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

func (p *Parser) mightBeCall (node astNode) (bool, astNode) {
  if p.isNextPunctuation("(") {
    // NOTE: This doesn't stop you calling something stupid using property access,
    //       that's for the interpreter. But it handles what you can from the parser side
    var validCalls = []astType { astIdentifier, astPropertyAccess, astFunctionCall }
    if !inAstTypeArray(node.getNodeType(), validCalls) {
      panic("Called something unreasonable (eg. 3.14())")
    }

    return true, astNodeFunctionCall{
      funcName: node,
      args: p.parseDelimited("(", ",", ")"),
    }
  }
  return false, node
}

func (p *Parser) mightBeBinary (me astNode, myPrecedence int) astNode {
  if p.isNextOperator() {
    var op = p.tokens.peek().(operatorToken).operator

    var theirPrecedence = operatorPrecedence[op]

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

func (p *Parser) mightBePropertyAccess (me astNode) (bool, astNode) {
  if p.isNextOperator(".") {
    p.tokens.read()
    return true, astNodePropertyAccess{
      object: me,
      property: p.parseAtom(false),
    }
  }

  return false, me
}

func (p *Parser) mightBeAssignment (me astNode) astNode {
  var t = p.tokens.peek()

  if t.getTokenType() != tkOperator {
    return me
  }

  var op = t.(operatorToken).operator
  if !inStringArray(op, assignmentOperators) {
    return me
  }

  // It's an assignment
  p.tokens.read()

  var validToAssignTo = []astType{ astIdentifier, astPropertyAccess }
  var meType = me.getNodeType()
  if !inAstTypeArray(meType, validToAssignTo) {
    panic("Invalid left side to assignment (eg. 3 = 4)")
  }

  var assigningValue = p.parseComponent(false)

  return astNodeAssignment{
    variable: me,
    value:    assigningValue,
    operator: op,
  }
}

func (p *Parser) parseComponent (acceptStatements bool) astNode {
  var nd = p.parseAtom(acceptStatements)
  for {
    var changed = false
    var wasAccess, newNode1 = p.mightBePropertyAccess(nd)
    changed = changed || wasAccess
    nd = newNode1

    var wasCall, newNode2 = p.mightBeCall(nd)
    changed = changed || wasCall
    nd = newNode2

    if !changed {
      break
    }
  }

  var final = p.mightBeBinary(p.mightBeAssignment(nd), 0)
  return final
}

func (p *Parser) parseObject () astNode {
  var obj = make(map[astNode]astNode)
  for {
    var key = p.parseAtom(false)
    p.expectPunctuation(":")
    var value = p.parseComponent(false)

    obj[key] = value

    if !p.isNextPunctuation(",") {
      break
    }
    // Read in the comma
    p.tokens.read()
  }

  p.expectPunctuation("}")

  return astNodeObject {
    valueMap: obj,
  }
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
  case tkBoolean:
    return astNodeBoolean{value:t.(booleanToken).value}
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
    if t.getTokenType() == tkPunctuation &&
      t.(punctuationToken).punctuation == "{" {
      return p.parseObject()
    }
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

  var body = p.parseBlockStatement(true, true)

  return astNodeFunctionDefinition{
    name: name,
    params: params,
    body: body,
  }
}

// astNodeVariableDeclaration is basically an assignment with loads of verification
func (p *Parser) parseVariableDeclaration (isConstant bool, isHoisted bool) astNode {
  var nxt = p.parseComponent(false)

  if isConstant && isHoisted {
    panic("Variable declaration is both constant (const) *and* hoisted (var)? Mistake in parser?")
  }

  // TODO: Allow an undefined var dec (eg. let x;)
  if nxt.getNodeType() != astAssignment {
    panic("Variable declaration not followed by assignment")
  }

  var assign = nxt.(astNodeAssignment)
  if assign.operator != "=" {
    panic("Variable declaration with non = operator (eg. let x *= 3)")
  }

  if assign.variable.getNodeType() != astIdentifier {
    panic("Variable declaration of non-identifier. (eg. let func() = 3)")
  }
  var ident = assign.variable.(astNodeIdentifier)

  return astNodeVariableDeclaration {
    varName: ident.name,
    value: assign.value,
    isConstant: isConstant,
    isHoisted: isHoisted,
  }
}

func (p *Parser) parseBlockStatement (expectBraces bool, ignoreFirst bool) astNodeBlock {
  if expectBraces && !ignoreFirst {
    p.expectPunctuation("{")
  }

  // Parse statements
  var stmts []astNode
  for !p.tokens.endOfStream {
    stmts = append(stmts, p.parseComponent(true))
    // TODO: Find out how we need to enforce line separation
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

func (p *Parser) parseIfStatement () astNode {
  p.expectPunctuation("(")
  var cond = p.parseComponent(false)
  p.expectPunctuation(")")
  var thenPart = p.parseComponent(true)

  var elsePart astNode = astNodeEmptyStatement{}
  if p.isNextKeyword("else") {
    p.tokens.read()
    elsePart = p.parseComponent(true)
  }

  return astNodeIfStatement{
    condition: cond,
    thenPart: thenPart,
    elsePart: elsePart,
  }
}

func (p *Parser) parseWhileLoop () astNode {
  p.expectPunctuation("(")
  var cond = p.parseComponent(false)
  p.expectPunctuation(")")

  var body = p.parseComponent(true)

  return astNodeWhileLoop{
    condition: cond,
    body: body,
  }
}

func (p *Parser) parseStatement (t token) astNode {
  if t.getTokenType() == tkKeyword {
    var keyword = t.(keywordToken).value
    switch keyword {
    case "let":
      return p.parseVariableDeclaration(false, false)
    case "const":
      return p.parseVariableDeclaration(true, false)
    case "var":
      return p.parseVariableDeclaration(false, true)
    case "return":
      var arg = p.parseComponent(false)
      return astNodeReturnStatement{arg: arg}
    case "if":
      return p.parseIfStatement()
    case "while":
      return p.parseWhileLoop()
    }
  }

  if t.getTokenType() == tkPunctuation {
    if t.(punctuationToken).punctuation == "{" {
      return p.parseBlockStatement(true, true)
    }
  }

  panic("Unhandled AST node type in parseStatement!")
}

func (p *Parser) ParseAST () []astNode {
  p.ast = p.parseBlockStatement(false, false).nodes

  return p.ast
}

func New (code string) Parser {
  var codeStream = stringStream { code: code }
  var tkStream = tokenStream { code: codeStream }
  tkStream.read()
  var p = Parser { tokens: tkStream }

  return p
}
