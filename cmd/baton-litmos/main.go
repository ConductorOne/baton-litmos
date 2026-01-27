package main

import (
	"context"
	"fmt"
	"os"

	cfg "github.com/conductorone/baton-litmos/pkg/config"
	"github.com/conductorone/baton-litmos/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/connectorrunner"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

var version = "dev"

func main() {
	ctx := context.Background()

	_, cmd, err := config.DefineConfiguration(
		ctx,
		"baton-litmos",
		getConnector,
		cfg.Config,
		connectorrunner.WithDefaultCapabilitiesConnectorBuilder(&connector.LitmosConnector{}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, c *cfg.Litmos) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)
	if err := cfg.ValidateConfig(c); err != nil {
		return nil, err
	}

	cb, err := connector.New(ctx, c.ApiKey, c.Source, c.LimitedCourses)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	conn, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	return conn, nil
}
