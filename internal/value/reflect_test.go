package value

import "testing"

func TestReflectResolver_WriteColumns(t *testing.T) {
	writeColumnsTestFunc(t, NewReflectResolver)
}
