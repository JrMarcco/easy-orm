package val

import "testing"

func TestRefVal_WriteCols(t *testing.T) {
	testValWriteCols(t, NewRefValWriter)
}
