package sni

import (
	"context"
	"net"

	"github.com/go-gost/gost/pkg/connector"
	"github.com/go-gost/gost/pkg/logger"
	md "github.com/go-gost/gost/pkg/metadata"
	"github.com/go-gost/gost/pkg/registry"
)

func init() {
	registry.RegiserConnector("sni", NewConnector)
}

type sniConnector struct {
	md     metadata
	logger logger.Logger
}

func NewConnector(opts ...connector.Option) connector.Connector {
	options := &connector.Options{}
	for _, opt := range opts {
		opt(options)
	}

	return &sniConnector{
		logger: options.Logger,
	}
}

func (c *sniConnector) Init(md md.Metadata) (err error) {
	return c.parseMetadata(md)
}

func (c *sniConnector) Connect(ctx context.Context, conn net.Conn, network, address string, opts ...connector.ConnectOption) (net.Conn, error) {
	c.logger = c.logger.WithFields(map[string]interface{}{
		"remote":  conn.RemoteAddr().String(),
		"local":   conn.LocalAddr().String(),
		"network": network,
		"address": address,
	})
	c.logger.Infof("connect %s/%s", address, network)

	return &sniClientConn{Conn: conn, host: c.md.host}, nil
}