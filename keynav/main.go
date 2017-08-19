package keynav

import "fmt"
import "github.com/artificial-universe-maker/go-utilities/models"

func CompiledEntity(pubID uint64, entityID models.AEID, uniqueID uint64) string {
	return fmt.Sprintf("c:%v:e:%v:%v", pubID, entityID, uniqueID)
}

func CompiledDialogRootWithinZone(pubID, zoneID uint64) string {
	return fmt.Sprintf("c:%v:e:%v:%v:e:%v:i",
		pubID, models.AEIDZone, zoneID, models.AEIDDialogNode)
}

func CompiledDialogNodeWithinZone(pubID, zoneID, parentDialogID uint64) string {
	return fmt.Sprintf("c:%v:e:%v:%v:e:%v:%v:i",
		pubID, models.AEIDZone, zoneID, models.AEIDDialogNode, parentDialogID)
}

func CompiledDialogNodeActionBundle(pubID, dialogID, bundleID uint64) string {
	return fmt.Sprintf("c:%v:e:%v:%v:e:%v:%v",
		pubID, models.AEIDDialogNode, dialogID, models.AEIDActionBundle, bundleID)
}

func ProjectMetadataStaticProperty(pubID uint64, property string) string {
	return fmt.Sprintf("c:%v:m:s:%v",
		pubID, property)
}

func ProjectMetadataDynamicProperty(pubID uint64, property string) string {
	return fmt.Sprintf("c:%v:m:d:%v",
		pubID, property)
}
