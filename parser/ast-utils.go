package parser

type astType int
const (
	astString astType = iota
	astNumber
	astIdentifier
	astBinary
	astAssignment
	astFunctionCall
	astFunctionDefinition
	astBlock
	astPropertyAccess
	astEmptyStatement
)

type astNode interface {
	getNodeType() astType
}

type astNodeEmptyStatement struct {}
func (e astNodeEmptyStatement) getNodeType () astType {
	return astEmptyStatement
}

type astNodeFunctionDefinition struct {
	name string
	params []string
	body astNodeBlock
}
func (fD astNodeFunctionDefinition) getNodeType () astType {
	return astFunctionDefinition
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
	varNm string
	value astNode
}
func (a astNodeAssignment) getNodeType () astType {
	return astAssignment
}