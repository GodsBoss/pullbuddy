package pullbuddy

import "io"

type Client struct {
	Addr string
	Out  io.Writer
}

func (client *Client) Status() error {
	return nil
}

func (client *Client) Schedule(id string) error {
	return nil
}
