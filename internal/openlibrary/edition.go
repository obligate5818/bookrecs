package openlibrary

import (
	"github.com/obligate5818/bookrecs/internal/models"
	"gorm.io/datatypes"
)

// This struct maps exactly the OpenLibrary edition JSON
type Edition struct {
	Key            string          `json:"key"`
	Title          string          `json:"title"`
	Authors        datatypes.JSON  `json:"authors"`
	ISBN13         datatypes.JSON  `json:"isbn_13"`
	Languages      datatypes.JSON  `json:"languages"`
	Pagination     string          `json:"pagination"`
	PublishDate    string          `json:"publish_date"`
	Publishers     datatypes.JSON  `json:"publishers"`
	SourceRecords  datatypes.JSON  `json:"source_records"`
	Works          datatypes.JSON  `json:"works"`
	Weight         string          `json:"weight"`
	LatestRevision int             `json:"latest_revision"`
	Revision       int             `json:"revision"`
	Created        DateTimeWrapper `json:"created"`
	LastModified   DateTimeWrapper `json:"last_modified"`
}

// Convert Edition to your internal DB model
func (o *Edition) ToInternalModel() *models.Edition {
	return &models.Edition{
		Key:            o.Key,
		Title:          o.Title,
		Authors:        o.Authors,
		ISBN13:         o.ISBN13,
		Languages:      o.Languages,
		Pagination:     o.Pagination,
		PublishDate:    o.PublishDate,
		Publishers:     o.Publishers,
		SourceRecords:  o.SourceRecords,
		Works:          o.Works,
		Weight:         o.Weight,
		LatestRevision: o.LatestRevision,
		Revision:       o.Revision,
		Created:        o.Created.Value,
		LastModified:   o.LastModified.Value,
	}
}
