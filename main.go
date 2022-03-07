package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var HOME string = os.Getenv("HOME")
var HTTP_PROXY string = os.Getenv("HTTP_PROXY")
var DATA_EP string = "https://api.gdc.cancer.gov/data"
var FILES_EP string = "https://api.gdc.cancer.gov/files"
var FIELDS = []string{
	// "file_name",
	"associated_entities.entity_submitter_id",
	// "associated_entities.case_id",
}

var TCGA_PROJ = []string{
	"TCGA-BRCA",
	"TCGA-GBM",
	"TCGA-OV",
	"TCGA-LUAD",
	"TCGA-UCEC",
	"TCGA-KIRC",
	"TCGA-HNSC",
	"TCGA-LGG",
	"TCGA-THCA",
	"TCGA-LUSC",
	"TCGA-PRAD",
	"TCGA-SKCM",
	"TCGA-COAD",
	"TCGA-STAD",
	"TCGA-BLCA",
	"TCGA-LIHC",
	"TCGA-CESC",
	"TCGA-KIRP",
	"TCGA-SARC",
	"TCGA-LAML",
	"TCGA-ESCA",
	"TCGA-PAAD",
	"TCGA-PCPG",
	"TCGA-READ",
	"TCGA-TGCT",
	"TCGA-THYM",
	"TCGA-KICH",
	"TCGA-ACC",
	"TCGA-MESO",
	"TCGA-UVM",
	"TCGA-DLBC",
	"TCGA-UCS",
	"TCGA-CHOL",
}

type Value interface {
	String() string
	Set(string) error
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return fmt.Sprint(*i)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {

	dlCmd := flag.NewFlagSet("dl", flag.ExitOnError)
	var dlProj arrayFlags
	var dlAll bool
	var dlDir string
	dlCmd.Var(&dlProj, "proj", "TCGA project name")
	dlCmd.BoolVar(&dlAll, "all", false, "Download all projects")
	dlCmd.StringVar(&dlDir, "dir", path.Join(HOME, "tcgadl"), "Download directory")

	if len(os.Args) < 2 {
		fmt.Println("Please specify a command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "dl":
		dlCmd.Parse(os.Args[2:])
		invalidProj := Difference(dlProj, TCGA_PROJ)
		if dlAll {
			DownloadMany(TCGA_PROJ, dlDir)
		} else if len(invalidProj) > 0 {
			fmt.Println("Unrecognized TCGA project name:", invalidProj)
			os.Exit(1)
		}
		DownloadMany(dlProj, dlDir)
	}

}
