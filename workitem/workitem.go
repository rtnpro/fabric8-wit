package workitem

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/fabric8-services/fabric8-wit/log"
	"github.com/manyminds/api2go/jsonapi"

	uuid "github.com/satori/go.uuid"
)

// WorkItem the model structure for the work item.
type WorkItem struct {
	// unique id per installation (used for references at the DB level)
	ID uuid.UUID
	// unique number per _space_
	Number int
	// ID of the type of this work item
	Type uuid.UUID
	// Version for optimistic concurrency control
	Version int
	// ID of the space to which this work item belongs
	SpaceID uuid.UUID
	// The field values, according to the field type
	Fields map[string]interface{}
	// optional, private timestamp of the latest addition/removal of a relationship with this workitem
	// this field is used to generate the `ETag` and `Last-Modified` values in the HTTP responses and conditional requests processing
	relationShipsChangedAt *time.Time
}

// WICountsPerIteration counting work item states by iteration
type WICountsPerIteration struct {
	IterationID string `gorm:"column:iterationid"`
	Total       int
	Closed      int
}

// GetETagData returns the field values to use to generate the ETag
func (wi WorkItem) GetETagData() []interface{} {
	return []interface{}{wi.ID, wi.Version, wi.relationShipsChangedAt}
}

// GetLastModified returns the last modification time
func (wi WorkItem) GetLastModified() time.Time {
	var lastModified *time.Time // default value
	if updatedAt, ok := wi.Fields[SystemUpdatedAt].(time.Time); ok {
		lastModified = &updatedAt
	}
	// also check the optional 'relationShipsChangedAt' field
	if wi.relationShipsChangedAt != nil && (lastModified == nil || wi.relationShipsChangedAt.After(*lastModified)) {
		lastModified = wi.relationShipsChangedAt
	}

	log.Debug(nil, map[string]interface{}{"wi_id": wi.ID}, "Last modified value: %v", lastModified)
	return *lastModified
}

// JSONAPI encoding functions

// MarshalJSON is the custom Marshaller for dealing with the variable fields in attributes.
func (wi WorkItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(wi.Fields)
}

// GetID returns the ID for marshalling to json.
func (wi WorkItem) GetID() string {
	return wi.ID.String()
}

// SetID sets the ID for marshalling to json.
func (wi WorkItem) SetID(id string) error {
	wi.ID, _ = uuid.FromString(id)
	return nil
}

// GetName returns the entity type name for marshalling to json.
func (wi WorkItem) GetName() string {
	return "workitems"
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (wi WorkItem) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type:        "areas",
			Name:        "area",
			IsNotLoaded: false, // we want to have the data field
		},
		{
			Type:        "identities",
			Name:        "assignees",
			IsNotLoaded: false, // we want to have the data field
		},
		{
			Type:        "workitemtypes",
			Name:        "baseType",
			IsNotLoaded: false, // we want to have the data field
		},
		{
			Type:        "workitems",
			Name:        "children",
			IsNotLoaded: true, // we do not want to have the data field
		},
		{
			Type:        "comments",
			Name:        "comments",
			IsNotLoaded: true, // we do not want to have the data field
		},
		{
			Type:        "identities",
			Name:        "creator",
			IsNotLoaded: false, // we want to have the data field
		},
		{
			Type:        "iterations",
			Name:        "iteration",
			IsNotLoaded: false, // we want to have the data field
		},
		{
			Type:        "spaces",
			Name:        "space",
			IsNotLoaded: false, // we want to have the data field
		},
		{
			Type:        "workitemlinktypes",
			Name:        "source-link-types",
			IsNotLoaded: true, // we do not want to have the data field
		},
		{
			Type:        "workitemlinktypes",
			Name:        "target-link-types",
			IsNotLoaded: true, // we do not want to have the data field
		},
		{
			Type:        "links",
			Name:        "links",
			IsNotLoaded: true, // we do not want to have the data field
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (wi WorkItem) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{
		jsonapi.ReferenceID{
			ID:   wi.SpaceID.String(),
			Type: "spaces",
			Name: "space",
		},
		jsonapi.ReferenceID{
			ID:   wi.Fields["system.area"].(string),
			Type: "areas",
			Name: "area",
		},
		jsonapi.ReferenceID{
			ID:   wi.Fields["system.creator"].(string),
			Type: "identities",
			Name: "creator",
		},
		jsonapi.ReferenceID{
			ID:   wi.Fields["system.iteration"].(string),
			Type: "iterations",
			Name: "iteration",
		},
		jsonapi.ReferenceID{
			ID:   wi.Type.String(),
			Type: "workitemtypes",
			Name: "baseType",
		},
	}
	assignees := reflect.ValueOf(wi.Fields["system.assignees"])
	for i := 0; i < assignees.Len(); i++ {
		result = append(result, jsonapi.ReferenceID{
			ID:   assignees.Index(i).Interface().(uuid.UUID).String(),
			Type: "identities",
			Name: "assignees",
		})
	}
	return result
}

// GetCustomLinks returns the custom links, namely the self link.
func (wi WorkItem) GetCustomLinks(linkURL string) jsonapi.Links {
	links := jsonapi.Links{
		"self":            jsonapi.Link{linkURL, nil},
		"sourceLinkTypes": jsonapi.Link{linkURL + "/source-link-types", nil},
		"targetLinkTypes": jsonapi.Link{linkURL + "/target-link-types", nil},
	}
	return links
}

/*
// GetCustomMeta returns the custom meta.
// TODO this looks like it is being called 10 times for each serialization. Check why!
func (wi WorkItem) GetCustomMeta(linkURL string) jsonapi.Metas {
	hasChildren, _ := utils.DatabaseMetaService.HasChildren(wi.GetName(), wi.ID)
	meta := map[string]map[string]interface{}{
		"children": map[string]interface{}{
			"hasChildren": hasChildren,
		},
	}
	return meta
}
*/
