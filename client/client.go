package client

import (
	"errors"
	"fmt"
	"gproxy_pool/request"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

type Sources []string

type Client struct {
	SourcePath string
}

var sources Sources
var sourcesCounter int

func (c *Client) sourcesLoad() error {
	files, err := filepath.Glob(c.SourcePath)
	if err != nil {
		return err
	}
	if len(files) < 1 {
		return errors.New("Empty directory")
	}
	for _, fn := range files {
		fmt.Printf(""+
			"Loading file : %v\n", fn)
		fd, _ := ioutil.ReadFile(fn)
		l := strings.Split(string(fd), "\n")
		if len(l) > 0 {
			sources = append(sources, l...)
		}
	}
	if len(sources) < 1 {
		return errors.New("Sources cannot be Empty")
	}
	sourcesCounter = len(sources)
	return nil
}

func (c *Client) Start() error {
	err := c.sourcesLoad()
	if err != nil {
		return err
	}
	ch := make(chan string)
	chc := make(chan int)
	fmt.Println(sourcesCounter)
	go func(c chan int, counter int) {
		c <- counter
	}(chc, sourcesCounter)
	for _, _url := range sources {
		go request.Get(_url, ch, chc)
	}
	time.Sleep(time.Second * 5)
	for msg := range ch {
		c := <-chc
		fmt.Println(c)
		if c < 1 {
			close(chc)
			close(ch)
			break
		}
		fmt.Println(msg)
	}
	return nil
}
