package value

import "testing"

func TestUnsafeResolver_WriteColumns(t *testing.T) {
	writeColumnsTestFunc(t, NewUnsafeResolver)
}
