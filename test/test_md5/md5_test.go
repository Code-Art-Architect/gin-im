package test_md5

import (
	"fmt"
	"testing"

	"github.com/code-art/gin-im/util"
)

func TestMD5(t *testing.T) {
	fmt.Println(util.MakePassword("123456", ""))
}
