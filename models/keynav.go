// Package keynav is a utility that generates consistent Redis keys
package models

import (
	"fmt"
)

// compiledNamespaceV1 returns Version 1 of the top level compiled namespace
func compiledNamespaceV1() string {
	return "c:v1"
}

// KeynavCompiledEntity generates the key for an entity following the standard pattern
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
func KeynavCompiledEntity(pubID uint64, entityID AEID, uniqueID uint64) string {
	return fmt.Sprintf("%v:%v:e:%v:%v", compiledNamespaceV1(), pubID, entityID, uniqueID)
}

// KeynavCompiledDialogRootWithinActor generates the key for a dialog root node within a actor
// Notice that we're not using a node ID. This is because the list of nodes within a actor
// are not readily available, for performance reasons.
func KeynavCompiledDialogRootWithinActor(pubID, actorID uint64) string {
	return fmt.Sprintf("%v:e:%v:i",
		KeynavCompiledEntity(pubID, AEIDActor, actorID),
		AEIDDialogNode)
}

// KeynavCompiledDialogNodeWithinActor generates the key for a dialog node within a actor
func KeynavCompiledDialogNodeWithinActor(pubID, actorID, parentDialogID uint64) string {
	return fmt.Sprintf("%v:e:%v:%v:i",
		KeynavCompiledEntity(pubID, AEIDActor, actorID),
		AEIDDialogNode, parentDialogID)
}

func KeynavCompiledActorsWithinZone(pubID, zoneID uint64) string {
	return fmt.Sprintf("%v:e:%v",
		KeynavCompiledEntity(pubID, AEIDZone, zoneID),
		AEIDActor)
}

// KeynavCompiledDialogNodeActionBundle generates the key for
// an action bundle within a dialog node
func KeynavCompiledDialogNodeActionBundle(pubID, dialogID, bundleID uint64) string {
	return fmt.Sprintf("%v:e:%v:%v",
		KeynavCompiledEntity(pubID, AEIDDialogNode, dialogID),
		AEIDActionBundle, bundleID)
}

// KeynavProjectMetadataStatic generates the key to access the static metadata hash
// Static means these values are not updated after published.
func KeynavProjectMetadataStatic(pubID uint64) string {
	return fmt.Sprintf("%v:%v:m:s",
		compiledNamespaceV1(),
		pubID)
}

// KeynavProjectMetadataDynamic generates the key to access the dynamic metadata hash
// Dynamic means these values may be updated after published.
func KeynavProjectMetadataDynamic(pubID uint64) string {
	return fmt.Sprintf("%v:%v:m:d",
		compiledNamespaceV1(),
		pubID)
}

// KeynavGlobalMetaProjects generates the key to access the hash of all published projects
// Mapping project name to project ID
func KeynavGlobalMetaProjects() string {
	return fmt.Sprintf("%v:live:projects", compiledNamespaceV1())
}

// KeynavParseFromKeyBundleID TODO: Consider if this can be helpful for event sourcing
// i.e. capture every action bundle that mutates states
func KeynavParseFromKeyBundleID(key string) string {
	fmt.Println(key)
	return ""
}

// KeynavCompiledTriggerActionBundle generates the key for
// an action bundle within a trigger
func KeynavCompiledTriggerActionBundle(pubID, zoneID, triggerType, bundleID uint64) string {
	return fmt.Sprintf("%v:e:%v:%v:e%v:%v",
		KeynavCompiledEntity(pubID, AEIDZone, zoneID),
		AEIDTrigger, triggerType,
		AEIDActionBundle, bundleID)
}

// KeynavCompiledTriggersWithinZone generates a key to a hash of all triggers within the zone
// and their keys therein. Each trigger has an associated action bundle with can be accessed
// via another read operation.
func KeynavCompiledTriggersWithinZone(pubID, zoneID uint64) string {
	return fmt.Sprintf("%v:e:%v",
		KeynavCompiledEntity(pubID, AEIDZone, zoneID),
		AEIDTrigger)
}
