package val

import "testing"

func TestUnsafeVal_WriteCols(t *testing.T) {
	testValWriteCols(t, NewUnsafeValWriter)
}
