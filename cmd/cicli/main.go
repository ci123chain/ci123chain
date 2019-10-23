package main

import (
	"CI123Chain/pkg/client/cmd"
	"github.com/spf13/cobra"
)

func main()  {
	cobra.EnableCommandSorting = false
	cmd.Execute()
}
