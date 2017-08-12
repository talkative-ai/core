package keynav

import "fmt"
import "github.com/artificial-universe-maker/go-utilities/models"

func CompiledEntity(pubID uint64, entityID models.AEID, uniqueID uint64) string {
	return fmt.Sprintf("compiled:%v:e:%v:%v", pubID, entityID, uniqueID)
}

func CompiledDialogRootWithinZone(pubID, zoneID uint64, dialogEntry string) string {
	return fmt.Sprintf("compiled:%v:e:%v:%v:e:%v:%v",
		pubID, models.AEIDZone, zoneID, models.AEIDDialogNode, dialogEntry)
}

func CompiledDialogNodeWithinZone(pubID, zoneID, parentDialogID uint64, dialogEntry string) string {
	return fmt.Sprintf("compiled:%v:e:%v:%v:e:%v:%v:%v",
		pubID, models.AEIDZone, zoneID, models.AEIDDialogNode, parentDialogID, dialogEntry)
}

func CompiledDialogNodeActionBundle(pubID, dialogID, bundleID uint64) string {
	return fmt.Sprintf("compiled:%v:e:%v:%v:e:%v:%v",
		pubID, models.AEIDDialogNode, dialogID, models.AEIDActionBundle, bundleID)
}
