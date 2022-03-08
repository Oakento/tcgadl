# tcgadl

Download TCGA files(faster than gdc-client)

### Usage


```shell
# download TCGA-LUSC and TCGA-LUAD to ./tcga/ 
tcgadl dl --proj TCGA-LUSC --proj TCGA-LUAD --dir ./tcga/
# use proxy
tcgadl dl --proj TCGA-LUSC --proxy http://127.0.0.1:7890
# download all TCGA projects
tcgadl dl --all
```



#### subcommand `dl`
`--proj`: TCGA project to download

`--all`: download all TCGA projects(shown in the table below)

`--dir`: specify downloading directory(default `$HOME/tcgadl`)

`--proxy`: use proxy. Having priority over `HTTP_PROXY`
#### effective environment variables
`$HTTP_PROXY` `$HOME`



| TCGA abbr. | TCGA Name | 
| --- | --- |
| [TCGA-BRCA](https://portal.gdc.cancer.gov/projects/TCGA-BRCA) | Breast Invasive Carcinoma |
| [TCGA-GBM](https://portal.gdc.cancer.gov/projects/TCGA-GBM) | Glioblastoma Multiforme |
| [TCGA-OV](https://portal.gdc.cancer.gov/projects/TCGA-OV) | Ovarian Serous Cystadenocarcinoma |
| [TCGA-LUAD](https://portal.gdc.cancer.gov/projects/TCGA-LUAD) | Lung Adenocarcinoma |
| [TCGA-UCEC](https://portal.gdc.cancer.gov/projects/TCGA-UCEC) | Uterine Corpus Endometrial Carcinoma |
| [TCGA-KIRC](https://portal.gdc.cancer.gov/projects/TCGA-KIRC) | Kidney Renal Clear Cell Carcinoma |
| [TCGA-HNSC](https://portal.gdc.cancer.gov/projects/TCGA-HNSC) | Head and Neck Squamous Cell Carcinoma |
| [TCGA-LGG](https://portal.gdc.cancer.gov/projects/TCGA-LGG) | Brain Lower Grade Glioma |
| [TCGA-THCA](https://portal.gdc.cancer.gov/projects/TCGA-THCA) | Thyroid Carcinoma |
| [TCGA-LUSC](https://portal.gdc.cancer.gov/projects/TCGA-LUSC) | Lung Squamous Cell Carcinoma |
| [TCGA-PRAD](https://portal.gdc.cancer.gov/projects/TCGA-PRAD) | Prostate Adenocarcinoma |
| [TCGA-SKCM](https://portal.gdc.cancer.gov/projects/TCGA-SKCM) | Skin Cutaneous Melanoma |
| [TCGA-COAD](https://portal.gdc.cancer.gov/projects/TCGA-COAD) | Colon Adenocarcinoma |
| [TCGA-STAD](https://portal.gdc.cancer.gov/projects/TCGA-STAD) | Stomach Adenocarcinoma |
| [TCGA-BLCA](https://portal.gdc.cancer.gov/projects/TCGA-BLCA) | Bladder Urothelial Carcinoma |
| [TCGA-LIHC](https://portal.gdc.cancer.gov/projects/TCGA-LIHC) | Liver Hepatocellular Carcinoma |
| [TCGA-CESC](https://portal.gdc.cancer.gov/projects/TCGA-CESC) | Cervical Squamous Cell Carcinoma and Endocervical Adenocarcinoma |
| [TCGA-KIRP](https://portal.gdc.cancer.gov/projects/TCGA-KIRP) | Kidney Renal Papillary Cell Carcinoma |
| [TCGA-SARC](https://portal.gdc.cancer.gov/projects/TCGA-SARC) | Sarcoma |
| [TCGA-LAML](https://portal.gdc.cancer.gov/projects/TCGA-LAML) | Acute Myeloid Leukemia |
| [TCGA-ESCA](https://portal.gdc.cancer.gov/projects/TCGA-ESCA) | Esophageal Carcinoma |
| [TCGA-PAAD](https://portal.gdc.cancer.gov/projects/TCGA-PAAD) | Pancreatic Adenocarcinoma |
| [TCGA-PCPG](https://portal.gdc.cancer.gov/projects/TCGA-PCPG) | Pheochromocytoma and Paraganglioma |
| [TCGA-READ](https://portal.gdc.cancer.gov/projects/TCGA-READ) | Rectum Adenocarcinoma |
| [TCGA-TGCT](https://portal.gdc.cancer.gov/projects/TCGA-TGCT) | Testicular Germ Cell Tumors |
| [TCGA-THYM](https://portal.gdc.cancer.gov/projects/TCGA-THYM) | Thymoma |
| [TCGA-KICH](https://portal.gdc.cancer.gov/projects/TCGA-KICH) | Kidney Chromophobe |
| [TCGA-ACC](https://portal.gdc.cancer.gov/projects/TCGA-ACC) | Adrenocortical Carcinoma |
| [TCGA-MESO](https://portal.gdc.cancer.gov/projects/TCGA-MESO) | Mesothelioma |
| [TCGA-UVM](https://portal.gdc.cancer.gov/projects/TCGA-UVM) | Uveal Melanoma |
| [TCGA-DLBC](https://portal.gdc.cancer.gov/projects/TCGA-DLBC) | Lymphoid Neoplasm Diffuse Large B-cell Lymphoma |
| [TCGA-UCS](https://portal.gdc.cancer.gov/projects/TCGA-UCS) | Uterine Carcinosarcoma |
| [TCGA-CHOL](https://portal.gdc.cancer.gov/projects/TCGA-CHOL) | Cholangiocarcinoma |
