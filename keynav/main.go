package keynav

import "fmt"
import "github.com/artificial-universe-maker/go-utilities/models"

func compiledNamespaceV1() string {
	return "c:v1"
}

func CompiledEntity(pubID uint64, entityID models.AEID, uniqueID uint64) string {
	return fmt.Sprintf("%v:%v:e:%v:%v", compiledNamespaceV1(), pubID, entityID, uniqueID)
}

func CompiledDialogRootWithinZone(pubID, zoneID uint64) string {
	return fmt.Sprintf("%v:%v:e:%v:%v:e:%v:i",
		compiledNamespaceV1(), pubID, models.AEIDZone, zoneID, models.AEIDDialogNode)
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
