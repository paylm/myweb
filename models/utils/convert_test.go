package utils

import "testing"

func Test_Convert(t *testing.T) {
	done := make(chan error)
	go convert("../../upload/pp.pdf", "400", done)
	err := <-done
	if err != nil {
		t.Errorf("test fail with err:%v\n", err)
	}
}
