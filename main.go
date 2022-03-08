package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var HOME string = os.Getenv("HOME")
var DlDir string
var DlDecompress bool
var DlSkip bool
var Proxy string

func main() {

	var dlProj arrayFlags
	var dlAll bool

	dlCmd := flag.NewFlagSet("dl", flag.ExitOnError)
	dlCmd.Var(&dlProj, "proj", "TCGA project name. Use tcgadl abbr to check all available TCGA projects.")
	dlCmd.Var(&dlProj, "p", "TCGA project name. Use tcgadl abbr to check all available TCGA projects.")
	dlCmd.BoolVar(&dlAll, "all", false, "Download all available TCGA projects.")
	dlCmd.BoolVar(&dlAll, "o", false, "Download all available TCGA projects.")
	dlCmd.StringVar(&DlDir, "dir", path.Join(HOME, "tcgadl"), "Downloading directory. Default: $HOME/tcgadl")
	dlCmd.StringVar(&DlDir, "d", path.Join(HOME, "tcgadl"), "Downloading directory. Default: $HOME/tcgadl")
	dlCmd.StringVar(&Proxy, "proxy", os.Getenv("HTTP_PROXY"), "HTTP_PROXY")
	dlCmd.StringVar(&Proxy, "x", os.Getenv("HTTP_PROXY"), "HTTP_PROXY")
	dlCmd.BoolVar(&DlDecompress, "decompress", false, "Decompress gzipped files.")
	dlCmd.BoolVar(&DlDecompress, "u", false, "Decompress gzipped files.")
	dlCmd.BoolVar(&DlSkip, "skip", false, "Skip existing files. This option will automatically set --decompress.")
	dlCmd.BoolVar(&DlSkip, "k", false, "Skip existing files. This option will automatically set --decompress.")

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
		if DlSkip && !DlDecompress {
			DlDecompress = true
		}

		HandleDl(dlProj)

	}

}
