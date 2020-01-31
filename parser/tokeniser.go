package parser

import (
  "strconv"
)

type readDecider func(string) bool
type tokenStream struct {
  code stringStream
  current token
  // Allows you to read() the second to last token
  consumedEnd bool
  endOfStream bool
}

func (t *tokenStream) peek () token {
  return t.current
}

func (t *tokenStream) read () token {
  var c = t.peek()
  t.current = t.readNext()

  t.readWhile(isWhitespace)

  if t.code.endOfStream {
    t.endOfStream = true
  }
  return c
}

func (t *tokenStream) readWhile (shouldRead readDecider) string {
  var final = ""
  for !t.code.endOfStream && shouldRead(t.code.peek()) {
    final += t.code.read()
  }
  return final
}

func (t *tokenStream) readString () token {
  var outStr = ""
  t.code.read()
  for !t.code.endOfStream {
    var next = t.code.read()
    if next == "\"" {
      break
    }
    outStr += next
  }
  return stringToken { value: outStr }
}

func (t *tokenStream) readIdentifierOrKeyword () token {
  var idenStr = t.readWhile(isIdentifierChar)
  if isKeyword(idenStr) {
    return keywordToken{value:idenStr}
  }
  return identifierToken { value: idenStr }
}

func (t *tokenStream) readNumber () token {
  var numString = t.readWhile(isNumberChar)
  var num, err = strconv.ParseFloat(numString, 64)

  if err != nil {
    panic(err)
  }

  return numberToken{value: num}
}

func (t *tokenStream) readOperator () token {
  var opStr = t.readWhile(isOperatorChar)
  if !isOperator(opStr) {
    panic("Unknown operator seen in source: " + opStr)
  }
  return operatorToken{operator:opStr}
}

func (t *tokenStream) readNext () token {
  t.readWhile(isWhitespace)
  if t.endOfStream && !t.consumedEnd {
    t.consumedEnd = true
    return lineTerminatorToken{}
  }

  if t.endOfStream {
    panic("Attempted a token read past the end of the code!")
  }

  var ch = t.code.peek()

  if ch == "\"" {
    return t.readString()
  }

  if t.current != nil &&
      t.current.getTokenType() != tkIdentifier &&
      t.current.getTokenType() != tkPunctuation {
    // This reads numbers like .5 before it assumes them to be
    // property access, which can only occur if the previous token
    // was an identifier or the ) punctuation
    if isNumberChar(ch) {
      return t.readNumber()
    }
  }

  if isOperatorChar(ch) {
    return t.readOperator()
  }
  // Numbers and identifiers overlap, but numbers come first
  if isNumberChar(ch) {
    return t.readNumber()
  }
  if isIdentifierChar(ch) {
    return t.readIdentifierOrKeyword()
  }
  if isPunctuation(ch) {
    t.code.read()
    return punctuationToken{punctuation:ch}
  }
  if isLineTerminator(ch) {
    t.code.read()
    return lineTerminatorToken{}
  }

  panic("Parser saw an unexpected character \"" + ch + "\"!")
}
