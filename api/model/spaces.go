package model

/*

//Space the Space type of resource to (un)marshall in the JSON-API requests/responses
type Space struct {
	ID          string     `jsonapi:"primary,spaces"`
	Name        string     `jsonapi:"attr,name"`
	Description string     `jsonapi:"attr,description"`
	BackLog     []WorkItem `jsonapi:"relation,backlog"`
	BackLogSize int        // carried in this struct to be exposed as the `meta/count` attribute in the `workitems`` links
}

// JSONAPILinks returns the links to the space
func (s Space) JSONAPILinks() *jsonapi.Links {
	config := configuration.Get()
	return &jsonapi.Links{
		"self": jsonapi.Link{
			Href: fmt.Sprintf("%[1]s/api/spaces/%[2]s", config.GetAPIServiceURL(), s.ID),
		},
	}
}

//JSONAPIRelationshipLinks is invoked for each relationship defined on the Space struct when marshaled
func (s Space) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	config := configuration.Get()
	if relation == "backlog" {
		return &jsonapi.Links{
			"related": jsonapi.Link{
				Href: fmt.Sprintf("%[1]s/api/spaces/%[2]s/backlog", config.GetAPIServiceURL(), s.ID),
			},
		}
	}
	return nil
}

// JSONAPIRelationshipMeta Invoked for each relationship defined on the Post struct when marshaled
func (s Space) JSONAPIRelationshipMeta(relation string) *jsonapi.Meta {
	if relation == "backlog" {
		return &jsonapi.Meta{
			"totalCount": s.BackLogSize,
		}
	}
	return nil
}

// Space represents a Space on the domain and db layer
type Space struct {
	gormsupport.Lifecycle
	ID          uuid.UUID `json:"id"`
	Version     int       `json:"version"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerId     uuid.UUID `sql:"type:uuid" json:"-"` // Belongs To Identity
}

// Ensure Fields implements the Equaler interface
var _ convert.Equaler = Space{}
var _ convert.Equaler = (*Space)(nil)

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (p Space) GetID() string {
	return p.ID.String()
}

// SetID to satisfy jsonapi.MarshalIdentifier interface
func (p Space) SetID(id string) error {
	var err error
	p.ID, err = uuid.FromString(id)
	return err
}

// Equal returns true if two Space objects are equal; otherwise false is returned.
func (p Space) Equal(u convert.Equaler) bool {
	other, ok := u.(Space)
	if !ok {
		return false
	}
	lfEqual := p.Lifecycle.Equal(other.Lifecycle)
	if !lfEqual {
		return false
	}
	if p.Version != other.Version {
		return false
	}
	if p.Name != other.Name {
		return false
	}
	if p.Description != other.Description {
		return false
	}
	if !uuid.Equal(p.OwnerId, other.OwnerId) {
		return false
	}
	return true
}

// GetETagData returns the field values to use to generate the ETag
func (p Space) GetETagData() []interface{} {
	return []interface{}{p.ID, p.Version}
}

// GetLastModified returns the last modification time
func (p Space) GetLastModified() time.Time {
	return p.UpdatedAt
}

// TableName overrides the table name settings in Gorm to force a specific table name
// in the database.
func (p Space) TableName() string {
	return "spaces"
}
*/
