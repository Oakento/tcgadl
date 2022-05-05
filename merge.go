package main

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type FileInfo struct {
	FileId   string
	FileName string
	Md5sum   string
	TCGA     string
}

type RNACounts struct {
	Index  int
	Counts []string
	Sample string
}

func checkExist(proj string) int {
	_, err := os.Stat(path.Join(Dir))
	if os.IsNotExist(err) {
		return -1
	}
	_, err = os.Stat(path.Join(Dir, "manifest", proj+".csv"))
	info1, err1 := os.Stat(path.Join(Dir, proj))
	info2, err2 := os.Stat(path.Join(Dir, proj+".tar.gz"))

	res := 0
	if err == nil {
		res += 100
	}
	if err1 == nil && info1.IsDir() {
		res += 10
	}
	if err2 == nil && !info2.IsDir() {
		res += 1
	}
	return res
}

func readManifest(proj string) ([]FileInfo, error) {
	manifestFile, err := os.Open(path.Join(Dir, "manifest", proj+".csv"))
	if err != nil {
		return nil, fmt.Errorf("Error opening manifest file")
	}
	defer manifestFile.Close()

	var filesInfo []FileInfo

	content, _ := ioutil.ReadAll(manifestFile)
	for _, line := range strings.Split(string(content), "\n")[1:] {
		if line == "" {
			continue
		}
		lineSplit := strings.Split(line, ",")
		if len(lineSplit) != 7 {
			return nil, fmt.Errorf("Invalid manifest file.")
		}
		filesInfo = append(filesInfo, FileInfo{
			lineSplit[0], lineSplit[1], lineSplit[2], lineSplit[3],
		})
	}

	return filesInfo, nil
}

func checkFileValid(fileInfo []FileInfo, proj string) (bool, []string) {
	dspCh := make(chan string)
	for _, info := range fileInfo {
		go func(info FileInfo) {
			_, err := os.Stat(path.Join(Dir, proj, info.FileId, info.FileName))
			if err != nil {
				dspCh <- info.FileId + "/" + info.FileName
			} else {
				md5sum := Md5sum(path.Join(Dir, proj, info.FileId, info.FileName))
				if md5sum != info.Md5sum {
					dspCh <- info.FileId + "/" + info.FileName
				} else {
					dspCh <- ""
				}
			}
		}(info)
	}

	res := true
	d := []string{}
	for i := 0; i < len(fileInfo); i++ {
		dsp := <-dspCh
		if dsp != "" {
			res = false
			d = append(d, dsp)
		}
	}
	return res, d
}

func merge(filesInfo []FileInfo, proj string) {
	// id, file_name, TCGA_Barcode, md5sum, vital_status, days_to_death, days_to_last_follow_up
	file, err := os.Open(path.Join(Dir, proj, filesInfo[0].FileId, filesInfo[0].FileName))
	if err != nil {
		fmt.Println("Error opening file:" + filesInfo[0].FileId + "/" + filesInfo[0].FileName)
		return
	}
	defer file.Close()
	// reader, _ := gzip.NewReader(file)
	// defer reader.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		// fmt.Println("Broken gzip file:" + filesInfo[0].FileId + "/" + filesInfo[0].FileName)
		return
	}
	lines := strings.Split(string(content), "\n")

	rnaCounts := make([][]string, len(lines)-5)
	for i := range rnaCounts {
		rnaCounts[i] = make([]string, len(filesInfo)+1)
	}
	rnaCounts[0][0] = ""
	rnaCounts[0][1] = filesInfo[0].TCGA
	for i, line := range lines[:len(lines)-6] {
		lineSplit := strings.Split(line, "\t")
		rnaCounts[i+1][0] = lineSplit[0]
		rnaCounts[i+1][1] = lineSplit[1]
	}
	fileCh := make(chan *RNACounts)
	for i, info := range filesInfo[1:] {
		go func(info FileInfo, i int) {
			file, err := os.Open(path.Join(Dir, proj, info.FileId, info.FileName))
			if err != nil {
				fmt.Println("Error opening file:" + info.FileId + "/" + info.FileName)
				fileCh <- &RNACounts{-1, nil, ""}
				return
			}
			defer file.Close()
			reader, _ := gzip.NewReader(file)
			defer reader.Close()
			content, err := ioutil.ReadAll(reader)
			if err != nil {
				fmt.Println("Broken gzip file:" + info.FileId + "/" + info.FileName)
				fileCh <- &RNACounts{-1, nil, ""}
				return
			}
			lines := strings.Split(string(content), "\n")
			fileCh <- &RNACounts{i, lines[:len(lines)-6], info.TCGA}
		}(info, i)
	}
	for i := 1; i < len(filesInfo); i++ {
		var cnts *RNACounts
		cnts = <-fileCh
		if cnts.Index == -1 {
			fmt.Printf("Sample %s reading error, skipped.\n", cnts.Sample)
			continue
		}
		rnaCounts[0][cnts.Index+2] = cnts.Sample

		for j, line := range cnts.Counts {
			lineSplit := strings.Split(line, "\t")
			rnaCounts[j+1][cnts.Index+2] = lineSplit[1]
		}
		// rnaCounts[cnts.Index] = append([]string{cnts.Sample}, cnts.Counts...)
	}
	close(fileCh)

	_ = os.MkdirAll(path.Join(Dir, "merge"), os.ModePerm)
	mergeOut, _ := os.Create(path.Join(Dir, "merge", proj+".csv"))
	defer mergeOut.Close()

	mergeWriter := csv.NewWriter(mergeOut)
	mergeWriter.Comma = '\t'
	for i := range rnaCounts {
		_ = mergeWriter.Write(rnaCounts[i])
	}
	mergeWriter.Flush()

}

func decompress(proj string) {
	gzipFile, err := os.Open(path.Join(Dir, proj+".tar.gz"))
	if err != nil {
		fmt.Println("Error opening gzip file.")
		os.Exit(1)
	}
	defer gzipFile.Close()

	content, _ := ioutil.ReadAll(gzipFile)
	WriteDecompressed(bytes.NewReader(content), proj, true)
}

func HandleMerge(projects []string) {

	for _, proj := range projects {
		fileExist := checkExist(proj)
		whereData := fileExist % 100

		// -1, 0, 100
		if whereData <= 0 {
			fmt.Printf("Project %s file not found.\n"+
				"Use `tcga dl --proj %s --dir %s --decompress` to donwload it.\n",
				proj, proj, Dir)
			return
		}

		if fileExist < 100 {
			fmt.Printf("Missing manifest/%s.csv\n"+
				" Use `tcga dl --proj %s --dir %s --decompress` to re-download the project ",
				proj, proj, Dir)
			return
		}
		// 10, 01, 11
		if whereData == 1 {
			decompress(proj)
		}
		fileInfo, err := readManifest(proj)
		if err != nil {
			fmt.Println(err)
			return
		}
		valid, dsp := checkFileValid(fileInfo, proj)
		if !valid {
			fmt.Println("Invalid file:", dsp)
			return
		}
		merge(fileInfo, proj)

	}

}
