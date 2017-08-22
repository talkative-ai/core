// Package keynav is a utility that generates consistent Redis keys
package keynav

import "fmt"
import "github.com/artificial-universe-maker/go-utilities/models"

// compiledNamespaceV1 returns Version 1 of the top level compiled namespace
func compiledNamespaceV1() string {
	return "c:v1"
}

// CompiledEntity generates the key for an entity following the standard pattern
// Because all this data is stored in memory, character count is kept to a bare minimum
// And many terms are severely truncated
//
// The pattern is as follows:
// [compiled]:[version]:[published_id]:[data_type]:[entity_type_id]:[entity_unique_id]
// With:
// [compiled]					Always just "c" to designate that this is compiled data
// [version]					Data version
// [published_id] 		64 bit integer ID of the published project
// [data_type]				"e" to designate "entity" or "m" to designate metadata
// [entity_type_id]		Integer designating the entity type. Found in models, as "AEID"
// [entity_unique_id]	64 bit integer ID of the entity
//
// Subentities may exist, and would therefore append to all of this in the same pattern,
// starting with [data_type] etc. etc.
func CompiledEntity(pubID uint64, entityID models.AEID, uniqueID uint64) string {
	return fmt.Sprintf("%v:%v:e:%v:%v", compiledNamespaceV1(), pubID, entityID, uniqueID)
}

// CompiledDialogRootWithinZone generates the key for a dialog root node within a zone
func CompiledDialogRootWithinZone(pubID, zoneID uint64) string {
	return fmt.Sprintf("%v:e:%v:i",
		CompiledEntity(pubID, models.AEIDZone, zoneID), models.AEIDDialogNode)
}

func CompiledDialogNodeWithinZone(pubID, zoneID, parentDialogID uint64) string {
	return fmt.Sprintf("%v:%v:e:%v:%v:e:%v:%v:i",
		compiledNamespaceV1(), pubID, models.AEIDZone, zoneID, models.AEIDDialogNode, parentDialogID)
}

func CompiledDialogNodeActionBundle(pubID, dialogID, bundleID uint64) string {
	return fmt.Sprintf("%v:%v:e:%v:%v:e:%v:%v",
		compiledNamespaceV1(), pubID, models.AEIDDialogNode, dialogID, models.AEIDActionBundle, bundleID)
}

func ProjectMetadataStatic(pubID uint64) string {
	return fmt.Sprintf("%v:%v:m:s",
		compiledNamespaceV1(), pubID)
}

func ProjectMetadataDynamicProperty(pubID uint64, property string) string {
	return fmt.Sprintf("%v:%v:m:d:%v",
		compiledNamespaceV1(), pubID, property)
}

func GlobalMetaProjects() string {
	return fmt.Sprintf("%v:live:projects", compiledNamespaceV1())
}
