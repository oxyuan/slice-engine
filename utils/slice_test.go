package utils

import "testing"

func TestDemoSlice(t *testing.T) {
	//slice := []string{"1", "2", "3"}
	//slice := make([]string, 0)
	//for i, s := range slice {
	//	t.Logf("slice[%d]=%s", i, s)
	//}

	slice := new([]string)
	appendPoint(slice, "4")

	for i, s := range *slice {
		t.Logf("slice[%d]=%s", i, s)
	}
}

func appendPoint(slice *[]string, value string) {
	*slice = append(*slice, value)
}

func appendTo(slice []string, value string) {
	slice = append(slice, value)
}
