package parser

type astType int
const (
	astString astType = iota
	astNumber
	astIdentifier
	astBinary
	astAssignment
	astFunctionCall
	astBlock
	astPropertyAccess
)

type astNode interface {
	getNodeType() astType
}

type astNodePropertyAccess struct {
	object astNode
	property astNode
}
func (p astNodePropertyAccess) getNodeType () astType {
	return astPropertyAccess
}

type astNodeBlock struct {
	nodes []astNode
}
func (b astNodeBlock) getNodeType () astType {
	return astBlock
}

type astNodeString struct {
	value string
}
func (s astNodeString) getNodeType () astType {
	return astString
}

type astNodeFunctionCall struct {
	funcName astNode
	args []astNode
}
func (f astNodeFunctionCall) getNodeType () astType {
	return astFunctionCall
}

type astNodeNumber struct {
	value float64
}
func (n astNodeNumber) getNodeType () astType {
	return astNumber
}

type astNodeIdentifier struct {
	name string
}
func (i astNodeIdentifier) getNodeType () astType {
	return astIdentifier
}

type astNodeBinary struct {
	operator string
	left astNode
	right astNode
}
func (b astNodeBinary) getNodeType () astType {
	return astBinary
}

type astNodeAssignment struct {
	vars []string
	value astNode
}
func (a astNodeAssignment) getNodeType () astType {
	return astAssignment
}