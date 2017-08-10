package models

// LBlock is the RawLBlock after action bundles have been built
type LBlock struct {
	AlwaysExec *uint64         `json:"always"`
	Statements *[][]LStatement `json:"statements"`
}

// LStatement is the RawLStatement after action bundles have been built
type LStatement struct {
	Operators *OpArray `json:"conditions"`
	Exec      []int32
}

// LBlock contian
type RawLBlock struct {
	AlwaysExec *map[string]interface{} `json:"always"`
	Statements *[][]RawLStatement      `json:"statements"`
}

type RawLStatement struct {
	Operators *OpArray                 `json:"conditions"`
	Exec      []map[string]interface{} `json:"then"`
}

type VarValMap map[int]interface{}
type OpArray []map[OperatorStr]VarValMap

type OperatorStr string

const (
	OpStrEQ OperatorStr = "eq"
	OpStrLT OperatorStr = "lt"
	OpStrGT OperatorStr = "gt"
	OpStrLE OperatorStr = "le"
	OpStrGE OperatorStr = "ge"
	OpStrNE OperatorStr = "ne"
)

type OperatorInt int8

const (
	OpIntEQ OperatorInt = 1 << iota
	OpIntLT
	OpIntGT
	OpIntLE
	OpIntGE
	OpIntNE
)

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

type StatementInt int8

const (
	StatementIF StatementInt = 1 << iota
	StatementELIF
	StatementELSE
)
