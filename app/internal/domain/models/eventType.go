package models

type EventType string

const (
	EventTypeOnCrmDealAdd            = "ONCRMDEALADD"
	EventTypeOnCrmDealUpdate         = "ONCRMDEALUPDATE"
	EventTypeOnCrmDealDelete         = "ONCRMDEALDELETE"
	EventTypeOnCrmDealMoveToCategory = "EVENTONCRMDEALMOVETOCATEGORY"

	EventTypeOnCrmContactAdd    = "ONCRMCONTACTADD"
	EventTypeOnCrmContactUpdate = "ONCRMCONTACTUPDATE"
	EventTypeOnCrmContactDelete = "ONCRMCONTACTDELETE"
)
