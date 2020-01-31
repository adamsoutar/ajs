package parser

import "strings"

type tokenType int
const (
  tkString tokenType = iota
  tkNumber
  tkIdentifier
  tkKeyword
  tkOperator
  tkLineTerminator
  tkPunctuation
)

type token interface {
  getTokenType() tokenType
}

type lineTerminatorToken struct {}
func (l lineTerminatorToken) getTokenType () tokenType {
  return tkLineTerminator
}

type punctuationToken struct {
  punctuation string
}
func (p punctuationToken) getTokenType () tokenType {
  return tkPunctuation
}

type stringToken struct {
  value string
}
func (s stringToken) getTokenType () tokenType {
  return tkString
}

type numberToken struct {
  value float64
}
func (n numberToken) getTokenType () tokenType {
  return tkNumber
}

type identifierToken struct {
  value string
}
func (s identifierToken) getTokenType () tokenType {
  return tkIdentifier
}

type keywordToken struct {
  value string
}
func (k keywordToken) getTokenType () tokenType {
  return tkKeyword
}

type operatorToken struct {
  operator string
}
func (o operatorToken) getTokenType () tokenType {
  return tkOperator
}

func inStringArray(str string, arr []string) bool {
  for _, el := range arr {
    if el == str {
      return true
    }
  }
  return false
}
// These two functions are identical.
// Can't merge them because go doesn't have generics...
// or an stdLib way of telling if something's in an array it seems.
func inAstTypeArray (as astType, arr []astType) bool {
  for _, el := range arr {
    if el == as {
      return true
    }
  }
  return false
}

func isOperatorChar (str string) bool {
  opChars := strings.Split("=!&|+-/*%><^.", "")
  return inStringArray(str, opChars)
}
func isOperator (str string) bool {
  ops := getOperators()
  return inStringArray(str, ops)
}

func isKeyword (str string) bool {
  res := []string { "let", "const", "var", "function" }
  return inStringArray(str, res)
}

func isPunctuation (str string) bool {
  punc := strings.Split("(,){}", "")
  return inStringArray(str, punc)
}

// TODO: Automatic Semicolon Insertion
var lineTerminators = ";\n"
func isLineTerminator (str string) bool {
  terms := strings.Split(lineTerminators, "")
  return inStringArray(str, terms)
}

func isIdentifierChar (str string) bool {
  idenChars := strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz$_0123456789", "")
  return inStringArray(str, idenChars)
}

func isNumberChar (str string) bool {
  numChars := strings.Split("0123456789.", "")
  return inStringArray(str, numChars)
}

// NOTE: lineTerminators are also whitespace
var whitespace = " \t\r" + lineTerminators
func isWhitespace (str string) bool {
  whitespaceChars := strings.Split(whitespace, "")
  return inStringArray(str, whitespaceChars)
}
