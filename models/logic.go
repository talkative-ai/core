package models

import (
	"database/sql/driver"
	"encoding/json"
)

// LBlock is the RawLBlock after AumActionSet has been bundled
type LBlock struct {
	// AlwaysExec are actions that have no conditions
	AlwaysExec string

	// Statements is an array of arrays containing objects that are
	// implicit "if/elif/else" blocks. e.g.
	// Statements:
	//	[
	//		[ If: []LBlock, ElIf: []LBlock, Else: []LBlock ]LStatement,
	//		[ If: []LBlock, ElIf: []LBlock, Else: []LBlock ]LStatement,
	//	]
	// Consider an []LStatement of length n
	// LStatement[0] is the "if" statement
	// LStatement[1 : n-1] are "else if" statements
	// LStatement[n-1] is the "else" statement
	// If n == 2, then LStatement[n-1] might be "elif" or "else"
	// depending on whether Operators == nil
	// If Operators == nil then we expect Exec to be executed right away
	//
	// Each []LStatement is executed in order as they mutate the runtime state
	Statements *[][]LStatement
}

// LStatement is the RawLStatement after AumActionSet has been bundled
// The OrGroup must yield true for the Exec AumActionSet to execute
// Exec mutates the runtimes state of an AUM instance
type LStatement struct {
	Operators *OrGroup
	Exec      string
}

// RawLBlock contains every execution block
type RawLBlock struct {
	AlwaysExec AumActionSet
	Statements *RawLStatementUnified
}

type RawLStatementUnified [][]RawLStatement

func (a *RawLStatementUnified) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &a)
}

func (a *RawLStatementUnified) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

// RawLStatement contains an OrGroup of AndGroups
type RawLStatement struct {
	Operators *OrGroup
	Exec      AumActionSet
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
