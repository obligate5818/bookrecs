package models

import (
	"time"

	"gorm.io/datatypes"
)

type Edition struct {
	ID             uint   `gorm:"primaryKey"`
	Key            string `gorm:"uniqueIndex"`
	Title          string
	Authors        datatypes.JSON
	ISBN13         datatypes.JSON
	Languages      datatypes.JSON
	Pagination     string
	PublishDate    string
	Publishers     datatypes.JSON
	SourceRecords  datatypes.JSON
	Works          datatypes.JSON
	Weight         string
	LatestRevision int
	Revision       int
	Created        time.Time
	LastModified   time.Time
}
