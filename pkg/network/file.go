package network

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/nice-pink/goutil/pkg/log"
)

func DownloadHttp(url string, filepath string) error {
	log.Info("http download:", url)

	out, err := os.Create(filepath)
	if err != nil {
		log.Err(err, "Could not create file.")
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		log.Err(err, "Could not request url.")
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		log.Error("bad status: %s", resp.Status)
		return errors.New("bad status: " + resp.Status)
	}

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Err(err, "Could not copy data to file.")
		return err
	}
	log.Info("Downloaded file with", n, "bytes")
	return nil
}

func UploadHttp(url string, filepath string, contentType string) error {
	log.Info("http upload", filepath, "to", url, "with content type", contentType)

	file, err := os.Open(filepath)
	if err != nil {
		log.Err(err, "Could not read file", filepath)
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, file)
	if err != nil {
		log.Err(err, "Could not create put request.")
		return err
	}
	req.Header.Add("Content-Type", contentType)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Err(err, "Could not send request.")
		return err
	}
	defer res.Body.Close()
	return nil
}
