package xcrypt

import (
	"fmt"
	"testing"
)

func Test_GenerateWordPass(t *testing.T) {
	for i := 0; i < 16; i++ {
		fmt.Println(GenerateWordPass(16))
	}
}
