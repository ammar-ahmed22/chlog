package models

import (
	"github.com/invopop/jsonschema"
)

type ChangelogChange struct {
	Details    string   `json:"details" jsonschema:"description=Summarized details of the change. Should only be a single sentence starting with a past-tense verb."`
	CommitHash string   `json:"commit_hash" jsonschema:"description=The commit hash associated with this change."`
	Tags       []string `json:"tags" jsonschema:"enum=added,enum=changed,enum=removed,enum=deprecated,enum=security,enum=fixed,description=Tags associated with this change"`
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
