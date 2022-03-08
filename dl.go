package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/schollz/progressbar/v3"
)

type Filter struct {
	Op      string    `json:"op"`
	Content []Content `json:"content"`
}

type Content struct {
	Op      string     `json:"op"`
	Content SubContent `json:"content"`
}

type SubContent struct {
	Field string   `json:"field"`
	Value []string `json:"value"`
}

type Params struct {
	Filters string `url:"filters"`
	Fields  string `url:"fields"`
	Format  string `url:"format"`
	Size    string `url:"size"`
}

type Metadata struct {
	Data struct {
		Hits []struct {
			ID       string `json:"id"`
			Md5sum   string `json:"md5sum"`
			FileName string `json:"file_name"`
			Entities []struct {
				Barcode string `json:"entity_submitter_id"`
			} `json:"associated_entities"`
		} `json:"hits"`
		Pagination struct {
			Count int `json:"count"`
			Total int `json:"total"`
			From  int `json:"from"`
			Page  int `json:"page"`
			Pages int `json:"pages"`
		} `json:"pagination"`
	} `json:"data"`
	Project string
}

type Manifest struct {
	FileID   string `json:"file_id"`
	FileName string `json:"file_name"`
	TCGA     string `json:"TCGA_barcode"`
	Md5sum   string `json:"md5sum"`
}

var Client *http.Client

func initClient() *http.Client {
	if Proxy == "" {
		return &http.Client{}
	}
	proxyURL, err := url.ParseRequestURI(Proxy)
	if err != nil {
		fmt.Printf(`Proxy URL "%s" is invalid. Use direct connection instead.\n`, Proxy)
		return &http.Client{}
	}
	fmt.Println("Connecting to https://api.gdc.cancer.gov via", Proxy)
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), os.ModePerm)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

func writeGzip(source io.Reader, proj, p string) {
	err := os.MkdirAll(DlDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error create directories. Check your permissions.")
		os.Exit(1)
	}
	target, _ := os.Create(p + ".tmp")
	defer target.Close()

	fmt.Print("\n")
	bar := progressbar.DefaultBytes(-1, "Downloading "+p)

	_, err = io.Copy(io.MultiWriter(target, bar), source)
	if err != nil {
		target.Close()
		os.Remove(p + ".tmp")
		fmt.Println("Error writing file to directory. Check your permissions.")
		os.Exit(1)
	}
	target.Close()
	_ = os.Rename(p+".tmp", p)
}

func writeDecompressed(source io.Reader, proj string) {
	err := os.MkdirAll(path.Join(DlDir, proj), os.ModePerm)
	if err != nil {
		fmt.Println("Error create directories. Check your permissions.")
		os.Exit(1)
	}
	fmt.Print("\n")
	bar := progressbar.DefaultBytes(-1, "Downloading "+proj)

	var buf bytes.Buffer
	_, err = io.Copy(io.MultiWriter(&buf, bar), source)
	if err != nil {
		fmt.Println("Error receiving files. Check your network connection and proxy settings.")
		os.Exit(1)
	}
	fmt.Print("Decompressing...")
	reader, err := gzip.NewReader(&buf)
	if err != nil {
		fmt.Println("Broken gzip file.")
		return
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading tar file.")
			return
		}
		fileName := path.Join(DlDir, proj, header.Name)
		target, _ := createFile(fileName)
		defer target.Close()
		_, err = io.Copy(target, tarReader)
		if err != nil {
			target.Close()
			os.Remove(path.Join(DlDir, proj, header.Name))
			fmt.Println("Error writing to file.")
			return
		}
	}
	fmt.Print("Done.")
}

func writeManifest(proj string, manifest []Manifest) {
	_ = os.MkdirAll(path.Join(DlDir, "manifest"), os.ModePerm)
	manifestOut, _ := os.Create(path.Join(DlDir, "manifest", proj+".csv"))
	defer manifestOut.Close()

	manifestWriter := csv.NewWriter(manifestOut)
	_ = manifestWriter.Write([]string{"file_id", "file_name", "TCGA_manifest"})
	for _, record := range manifest {
		_ = manifestWriter.Write([]string{
			record.FileID,
			record.FileName,
			record.TCGA,
		})
	}
	manifestWriter.Flush()
}

