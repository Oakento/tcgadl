package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var HOME string = os.Getenv("HOME")
var DlDir string
var Proxy string

func main() {

	var dlProj arrayFlags
	var dlAll bool

	dlCmd := flag.NewFlagSet("dl", flag.ExitOnError)
	dlCmd.Var(&dlProj, "proj", "TCGA project name")
	dlCmd.BoolVar(&dlAll, "all", false, "Download all projects")
	dlCmd.StringVar(&DlDir, "dir", path.Join(HOME, "tcgadl"), "Downloading directory")
	dlCmd.StringVar(&Proxy, "proxy", os.Getenv("HTTP_PROXY"), "HTTP_PROXY")

	if len(os.Args) < 2 {
		fmt.Println("Please specify a command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "dl":
		dlCmd.Parse(os.Args[2:])
		invalidProj := Difference(dlProj, TCGA_PROJ)
		if dlAll {
			dlProj = TCGA_PROJ
		} else if len(invalidProj) > 0 {
			fmt.Println("Unrecognized TCGA project name:", invalidProj)
			os.Exit(1)
		}
		HandleDl(dlProj)
	}

}
