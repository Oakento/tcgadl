# tcgadl

Fast downloading and merging TCGA files

## Usage


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
# merge project data to a table
tcgadl merge --proj TCGA-LUSC --dir ./tcga/
```

#### subcommand `dl`
`--proj`: TCGA project to download

`--all`: download all TCGA projects

`--dir`: specify downloading directory(default `$HOME/tcgadl`)

`--decompress`: decompress downloaded files(default `false`)

`--skip`: skip downloaded files(default `false`, will automatically enable `--decompress`)

`--proxy`: use proxy. Having priority over `HTTP_PROXY`

#### subcommand `merge`
`--proj`: TCGA project to download

`--all`: download all TCGA projects

`--dir`: specify downloading directory(default `$HOME/tcgadl`)

#### subcommand `show`
Display available TCGA projects.

#### subcommand `help`
Display help information.

#### effective environment variables
`$HTTP_PROXY` `$HOME`


## Compilation
#### Requirements
`go 1.17`

#### Linux / macOS
```shell
go build -o tcgadl -trimpath -ldflags "-s -w -buildid="
```
#### Windows
```shell
go build -o tcgadl.exe -trimpath -ldflags "-s -w -buildid="
```