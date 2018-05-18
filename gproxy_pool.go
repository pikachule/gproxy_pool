package gproxy_pool

import (
	//"fmt"
	"gproxy_pool/client"
)

type Options struct {
	SourcePath string
}

func New(opts Options) (client.Client, error) {
	c := client.Client{opts.SourcePath}
	c.Start()
	return c, nil
}
