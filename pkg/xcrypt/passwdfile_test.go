package xcrypt

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/terefang/gommons/pkg/xstrings"
)

func Test_ReadFromHtpasswd(t *testing.T) {
	t1 := time.Now().UnixMicro()
	_up, _ur, err := ReadFromHtpasswd("test/htpasswd", false)
	if err != nil {
		t.Fatal(err)
	}
	_ru := make(map[string][]string)
	for k, v := range _ur {
		roles := strings.Split(strings.ToUpper(v), ",")
		for _, role := range roles {
			_ru[role] = append(_ru[role], k)
		}
	}
	t2 := time.Now().UnixMicro()
	td := t2 - t1
	fmt.Printf("Read %d entries in %d us\n", len(_up), td)
	fmt.Println(xstrings.Stringify(_up))
	fmt.Println(xstrings.Stringify(_ur))
	fmt.Println(xstrings.Stringify(_ru))
}
