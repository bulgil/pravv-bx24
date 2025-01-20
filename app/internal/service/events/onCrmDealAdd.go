package events

import (
	"github.com/bulgil/pravv-bx24/app/internal/domain/models"
	"github.com/bulgil/pravv-bx24/app/package/logger"
)

// При создании сделки из контакта, нужно проверить наличие дубля в соответствующей воронке.
// Если дубль найден, необходимо его удалить и оповестить менеджера, создавшего сделку.
func onCrmDealAdd(log logger.Logger, data models.Data) {
	log.Info("lead created")

	_ = data
}
