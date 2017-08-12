package keynav

import "fmt"
import "github.com/artificial-universe-maker/go-utilities/models"

func CompiledEntities(pubID int, entityID models.AEID, uniqueID string) string {
	return fmt.Sprintf("compiled:%v:e:%v:%v", pubID, entityID, uniqueID)
}
