package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
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
			Md5Sum   string `json:"md5sum"`
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
	FileID string `json:"file_id"`
	TCGA   string `json:"TCGA_barcode"`
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

func download(fileIds []string, manifest []Manifest, proj string) {
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

	err = os.MkdirAll(DlDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error create directories. Check your permissions.")
		os.Exit(1)
	}
	target, _ := os.Create(path.Join(DlDir, proj+".tmp"))
	defer target.Close()

	fmt.Print("\n")
	bar := progressbar.DefaultBytes(-1, "Downloading "+proj)

	_, err = io.Copy(io.MultiWriter(target, bar), resp.Body)
	if err != nil {
		target.Close()
		os.Remove(path.Join(DlDir, proj+".tmp"))
		fmt.Println("Error writing file to directory. Check your permissions.")
		os.Exit(1)
	}
	target.Close()
	_ = os.Rename(path.Join(DlDir, proj+".tmp"), path.Join(DlDir, proj+".tar.gz"))

	_ = os.MkdirAll(path.Join(DlDir, "manifest"), os.ModePerm)
	manifestOut, _ := os.Create(path.Join(DlDir, "manifest", proj+".csv"))
	defer manifestOut.Close()

	manifestWriter := csv.NewWriter(manifestOut)
	_ = manifestWriter.Write([]string{"file_id", "TCGA_manifest"})
	for _, record := range manifest {
		_ = manifestWriter.Write([]string{
			record.FileID, record.TCGA,
		})
	}
	manifestWriter.Flush()

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
			// "file_name",
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
	resp, err := Client.Do(req)
	if err != nil {
		fmt.Println("Error sending request. Check your network connection and proxy settings.")
		os.Exit(1)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("GDC service may not operating normally. Please try again later.")
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
				FileID: file.ID,
				TCGA:   file.Entities[0].Barcode,
			})
		}
		download(fileIds, manifest, metadata.Project)
	}
	close(ch)

}
