package main

import (
  "ajs/parser"
  "fmt"
  "io/ioutil"
)

func main () {
  var codeBytes, err = ioutil.ReadFile("./test.js")
  if err != nil {
    panic(err)
  }

  var p = parser.New(string(codeBytes))
  /*var AST = */p.ParseAST()
  fmt.Println("AST generated")
}
