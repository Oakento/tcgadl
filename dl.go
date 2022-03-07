package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

func filter(proj string) *Filter {
	return &Filter{
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
	}
}

func fetchInfo(proj string, client *http.Client, ch chan map[string]interface{}) {
	filterJSON, _ := json.Marshal(filter(proj))
	params := Params{
		Filters: string(filterJSON),
		Fields:  strings.Join(FIELDS, `,`),
		Format:  "JSON",
		Size:    "1000000",
	}

	paramsStr, _ := query.Values(params)
	reqUrl := FILES_EP + "?" + paramsStr.Encode()

	req, _ := http.NewRequest("GET", reqUrl, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Panic("Error sending request:", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
	}
	metadata := make(map[string]interface{})
	_ = json.Unmarshal(bodyBytes, &metadata)
	metadata["project"] = proj
	ch <- metadata
}

func download(fileIds []string, manifest []map[string]string, dir string, proj string, client *http.Client) {
	payload, _ := json.Marshal(url.Values{"ids": fileIds})

	req, _ := http.NewRequest("POST", DATA_EP, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Range", "bytes=0-")
	// req.Header.Set("Host", "api.gdc.cancer.gov")
	resp, err := client.Do(req)
	if err != nil {
		log.Panic("Error sending request:", err)
	}
	// resp, _ := http.PostForm(DATA_EP, url.Values{"ids": fileIds})
	// fileSize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	defer resp.Body.Close()

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	target, _ := os.Create(path.Join(dir, proj+".tmp"))
	defer target.Close()

	fmt.Print("\n")
	bar := progressbar.DefaultBytes(-1, "Downloading "+proj)

	_, err = io.Copy(io.MultiWriter(target, bar), resp.Body)
	if err != nil {
		target.Close()
		os.Remove(path.Join(dir, proj+".tmp"))
		log.Panic("Error writing file to directory:", err)
	}
	target.Close()
	_ = os.Rename(path.Join(dir, proj+".tmp"), path.Join(dir, proj+".tar.gz"))

	_ = os.MkdirAll(path.Join(dir, "manifest"), os.ModePerm)
	manifestOut, _ := os.Create(path.Join(dir, "manifest", proj+".csv"))
	defer manifestOut.Close()

	manifestWriter := csv.NewWriter(manifestOut)
	_ = manifestWriter.Write([]string{"file_id", "TCGA_manifest"})
	for _, record := range manifest {
		_ = manifestWriter.Write([]string{
			record["file_id"], record["TCGA_barcode"],
		})
	}
	manifestWriter.Flush()

}

func DownloadMany(projects []string, dir string) {
	var client *http.Client
	if HTTP_PROXY != "" {
		proxyURL, err := url.Parse(HTTP_PROXY)
		if err != nil {
			log.Println("Proxy URL is invalid:", err)
			os.Exit(1)
		}
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
	} else {
		client = &http.Client{}
	}

	chs := make([]chan map[string]interface{}, len(projects))
	for i, proj := range projects {
		chs[i] = make(chan map[string]interface{})
		go fetchInfo(proj, client, chs[i])
		fmt.Printf("%s: Collecting metadata...\n", proj)
	}

	for _, ch := range chs {
		metadata := <-ch
		files := metadata["data"].(map[string]interface{})["hits"].([]interface{})
		fileIds := make([]string, 0)
		manifest := make([]map[string]string, 0)
		for _, file := range files {
			fileIds = append(fileIds, file.(map[string]interface{})["id"].(string))
			manifest = append(manifest, map[string]string{
				"file_id":      file.(map[string]interface{})["id"].(string),
				"TCGA_barcode": file.(map[string]interface{})["associated_entities"].([]interface{})[0].(map[string]interface{})["entity_submitter_id"].(string),
			})
		}
		download(fileIds, manifest, dir, metadata["project"].(string), client)
	}

}
