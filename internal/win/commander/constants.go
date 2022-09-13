package commander

// опреаторы сравнения
const (
	EQUAL      = "=="
	NOT_EQUAL  = "!="
	MORE_EQUAL = ">="
	LESS_EQUAL = "<="
	MORE       = ">"
	LESS       = "<"
	BIT        = "bit"
	NOT_BIT    = "!bit"
)

// логические операторы
const (
	AND = "and"
	OR  = "or"
)

// команды для win
const (
	SHUTDOWN    = "shutdown"
	RESTART     = "restart"
	RUN_PROGRAM = "run"
)

const (
	MIN_ACTION_TIMEOUT = 1
	MIN_SCAN_PERIOD    = 1
)
