package main

import (
	"github.com/ci123chain/ci123chain/pkg/client/cmd"
	"github.com/spf13/cobra"
)

func main()  {
	cobra.EnableCommandSorting = false
	cmd.Execute()
}
