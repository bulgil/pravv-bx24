package events

import (
	"github.com/bulgil/pravv-bx24/app/internal/domain/models"
	"github.com/bulgil/pravv-bx24/app/package/logger"
)

func EventRoute(log logger.Logger, e models.Event) {
	eType := e.EventType
	eData := e.Data

	switch eType {
	case models.EventTypeOnCrmDealAdd:
		onCrmDealAdd(log, eData)
	}
}
