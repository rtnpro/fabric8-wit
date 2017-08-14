package resource

import (
	"fmt"
	"net/http"

	"github.com/fabric8-services/fabric8-wit/application"
	"github.com/fabric8-services/fabric8-wit/errors"
	model "github.com/fabric8-services/fabric8-wit/space"
	"github.com/manyminds/api2go"

	uuid "github.com/satori/go.uuid"
)

// SpacesResourceConfiguration is the config interface for SpacesResource
type SpacesResourceConfiguration interface {
	GetCacheControlWorkItems() string
	GetAPIServiceURL() string
}

// SpacesResource the resource for spaces
type SpacesResource struct {
	db     application.DB
	config SpacesResourceConfiguration
}

// NewSpacesResource returns a new SpacesResource
func NewSpacesResource(db application.DB, config SpacesResourceConfiguration) SpacesResource {
	return SpacesResource{
		db:     db,
		config: config,
	}
}

// FindOne finds a Space item by ID
func (res SpacesResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	spaceID, err := uuid.FromString(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, "the space ID is not a valid UUID", http.StatusNotFound)
	}
	s, err := res.db.Spaces().Load(r.Context, spaceID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}
	return &Response{Res: s}, nil
}

// FindAll returns all Space items
func (res SpacesResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	/*
		_, err := login.ContextIdentity(ctx)
		if err != nil {
			return Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusUnauthorized)
		}
	*/
	spaces, _, err := res.db.Spaces().List(r.Context, nil, nil)
	if err != nil {
		return Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}
	return Response{Res: spaces}, nil
}

// PaginatedFindAll returns a page of Space items
func (res SpacesResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	offset, limit, _ := ParsePaging(r)
	spaces, cnt, err := res.db.Spaces().List(r.Context, &offset, &limit)
	if err != nil {
		return uint(0), Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}
	return uint(cnt), Response{Res: spaces}, nil
}

// Create a space item
func (res SpacesResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	space, ok := obj.(model.Space)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.NewBadParameterError("Invalid instance given", nil), "Invalid instance given", http.StatusBadRequest)
	}
	_, err := res.db.Spaces().Create(r.Context, &space)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusBadRequest)
	}
	fmt.Println("Space created", space.ID)
	return &Response{Res: space, Code: http.StatusCreated}, nil
}

// Delete a space item
func (res SpacesResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	spaceID, err := uuid.FromString(id)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusBadRequest)
	}
	err = res.db.Spaces().Delete(r.Context, spaceID)
	return &Response{Code: http.StatusNoContent}, err
}

// Update a space item
func (res SpacesResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	// FIXME: obj is a pointer to model.Space and it's not documented in upstream docs
	space, ok := obj.(*model.Space)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.NewBadParameterError("Invalid instance given", nil), "Invalid instance given", http.StatusBadRequest)
	}
	_, err := res.db.Spaces().Save(r.Context, space)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusBadRequest)
	}
	return &Response{Res: *space, Code: http.StatusNoContent}, nil
}
