package _const

import "math"

const (
	DEF_SEQ_INCREMENT_RATE = 1.2
	ENOUGH_LARGE_SEQ       = 7
	DESTINCT_INCREASE_RATE = 1.2
	DEF_SEQ                = float64(math.MaxInt32 >> 2)
	WORD_SEGMENT_LENGTH    = 7
	MAX_SEQ                = 2 << 4
	DEF_EXECUTE_TIMEOUT    = 5 * 1000
	DEF_MAX_CHAIN_DEPTH    = 2 << 5
)
