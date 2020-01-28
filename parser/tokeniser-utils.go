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

func inArray (str string, arr []string) bool {
  for _, el := range arr {
    if el == str {
      return true
    }
  }
  return false
}

func isOperatorChar (str string) bool {
  opChars := strings.Split("=!&|+-/*%><^.", "")
  return inArray(str, opChars)
}
func isOperator (str string) bool {
  ops := getOperators()
  return inArray(str, ops)
}

func isKeyword (str string) bool {
  res := []string { "let", "const" }
  return inArray(str, res)
}

func isPunctuation (str string) bool {
  punc := strings.Split("(,)", "")
  return inArray(str, punc)
}

// TODO: Automatic Semicolon Insertion
func isLineTerminator (str string) bool {
  terms := strings.Split(";\n", "")
  return inArray(str, terms)
}

func isIdentifierChar (str string) bool {
  idenChars := strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz$_0123456789", "")
  return inArray(str, idenChars)
}

func isNumberChar (str string) bool {
  numChars := strings.Split("0123456789.", "")
  return inArray(str, numChars)
}

func isWhitespace (str string) bool {
  whitespaceChars := strings.Split(" \t\r", "")
  return inArray(str, whitespaceChars)
}
