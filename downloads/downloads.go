package downloads

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

func download(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	filename, err := urlToFileName(url)
	if err != nil {
		return "", err
	}
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	return filename, err
}

func urlToFileName(rawurl string) (string, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	return filepath.Base(url.Path), nil
}

func writeZip(outFilename string, filenames []string) error {
	outf, err := os.Create(outFilename)
	if err != nil {
		return err
	}
	zw := zip.NewWriter(outf)
	for _, filename := range filenames {
		w, err := zw.Create(filename)
		if err != nil {
			return err
		}
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(w, f)
		if err != nil {
			return err
		}
	}
	return zw.Close()
}

func Start() {

	var wg sync.WaitGroup
	var urls = []string{
		"https://static.inews24.com/v1/e707be078c033c.jpg",
		"https://spnimage.edaily.co.kr/images/photo/files/NP/S/2022/01/PS22010900030.jpg",
		"https://cafeptthumb-phinf.pstatic.net/MjAyMjAxMDhfMjQ4/MDAxNjQxNjE3NjM0ODE3.c6w90hCHat8FMBCHIyALINmuC-fAKCYyBf49qCkZAjUg.mraZW4sv-A1XCiODr_3PKlxVZKclZAtNqiQK0_94n_0g.JPEG/IMG_9302.jpg",
	}

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			if _, err := download(url); err != nil {
				log.Fatal(err)
			}
		}(url)
	}
	wg.Wait()

	filenames, err := filepath.Glob("*.jpg")
	if err != nil {
		log.Fatal(err)
	}

	err = writeZip("right_face.zip", filenames)

	fmt.Println(filenames)
	for _, v := range filenames {
		os.Remove(v)
	}

	fmt.Println("Done!")

}
