package main

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	apiKeyField = field.StringField("api-key", field.WithDescription(`API Key`), field.WithRequired(true))
	sourceField = field.StringField("source", field.WithDescription(`Source`), field.WithRequired(true))
)

var configFields = []field.SchemaField{
	apiKeyField,
	sourceField,
}

var configRelations = []field.SchemaFieldRelationship{}

var cfg = field.Configuration{
	Fields:      configFields,
	Constraints: configRelations,
}
