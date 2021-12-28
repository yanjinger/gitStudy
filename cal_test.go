package gitStudy

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	if 2 == Add(1, 1) {
		fmt.Println("ok")
	} else {
		fmt.Println("failed")
	}
}
