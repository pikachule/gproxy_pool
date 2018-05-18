package client

import (
	_ "fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Sources []string

type Client struct {
	SourcePath string
}

var sources Sources

func (c *Client) sourcesLoad() error {
	files, err := filepath.Glob(c.SourcePath)
	if err != nil {
		return err
	}
	if len(files) < 1 {
		return errors.New("Empty directory")
	}
	for _, fn := range files {
		fd, _ := ioutil.ReadFile(fn)
		l := strings.Split(string(fd), `\n`)
		if len(l) > 0 {
			sources = append(sources, l...)
		}
	}
	if len(sources) < 1 {
		return errors.New("Sources cannot be Empty")
	}
	return nil
}

func (c *Client) Start() error {
	err := c.sourcesLoad()
	if err != nil {
		return err
	}
	return nil
}
