package main

import (
	"flag"
	"os"
	"net/http"
	"io"
	"archive/zip"
	"context"
	"cloud.google.com/go/storage"
	"fmt"
	"log"
	"sync"
)

var (
	bucketName = flag.String("bucket-name", "deeplesion-data", "the name of the GCP bucket you want to upload to")
	parallel = flag.Bool("parallel", false, "Download and upload data in parallel, generally requires more disk space")
	removeFiles = flag.Bool("remove-files", false, "remove each file after download and upload (only if parallel=false)")
)

func main() {
	flag.Parse()
	if *parallel {
		BeginParallel()
	} else {
		Begin()
	}

}

func Begin() {
	for i, url := range DownloadURLs {
		fn := fmt.Sprintf("Images_png_%02d.zip", i + 1)
		FetchUploadAndHandleFile(fn, url, *bucketName, nil)
		if *removeFiles {
			os.Remove(fn)
		}
	}
}

func BeginParallel() {
	var wg sync.WaitGroup
	for i, url := range DownloadURLs {
		fn := fmt.Sprintf("Images_png_%02d.zip", i + 1)
		wg.Add(1)
		go FetchUploadAndHandleFile(fn, url, *bucketName, &wg)
	}
	wg.Wait()
}

func FetchUploadAndHandleFile(filename, url, bucketName string, wg *sync.WaitGroup) error {
	log.Printf("Starting download of %s\n", filename)
	FetchFile(filename, url)
	log.Printf("Download of file %s complete, begining unzip and upload to GCP\n", filename)
	UnzipAndUploadFiles(filename, bucketName)
	if wg != nil {
		wg.Done()
	}
	return nil
}

func UnzipAndUploadFiles(filename, bucketName string) error {
	r, err := zip.OpenReader(filename)
	if err != nil {
		log.Printf("Unable to open zip %s\n", filename)
		return err
	}
	defer r.Close()

	// Setup connection to GCP bucket
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	bkt := client.Bucket(bucketName)

	for _, f := range r.File {
		imageFile, err := f.Open()
		if err != nil {
			return err
		}

		fmt.Printf("Uploading %s\n", f.Name)
		imageObj := bkt.Object(f.Name)
		w := imageObj.NewWriter(context.Background())
		_, err = io.Copy(w, imageFile)
		if err != nil {
			w.Close()
			imageFile.Close()
			log.Println("error copying to gcp")
		}

		w.Close()
		imageFile.Close()
	}

	return nil
}

func FetchFile(filename string, url string) error {
	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
