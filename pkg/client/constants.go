package client

import "zheleznovux.com/modbus-console/pkg/client/tag"

// состояния сети
const (
	GOOD = "good"
	BAD  = "bad"
)

// протоколы
const (
	MODBUS_TCP = "modbusTCP"
)

// функции modbus
const (
	FUNCTION_1 = 0x1
	FUNCTION_2 = 0x2
	FUNCTION_3 = 0x3
	FUNCTION_4 = 0x4
	FUNCTION_5 = 0x5
	FUNCTION_6 = 0x6
)

var FUNCTION__TAG_TYPE = map[string][]int{
	tag.COIL_TYPE:  {FUNCTION_1, FUNCTION_2},
	tag.WORD_TYPE:  {FUNCTION_3, FUNCTION_4},
	tag.DWORD_TYPE: {FUNCTION_3, FUNCTION_4},
}