package keynav

import "fmt"
import "github.com/artificial-universe-maker/go-utilities/models"

func CompiledEntities(pubID int, entityID models.AEID) string {
	return fmt.Sprintf("compiled:%i:e:%i", pubID, entityID)
}
