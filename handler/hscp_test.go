package handler

import (
	"testing"
	"fmt"
)

func TestExecShell(t *testing.T)  {
	result,err:=execShell("lsw")
	fmt.Println(result)
	fmt.Println("----")
	fmt.Println(err)
	
}