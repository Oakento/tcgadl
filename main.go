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
	// dlCmd.Var(&dlProj, "p", "TCGA project name. Use tcgadl abbr to check all available TCGA projects.")
	dlCmd.BoolVar(&dlAll, "all", false, "Download all available TCGA projects.")
	// dlCmd.BoolVar(&dlAll, "o", false, "Download all available TCGA projects.")
	dlCmd.StringVar(&DlDir, "dir", path.Join(HOME, "tcgadl"), "Downloading directory. Default: $HOME/tcgadl")
	// dlCmd.StringVar(&DlDir, "d", path.Join(HOME, "tcgadl"), "Downloading directory. Default: $HOME/tcgadl")
	dlCmd.StringVar(&Proxy, "proxy", os.Getenv("HTTP_PROXY"), "HTTP_PROXY")
	// dlCmd.StringVar(&Proxy, "x", os.Getenv("HTTP_PROXY"), "HTTP_PROXY")
	dlCmd.BoolVar(&DlDecompress, "decompress", false, "Decompress gzipped files.")
	// dlCmd.BoolVar(&DlDecompress, "u", false, "Decompress gzipped files.")
	dlCmd.BoolVar(&DlSkip, "skip", false, "Skip existing files. This option will automatically set --decompress.")
	// dlCmd.BoolVar(&DlSkip, "k", false, "Skip existing files. This option will automatically set --decompress.")

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
	case "show":
		fmt.Print(`=============================
TCGA proj.	Description
=============================
TCGA-BRCA	Breast Invasive Carcinoma
TCGA-GBM	Glioblastoma Multiforme
TCGA-OV		Ovarian Serous Cystadenocarcinoma
TCGA-LUAD	Lung Adenocarcinoma
TCGA-UCEC	Uterine Corpus Endometrial Carcinoma
TCGA-KIRC	Kidney Renal Clear Cell Carcinoma
TCGA-HNSC	Head and Neck Squamous Cell Carcinoma
TCGA-LGG	Brain Lower Grade Glioma
TCGA-THCA	Thyroid Carcinoma
TCGA-LUSC	Lung Squamous Cell Carcinoma
TCGA-PRAD	Prostate Adenocarcinoma
TCGA-SKCM	Skin Cutaneous Melanoma
TCGA-COAD	Colon Adenocarcinoma
TCGA-STAD	Stomach Adenocarcinoma
TCGA-BLCA	Bladder Urothelial Carcinoma
TCGA-LIHC	Liver Hepatocellular Carcinoma
TCGA-CESC	Cervical Squamous Cell Carcinoma and Endocervical Adenocarcinoma
TCGA-KIRP	Kidney Renal Papillary Cell Carcinoma
TCGA-SARC	Sarcoma
TCGA-LAML	Acute Myeloid Leukemia
TCGA-ESCA	Esophageal Carcinoma
TCGA-PAAD	Pancreatic Adenocarcinoma
TCGA-PCPG	Pheochromocytoma and Paraganglioma
TCGA-READ	Rectum Adenocarcinoma
TCGA-TGCT	Testicular Germ Cell Tumors
TCGA-THYM	Thymoma
TCGA-KICH	Kidney Chromophobe
TCGA-ACC	Adrenocortical Carcinoma
TCGA-MESO	Mesothelioma
TCGA-UVM	Uveal Melanoma
TCGA-DLBC	Lymphoid Neoplasm Diffuse Large B-cell Lymphoma
TCGA-UCS	Uterine Carcinosarcoma
TCGA-CHOL	Cholangiocarcinoma
`)
	default:
		dlCmd.Parse(os.Args[2:])
		dlCmd.Usage()
	}

}
