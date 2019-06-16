package pullbuddy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Addr string
	Out  io.Writer
}

func (client *Client) Status() error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/list", orDefaultAddr(client.Addr, DefaultClientAddress)), nil)
	if err != nil {
		return err
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	images := &listResponse{}
	err = json.NewDecoder(resp.Body).Decode(images)
	if err != nil {
		return err
	}
	if len(images.Images) == 0 {
		fmt.Fprintf(client.Out, "no images found\n")
		return nil
	}
	for i := range images.Images {
		fmt.Fprintf(client.Out, "%40s %10s %s\n", images.Images[i].ID, images.Images[i].Status, images.Images[i].Error)
	}
	return nil
}

func (client *Client) Schedule(id string) error {
	reqBodyStruct := scheduleRequest{
		ImageID: id,
	}
	reqBodyData, err := json.Marshal(reqBodyStruct)
	if err != nil {
		return err
	}
	reqBody := bytes.NewReader(reqBodyData)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/schedule", orDefaultAddr(client.Addr, DefaultClientAddress)), reqBody)
	if err != nil {
		return err
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got unexpected HTTP status %s\n", resp.Status)
	}
	fmt.Fprintf(client.Out, "scheduled Docker image %s\n", id)
	return nil
}

const DefaultClientAddress = "localhost:30666"
