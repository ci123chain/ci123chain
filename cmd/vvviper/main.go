package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func main() {
	viper.SetEnvPrefix("CI")
	viper.AutomaticEnv()

	_ = os.Setenv("CI_KE-Y","env")  //设置环境变量
	//_ = viper.BindEnv("KE-Y") //绑定环境变量

	fmt.Println(viper.Get("ke-y")) //获取环境变量
	fmt.Println()

	rootCmd := &cobra.Command{
		Use: 	"test",
		Short:  "",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	rootCmd.Flags().String("ke-y", "cmdFlag", "Run abci app with different log level")
	_ = rootCmd.Execute()

	_ = viper.BindPFlags(rootCmd.Flags())
	fmt.Println(viper.Get("KE-Y"))
	fmt.Println()

	//pflag.String("KEY","flag","")
	//_ = viper.BindPFlag("KEY", pflag.Lookup("KEY"))
	//fmt.Println(viper.Get("KEY"))
	//fmt.Println()
	//_ = viper.BindPFlags(pflag.CommandLine)

	//cmd > env > flagDefault
}

func fromEnv(env string) string {
	_ = viper.BindEnv(env) //绑定环境变量

	return viper.Get(env).(string) //获取环境变量
}

