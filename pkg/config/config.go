package config

import "github.com/conductorone/baton-sdk/pkg/field"

var (
	apiKeyField = field.StringField(
		"api-key",
		field.WithRequired(true),
		field.WithDisplayName("API Key"),
		field.WithDescription("API Key"),
		field.WithPlaceholder("Your Litmos API Key"),
		field.WithIsSecret(true),
	)
	sourceField = field.StringField(
		"source",
		field.WithRequired(true),
		field.WithDisplayName("Source"),
		field.WithDescription("Source"),
		field.WithPlaceholder("Your Litmos source"),
	)
	limitCoursesField = field.StringSliceField(
		"limited-courses",
		field.WithDisplayName("Course IDs"),
		field.WithDescription("Limit imported sources to a specific list by Course ID"),
	)
	baseURLField = field.StringField(
		"base-url",
		field.WithDescription("Override the Litmos API base URL (for testing)"),
		field.WithExportTarget(field.ExportTargetCLIOnly),
		field.WithHidden(true),
	)
)

//go:generate go run ./gen
var Config = field.NewConfiguration(
	[]field.SchemaField{
		apiKeyField,
		sourceField,
		limitCoursesField,
		baseURLField,
	},
	field.WithConnectorDisplayName("Litmos"),
	field.WithIconUrl("/static/app-icons/litmos.svg"),
	field.WithHelpUrl("/docs/baton/litmos"),
)
