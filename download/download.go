package download

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Download struct {
	Link        string
	Size        int
	Concurrency int
	Name        string
	Path        string
	Chunks      map[int][]byte
	StartTime   time.Time
}

func (d *Download) send(requestType string) (*http.Request, error) {
	req, err := http.NewRequest(
		requestType,
		d.Link,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return req, nil
}

func (d *Download) Connect() (int, bool, error) {
	req, err := d.send("HEAD")
	if err != nil {
		return -1, false, err
	}

	if err != nil {
		return -1, false, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, false, err
	}
	defer resp.Body.Close()

	if err != nil {
		return -1, false, err
	}

	fmt.Printf("Got response %v\n", resp.StatusCode)

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return -1, false, err
	}
	_, params, err := mime.ParseMediaType((resp.Header.Get("Content-Disposition")))

	if err == nil {
		filename := params["filename"]
		d.Name = filename
		fmt.Println(filename)
	}

	if resp.Header.Get("Accept-Ranges") == "bytes" {
		return size, true, nil
	}

	return size, false, nil
}

func (d *Download) Download() error {
	output, err := os.Create(d.Name)

	if err != nil {
		return err
	}
	defer output.Close()

	file, err := http.Get(d.Link)
	if err != nil {
		return err
	}
	defer file.Body.Close()

	if file.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", file.Status)
	}

	_, err = io.Copy(output, file.Body)

	if err != nil {
		return err
	}

	return nil
}
