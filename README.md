# tcgadl

Download TCGA files(faster than gdc-client)

### Usage


```shell
# download TCGA-LUSC and TCGA-LUAD to ./tcga/ 
tcgadl dl --proj TCGA-LUSC --proj TCGA-LUAD --dir ./tcga/
# with proxy
tcgadl dl --proj TCGA-LUSC --proxy http://127.0.0.1:7890
# download all TCGA projects
tcgadl dl --all
# download and decompress all TCGA projects
tcgadl dl --all --decompress
# download all and skip existing files
tcgadl dl --all --skip
```

#### subcommand `dl`
`--proj`: TCGA project to download

`--all`: download all TCGA projects

`--dir`: specify downloading directory(default `$HOME/tcgadl`)

`--decompress`: decompress downloaded files(default `false`)

`--skip`: skip downloaded files(default `false`)

`--proxy`: use proxy. Having priority over `HTTP_PROXY`

#### subcommand `show`
Display available TCGA projects.


#### effective environment variables
`$HTTP_PROXY` `$HOME`
