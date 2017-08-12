package models

// LBlock is the RawLBlock after AumActionSet has been bundled
type LBlock struct {
	AlwaysExec string          `json:"always"`
	Statements *[][]LStatement `json:"statements"`
}

// LStatement is the RawLStatement after AumActionSet has been bundled
type LStatement struct {
	Operators *OrGroup `json:"conditions"`
	Exec      string
}

// RawLBlock contains every execution block
// AlwaysExec are actions that have no conditions
// Statements are an array of arrays containing objects that might look like this
// Statements: [ [ If: [], ElIf: [], Else: [] ], [ If: [], ElIf: [], Else[] ] ]
// Each []RawLStatement is executed in order as they mutate the runtime state
// Each RawLStatement within []RawLStatement is an implicit "if"/"elif"/"else" block,
// where the conditions of each one are tested in order until one results in true
type RawLBlock struct {
	AlwaysExec AumActionSet       `json:"always"`
	Statements *[][]RawLStatement `json:"statements"`
}

// RawLStatement contains an OrGroup of AndGroups
// The OrGroup must yield true for the Exec AumActionSet to execute
// Exec mutates the runtimes state of an AUM instance
type RawLStatement struct {
	Operators *OrGroup     `json:"conditions"`
	Exec      AumActionSet `json:"then"`
}

// VarValMap contains a mapping of variables
// and their respective comparison values
type VarValMap map[int]interface{}

// AndGroup contains a map of operators with VarValMaps
// All variables when compared to their values with respect to OperatorStr
// must yield true in order for a RawLStatement to execute its Exec AumActionSet
type AndGroup map[OperatorStr]VarValMap

// OrGroup contains a list of AndGroups
// At least one AndGroup must yield true in order for
// a RawLStatement to execute its Exec AumActionSet
type OrGroup []AndGroup

// OperatorStr allows a workbench user to create conditional logic
type OperatorStr string

const (
	// OpStrEQ =
	OpStrEQ OperatorStr = "eq"
	// OpStrLT <
	OpStrLT OperatorStr = "lt"
	// OpStrGT >
	OpStrGT OperatorStr = "gt"
	// OpStrLE <=
	OpStrLE OperatorStr = "le"
	// OpStrGE >=
	OpStrGE OperatorStr = "ge"
	// OpStrNE !=
	OpStrNE OperatorStr = "ne"
)

// OperatorInt is a compiled OperatorStr
// Compiled by Lakshmi and for use by the runtime Brahman
type OperatorInt int8

const (
	// OpIntEQ =
	OpIntEQ OperatorInt = 1 << iota
	// OpIntLT <
	OpIntLT
	// OpIntGT >
	OpIntGT
	// OpIntLE <=
	OpIntLE
	// OpIntGE >=
	OpIntGE
	// OpIntNE !=
	OpIntNE
)

// GenerateOperatorStrIntMap is a helper function for Lakshmi's compilation process
func GenerateOperatorStrIntMap() map[OperatorStr]OperatorInt {
	return map[OperatorStr]OperatorInt{
		OpStrEQ: OpIntEQ,
		OpStrLT: OpIntLT,
		OpStrGT: OpIntGT,
		OpStrLE: OpIntLE,
		OpStrGE: OpIntGE,
		OpStrNE: OpIntNE,
	}
}
