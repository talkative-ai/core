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

// CompiledDialogRootWithinActor generates the key for a dialog root node within a actor
// Notice that we're not using a node ID. This is because the list of nodes within a actor
// are not readily available, for performance reasons.
func CompiledDialogRootWithinActor(pubID, actorID uint64) string {
	return fmt.Sprintf("%v:e:%v:i",
		CompiledEntity(pubID, models.AEIDActor, actorID),
		models.AEIDDialogNode)
}

// CompiledDialogNodeWithinActor generates the key for a dialog node within a actor
func CompiledDialogNodeWithinActor(pubID, actorID, parentDialogID uint64) string {
	return fmt.Sprintf("%v:e:%v:%v:i",
		CompiledEntity(pubID, models.AEIDActor, actorID),
		models.AEIDDialogNode, parentDialogID)
}

func CompiledActorsWithinZone(pubID, zoneID uint64) string {
	return fmt.Sprintf("%v:e:%v",
		CompiledEntity(pubID, models.AEIDZone, zoneID),
		models.AEIDActor)
}

// CompiledDialogNodeActionBundle generates the key for
// an action bundle within a dialog node
func CompiledDialogNodeActionBundle(pubID, dialogID, bundleID uint64) string {
	return fmt.Sprintf("%v:e:%v:%v",
		CompiledEntity(pubID, models.AEIDDialogNode, dialogID),
		models.AEIDActionBundle, bundleID)
}

// ProjectMetadataStatic generates the key to access the static metadata hash
// Static means these values are not updated after published.
func ProjectMetadataStatic(pubID uint64) string {
	return fmt.Sprintf("%v:%v:m:s",
		compiledNamespaceV1(),
		pubID)
}

// ProjectMetadataDynamic generates the key to access the dynamic metadata hash
// Dynamic means these values may be updated after published.
func ProjectMetadataDynamic(pubID uint64) string {
	return fmt.Sprintf("%v:%v:m:d",
		compiledNamespaceV1(),
		pubID)
}

// GlobalMetaProjects generates the key to access the hash of all published projects
// Mapping project name to project ID
func GlobalMetaProjects() string {
	return fmt.Sprintf("%v:live:projects", compiledNamespaceV1())
}
