package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	ApiKeyField = field.StringField(
		"api-key",
		field.WithDescription("API Key"),
		field.WithRequired(true),
	)
	SourceField = field.StringField(
		"source",
		field.WithDescription("Source"),
		field.WithRequired(true),
	)
	LimitedCoursesField = field.StringSliceField(
		"limited-courses",
		field.WithDescription("Limit imported sources to a specific list by Course ID"),
	)

	ConfigurationFields = []field.SchemaField{
		ApiKeyField,
		SourceField,
		LimitedCoursesField,
	}
)

//go:generate go run ./gen
var Config = field.NewConfiguration(ConfigurationFields)
