package models

import (
	"github.com/invopop/jsonschema"
)

type ChangelogChange struct {
	Title       string   `json:"title" jsonschema:"description=The title of the change. Should be succint."`
	Description string   `json:"description" jsonschema:"description=End-user friendly description of the change. Should be more verbose."`
	Impact      string   `json:"impact" jsonschema:"description=The impact of the change. Describe what and how the change affects the user or usage of the software."`
	Commits     []string `json:"commits" jsonschema:"description=List of commit hashes associated with this change. Must have at least one value."`
	Tags        []string `json:"tags" jsonschema:"description=Tags associated with this change"`
}

type ChangelogEntry struct {
	Version string            `json:"version" jsonschema:"description=The version number of the release. Leave as empty string."`
	Date    string            `json:"date" jsonschema:"description=The date of the release. Leave as empty string."`
	FromRef string            `json:"from_ref" jsonschema:"description=The starting commit reference for the changelog entry. Leave as empty string."`
	ToRef   string            `json:"to_ref" jsonschema:"description=The ending commit reference for the changelog entry. Leave as empty string."`
	Changes []ChangelogChange `json:"changes" jsonschema:"description=Generate a list of changes following the schema using the provided git commits and diffs."`
}

func GenerateSchema[T any]() *jsonschema.Schema {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(&v)
	return schema
}

var ChangelogEntrySchema = GenerateSchema[ChangelogEntry]()
