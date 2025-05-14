package cmd

import (
	"fmt"

	"github.com/CLBRITTON2/go-git-good/common"
)

func Init(flags []string) {
	path := "."
	if len(flags) >= 1 {
		path = flags[0]
	}
	_, err := common.CreateRepository(path)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}
