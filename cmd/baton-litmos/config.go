package main

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	apiKeyField       = field.StringField("api-key", field.WithDescription(`API Key`), field.WithRequired(true))
	sourceField       = field.StringField("source", field.WithDescription(`Source`), field.WithRequired(true))
	limitCoursesField = field.StringSliceField("limited-courses", field.WithDescription(`Limit imported sources to a specific list by Course ID`), field.WithRequired(false))
)

var configFields = []field.SchemaField{
	apiKeyField,
	sourceField,
	limitCoursesField,
}

var configRelations = []field.SchemaFieldRelationship{}

var cfg = field.Configuration{
	Fields:      configFields,
	Constraints: configRelations,
}