func dl(fileIds []string, manifest []Manifest, proj string) {
	payload, _ := json.Marshal(url.Values{"ids": fileIds})
	req, _ := http.NewRequest("POST", DATA_EP, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Range", "bytes=0-")
	// req.Header.Set("Host", "api.gdc.cancer.gov")
	resp, err := Client.Do(req)
	if err != nil {
		fmt.Println("Error sending request. Check your network connection and proxy settings.")
		os.Exit(1)
	}
	// resp, _ := http.PostForm(DATA_EP, url.Values{"ids": fileIds})
	// fileSize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	defer resp.Body.Close()
	if DlDecompress {
		writeDecompressed(resp.Body, proj)
	} else {
		writeGzip(resp.Body, proj, filepath.Join(DlDir, proj+".tar.gz"))
	}
	writeManifest(proj, manifest)
}

func appendDl(fileIds []string, manifest []Manifest, proj string) {

	_, err := os.Stat(path.Join(DlDir, proj))
	if err != nil && os.IsNotExist(err) {
		dl(fileIds, manifest, proj)
		return
	}
	manifestMap := make(map[string][]string)
	for _, record := range manifest {
		manifestMap[record.FileID] = []string{record.Md5sum, record.FileName}
	}

	fileIdList, _ := DirList(path.Join(DlDir, proj))
	reFiles := Difference(fileIds, fileIdList)
	wrongCh := make(chan string)
	for _, fId := range fileIdList {
		go func(fId string, manifestMap map[string][]string) {
			fs, _ := FileList(path.Join(DlDir, proj, fId))
			if len(fs) == 0 {
				wrongCh <- fId
			} else {
				md5sum := Md5sum(fs[0])
				if md5sum != manifestMap[fId][0] {
					wrongCh <- fId
				} else {
					wrongCh <- ""
				}
			}
		}(fId, manifestMap)
	}

	for i := 0; i < len(fileIdList); i++ {
		wrong := <-wrongCh
		if wrong == "" {
			continue
		}
		reFiles = append(reFiles, wrong)
	}
	close(wrongCh)

	if len(reFiles) == 0 {
		fmt.Println("Nothing changed. All files are already downloaded.")
		return
	}
	payload, _ := json.Marshal(url.Values{"ids": reFiles})

	req, _ := http.NewRequest("POST", DATA_EP, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Range", "bytes=0-")
	// req.Header.Set("Host", "api.gdc.cancer.gov")
	resp, err := Client.Do(req)
	if err != nil {
		fmt.Println("Error sending request. Check your network connection and proxy settings.")
		os.Exit(1)
	}
	// resp, _ := http.PostForm(DATA_EP, url.Values{"ids": fileIds})
	// fileSize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	defer resp.Body.Close()

	if len(reFiles) > 1 {
		writeDecompressed(resp.Body, proj)
	} else {
		writeGzip(resp.Body, proj, path.Join(DlDir, proj, reFiles[0], manifestMap[reFiles[0]][1]))
	}
	writeManifest(proj, manifest)
}

func fetchInfo(proj string) *Metadata {
	filterJSON, _ := json.Marshal(&Filter{
		Op: "and",
		Content: []Content{
			{
				Op: "in",
				Content: SubContent{
					Field: "cases.project.project_id",
					Value: []string{proj},
				},
			},
			{
				Op: "in",
				Content: SubContent{
					Field: "files.experimental_strategy",
					Value: []string{"RNA-Seq"},
				},
			},
			{
				Op: "in",
				Content: SubContent{
					Field: "files.analysis.workflow_type",
					Value: []string{"HTSeq - Counts"},
				},
			},
		},
	})
	params, _ := query.Values(Params{
		Filters: string(filterJSON),
		Fields: strings.Join([]string{
			"file_name",
			"md5sum",
			"associated_entities.entity_submitter_id",
			// "associated_entities.case_id",
		}, `,`),
		Format: "JSON",
		Size:   "1000000",
	})
	reqUrl := FILES_EP + "?" + params.Encode()

	req, _ := http.NewRequest("GET", reqUrl, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := Client.Do(req)
	if err != nil {
		fmt.Println("Error sending request. Check your network connection and proxy settings.")
		os.Exit(1)
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}
	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("GDC response broken.")
		os.Exit(1)
	}
	var metadata Metadata
	_ = json.Unmarshal(bodyBytes, &metadata)
	metadata.Project = proj
	return &metadata
}

func HandleDl(projects []string) {

	Client = initClient()

	ch := make(chan *Metadata, len(projects))

	for _, proj := range projects {
		go func(proj string) {
			fmt.Printf("%s: Collecting info...\n", proj)
			ch <- fetchInfo(proj)
		}(proj)
	}

	for i := 0; i < len(projects); i++ {
		metadata := <-ch
		files := metadata.Data.Hits
		fileIds := make([]string, 0)
		manifest := make([]Manifest, 0)
		for _, file := range files {
			fileIds = append(fileIds, file.ID)
			manifest = append(manifest, Manifest{
				FileID:   file.ID,
				FileName: file.FileName,
				TCGA:     file.Entities[0].Barcode,
				Md5sum:   file.Md5sum,
			})
		}
		if DlSkip {
			appendDl(fileIds, manifest, metadata.Project)
		} else {
			dl(fileIds, manifest, metadata.Project)
		}

	}
	close(ch)

}
