package resource

import (
	"net/http"

	"github.com/fabric8-services/fabric8-wit/application"
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
	ctx := r.Context
	spaceID, err := uuid.FromString(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, "the space ID is not a valid UUID", http.StatusNotFound)
	}
	s, err := res.db.Spaces().Load(ctx, spaceID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}
	return &Response{Res: s}, nil
}

// FindAll returns all Space items
func (res SpacesResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	ctx := r.Context
	/*
		_, err := login.ContextIdentity(ctx)
		if err != nil {
			return Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusUnauthorized)
		}
	*/
	spaces, _, err := res.db.Spaces().List(ctx, nil, nil)
	if err != nil {
		return Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusInternalServerError)
	}
	return Response{Res: spaces}, nil
}

/*
//GetByID gets a space resource by its ID
func (r SpacesResource) GetByID(ctx *gin.Context) {
	spaceID, err := uuid.FromString(ctx.Param("spaceID"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errs.Wrapf(err, "the space ID is not a valid UUID"))
	}
	s, err := r.db.Spaces().Load(ctx, spaceID)
	if err != nil {
		//TODO: retrieve the correct HTTP status for the given err
		ctx.AbortWithError(http.StatusInternalServerError, errs.Wrapf(err, "error while fetching the space with id=%s", spaceID.String()))
	}
	// convert the business-domain 'space' into a jsonapi-model space
	result := model.Space{
		ID:          s.ID.String(),
		Name:        s.Name,
		Description: s.Description,
		BackLogSize: 10,
	}

	// marshall the result into a JSON-API compliant doc
	ctx.Status(http.StatusOK)
	ctx.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(ctx.Writer, &result); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errs.Wrapf(err, "error while fetching the space with id=%s", spaceID.String()))
	}
}

// List runs the list action.
func (r SpacesResource) List(ctx *gin.Context) {
	_, err := login.ContextIdentity(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
	}
	pageNumber := GetQueryParamAsString(ctx, "page[number]")
	pageLimit, err := GetQueryParamAsInt(ctx, "page[limit]")
	offset, limit := computePagingLimits(pageNumber, pageLimit)
	spaces, cnt, err := r.db.Spaces().List(ctx, &offset, &limit)
	results := []*model.Space{}
	for _, s := range spaces {
		results = append(results, &model.Space{
			ID:          s.ID.String(),
			Name:        s.Name,
			Description: s.Description,
			BackLogSize: 10,
		})
	}
	p, err := jsonapi.Marshal(results)
	payload, ok := p.(*jsonapi.ManyPayload)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, errs.Wrap(err, "error while preparing the response payload"))
	}
	payload.Links = &jsonapi.Links{}
	first, prev, next, last := getPagingLinks(payload.Links, fmt.Sprintf("%[1]s/api/spaces/", r.config.GetAPIServiceURL()), int(cnt), offset, limit, 10, "")
	payload.Meta = &jsonapi.Meta{
		"total-count": cnt,
	}
	payload.Links = &jsonapi.Links{
		"first": first,
		"prev":  prev,
		"next":  next,
		"last":  last,
	}

	ctx.Status(http.StatusOK)
	ctx.Header("Content-Type", jsonapi.MediaType)
	if err := json.NewEncoder(ctx.Writer).Encode(payload); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errs.Wrapf(err, "Error while fetching the spaces"))
	}
}*/
