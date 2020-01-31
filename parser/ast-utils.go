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
	astVariableDeclaration
	astReturnStatement
	astObject
)

type astNode interface {
	getNodeType() astType
}

type astNodeObject struct {
	valueMap map[astNode]astNode
}
func (o astNodeObject) getNodeType () astType {
	return astObject
}

type astNodeReturnStatement struct {
	arg astNode
}
func (r astNodeReturnStatement) getNodeType () astType {
	return astReturnStatement
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

type astNodeVariableDeclaration struct {
	varName string
	// These two flags should NOT be set at once
	isConstant bool
	isHoisted bool
	value astNode//? (nullable)
}
func (vD astNodeVariableDeclaration) getNodeType () astType {
	return astVariableDeclaration
}

type astNodeAssignment struct {
	variable astNode
	value    astNode
	operator string
}
func (a astNodeAssignment) getNodeType () astType {
	return astAssignment
}