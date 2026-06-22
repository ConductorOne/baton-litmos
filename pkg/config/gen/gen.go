package main

import (
	cfg "github.com/conductorone/baton-litmos/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("litmos", cfg.Config)
}
