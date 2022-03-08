package main

import (
	"fmt"
)

const DATA_EP string = "https://api.gdc.cancer.gov/data"
const FILES_EP string = "https://api.gdc.cancer.gov/files"

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

// Rewrite of flag type Value to allow for multiple flag value
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
