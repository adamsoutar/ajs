package parser

type stringStream struct {
  code string
  index int
  row int
  col int
  endOfStream bool
}

func (s *stringStream) peek () string {
  return string(s.code[s.index])
}

func (s *stringStream) read () string {
  var p = s.peek()
  s.index++
  if s.index >= len(s.code) {
    s.endOfStream = true
  }
  return p
}
