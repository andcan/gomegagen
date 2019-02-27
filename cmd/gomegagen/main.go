package main

import (
	"github.com/andcan/gomegagen/generators"
	"github.com/spf13/pflag"
	"k8s.io/gengo/args"
	"k8s.io/klog"
	"path/filepath"
)

func main() {
	arguments := args.Default()

	arguments.GoHeaderFilePath = filepath.Join(args.DefaultSourceTree(), "github.com/andcan/gomegagen/boilerplate/boilerplate.go.txt")

	customArgs := &generators.CustomArgs{}
	pflag.CommandLine.StringSliceVar(&customArgs.WhitelistStructs, "whitelist-struct", nil, "")
	pflag.CommandLine.StringSliceVar(&customArgs.BlacklistStructs, "blacklist-struct", nil, "")
	//pflag.CommandLine.StringSliceVar(&customArgs.BlacklistFields, "blacklist-field", nil, "") // not implemented yet

	arguments.CustomArgs = customArgs

	err := arguments.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	)
	if nil != err {
		klog.Fatalf("Error: %v", err)
	}
	klog.Info("Completed successfully.")
}
