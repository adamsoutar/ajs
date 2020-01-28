package parser

import "fmt"

type Parser struct {
  tokens tokenStream
  ast []astNode
}

func (p *Parser) PrintAll() {
  for true {
    var tk = p.tokens.read()

    var tp = tk.getTokenType()
    switch tp {
      case tkString:
        fmt.Println("STRING: " + tk.(stringToken).value)
      case tkIdentifier:
        fmt.Println("IDENTIFIER: " + tk.(identifierToken).value)
      case tkNumber:
        var numStr = fmt.Sprintf("%f", tk.(numberToken).value)
        fmt.Println("NUMBER: " + numStr)
      case tkOperator:
        fmt.Println("OPERATOR:" + tk.(operatorToken).operator)
      default:
        fmt.Println("UNKNOWN TOKEN - Probably invalid token")
    }

    if p.tokens.endOfStream {
      break
    }
  }
}

func (p *Parser) isNextOperator () bool {
  return p.tokens.peek().getTokenType() == tkOperator
}

func (p *Parser) mightBeBinary (left astNode, myPrecedence int) astNode {
  if p.isNextOperator() {
    var op = p.tokens.peek().(operatorToken)

    var theirPrecedence =
  }
  return left
}

func (p *Parser) parseAssigment (isConst bool) astNode {
  // Read in the vars
  var vars []string
  for true {
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

  p.expectToken(tkLineTerminator)

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
  // TODO: Wrap this in mightBeBinary etc.
  return p.parseAtom(acceptStatements)
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
  }

  // Statement
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

func (p *Parser) ParseAST () {
  var ast []astNode
  for !p.tokens.endOfStream {
    // TODO: See if we actually want to parse an expression
    //       eg. a line that just has a function call
    ast = append(ast, p.parseComponent(true))
  }
  p.ast = ast
}

func New (code string) Parser {
  var codeStream = stringStream { code: code }
  var tkStream = tokenStream { code: codeStream }
  tkStream.read()
  var p = Parser { tokens: tkStream }

  return p
}
