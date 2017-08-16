package resource

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fabric8-services/fabric8-wit/application"
	"github.com/fabric8-services/fabric8-wit/errors"
	query "github.com/fabric8-services/fabric8-wit/query/simple"
	"github.com/fabric8-services/fabric8-wit/search"
	"github.com/manyminds/api2go"
	uuid "github.com/satori/go.uuid"
)

const (
	SpaceQueryParam               = "spacesID"
	FilterQueryParam              = "filter"                // a query language expression restricting the set of found work items
	ExpressionFilterQueryParam    = "filter[expression]"    // accepts query in JSON format and redirects to /api/search? API". Example: `{$AND: [{"space": "f73988a2-1916-4572-910b-2df23df4dcc3"}, {"state": "NEW"}]}`
	AssigneeFilterQueryParam      = "filter[assignee]"      // Work Items assigned to the given user
	IterationFilterQueryParam     = "filter[iteration]"     // IterationID to filter work items
	WorkItemTypeFilterQueryParam  = "filter[workitemtype]"  // ID of work item type to filter work items by
	WorkItemStateFilterQueryParam = "filter[workitemstate]" // work item state to filter work items by
	AreaFilterQueryParam          = "filter[area]"          // AreaID to filter work items
	ParentExistsFilterQueryParam  = "filter[parentexists]"  //if false list work items without any parent
	PageLimitQueryParam           = "page[limit]"           // Paging size
	PageOffsetQueryParam          = "page[offset]"          // Paging size
	none                          = "none"
)

// WorkItemControllerConfig the config interface for the WorkitemController
type WorkItemsResourceConfiguration interface {
	GetCacheControlWorkItems() string
	GetAPIServiceURL() string
}

// WorkItemsResource the resource for work items
type WorkItemsResource struct {
	db     application.DB
	config WorkItemsResourceConfiguration
}

// NewWorkItemsResource returns a new WorkItemsResource
func NewWorkItemsResource(db application.DB, config WorkItemsResourceConfiguration) WorkItemsResource {
	return WorkItemsResource{db: db, config: config}
}

type WorkItemsResourceListContext struct {
	api2go.APIContexter
	SpaceID             uuid.UUID `gin:"param,spaceID"`
	Filter              *string   `gin:"query,filter"`
	ExpressionFilter    *string
	AssigneeFilter      *string
	IterationFilter     *string
	WorkitemTypeFilter  *uuid.UUID
	AreaFilter          *string
	WorkitemStateFilter *string
	ParentExistsFilter  *bool
	PageOffset          *int
	PageLimit           *int
}

// Note: this kind of function could be generated, based on struct tags on WorkItemsResourceListContext fields.
func NewWorkItemsResourceListContext(r api2go.Request) (*WorkItemsResourceListContext, error) {
	spacesID, ok := r.QueryParams[SpaceQueryParam]
	if !ok {
		return nil, errors.NewBadParameterError("Invalid spaceid", nil)
	}
	spaceID, err := uuid.FromString(spacesID[0])
	if err != nil {
		return nil, err
	}

	var filter, expressionFilter, assigneeFilter, iterationFilter, areaFilter, workItemStateFilter string
	var workItemTypeFilter uuid.UUID
	var parentExistsFilter bool
	var pageOffset, pageLimit int

	for key, value := range r.QueryParams {
		switch key {
		case FilterQueryParam:
			filter = value[0]
		case ExpressionFilterQueryParam:
			expressionFilter = value[0]
		case AssigneeFilterQueryParam:
			assigneeFilter = value[0]
		case IterationFilterQueryParam:
			iterationFilter = value[0]
		case WorkItemTypeFilterQueryParam:
			workItemTypeFilter, err = uuid.FromString(value[0])
			if err != nil {
				return nil, err
			}
		case AreaFilterQueryParam:
			areaFilter = value[0]
		case WorkItemStateFilterQueryParam:
			workItemStateFilter = value[0]
		case ParentExistsFilterQueryParam:
			parentExistsFilter, err = strconv.ParseBool(value[0])
			if err != nil {
				return nil, err
			}
		case PageOffsetQueryParam:
			pageOffset, err = strconv.Atoi(value[0])
			if err != nil {
				return nil, err
			}
		case PageLimitQueryParam:
			pageLimit, err = strconv.Atoi(value[0])
			if err != nil {
				return nil, err
			}
		}
	}
	return &WorkItemsResourceListContext{
		APIContexter:        r.Context,
		SpaceID:             spaceID,
		Filter:              &filter,
		ExpressionFilter:    &expressionFilter,
		AssigneeFilter:      &assigneeFilter,
		IterationFilter:     &iterationFilter,
		WorkitemTypeFilter:  &workItemTypeFilter,
		AreaFilter:          &areaFilter,
		WorkitemStateFilter: &workItemStateFilter,
		ParentExistsFilter:  &parentExistsFilter,
		PageOffset:          &pageOffset,
		PageLimit:           &pageLimit,
	}, nil
}

func (res WorkItemsResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	ctx := r.Context
	listCtx, err := NewWorkItemsResourceListContext(r)
	if err != nil {
		return Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusBadRequest)
	}
	if err != nil {
		return Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusBadRequest)
	}
	//var additionalQuery []string
	*listCtx.ExpressionFilter = fmt.Sprintf(`{"%s":[{"space": "%s"}, %s]}`, search.Q_AND, listCtx.SpaceID, *listCtx.ExpressionFilter)
	exp, err := query.Parse(listCtx.Filter)
	/*
		if err != nil {
			return Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusBadRequest)
			//abortWithError(ctx, errors.NewBadParameterError("filter", err))
		}
		if listCtx.ExpressionFilter != nil {
			q := *listCtx.ExpressionFilter
			*listCtx.ExpressionFilter = fmt.Sprintf(`{"%s":[{"space": "%s"}, %s]}`, search.Q_AND, listCtx.SpaceID, q)
			additionalQuery = append(additionalQuery, "filter[expression]="+*listCtx.ExpressionFilter)
				q := *listCtx.ExpressionFilter
				// Better approach would be to convert string to Query instance itself.
				// Then add new AND clause with spaceID as another child of input query
				// Then convert new Query object into simple string
				queryWithSpaceID := fmt.Sprintf(`?filter[expression]={"%s":[{"space": "%s" }, %s]}`, search.Q_AND, listCtx.SpaceID, q)
				searchURL := app.SearchHref() + queryWithSpaceID
				ctx.Header("Location", searchURL)
				ctx.Status(http.StatusTemporaryRedirect)
				return
			// return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusBadRequest)
		}
		if listCtx.AssigneeFilter != nil {
			if *listCtx.AssigneeFilter == none {
				exp = criteria.And(exp, criteria.IsNull("system.assignees"))
				additionalQuery = append(additionalQuery, "filter[assignee]=none")
			} else {
				exp = criteria.And(exp, criteria.Equals(criteria.Field("system.assignees"), criteria.Literal([]string{*listCtx.AssigneeFilter})))
				additionalQuery = append(additionalQuery, "filter[assignee]="+*listCtx.AssigneeFilter)
			}
		}
		if listCtx.IterationFilter != nil {
			exp = criteria.And(exp, criteria.Equals(criteria.Field(workitem.SystemIteration), criteria.Literal(string(*listCtx.IterationFilter))))
			additionalQuery = append(additionalQuery, "filter[iteration]="+*listCtx.IterationFilter)
			// Update filter by adding child iterations if any
			application.Transactional(res.db, func(tx application.Application) error {
				iterationUUID, errConversion := uuid.FromString(*listCtx.IterationFilter)
				if errConversion != nil {
					//ctx.AbortWithError(http.StatusBadRequest, errors.NewBadParameterError("iterationID", errConversion))
					return errConversion
				}
				childrens, err := tx.Iterations().LoadChildren(ctx, iterationUUID)
				if err != nil {
					//ctx.AbortWithError(http.StatusBadRequest, err)
					return err
				}
				for _, child := range childrens {
					childIDStr := child.ID.String()
					exp = criteria.Or(exp, criteria.Equals(criteria.Field(workitem.SystemIteration), criteria.Literal(childIDStr)))
					additionalQuery = append(additionalQuery, "filter[iteration]="+childIDStr)
				}
				//return Response{},
				return nil
			})
		}
		if listCtx.WorkitemTypeFilter != nil {
			exp = criteria.And(exp, criteria.Equals(criteria.Field("Type"), criteria.Literal([]uuid.UUID{*listCtx.WorkitemTypeFilter})))
			additionalQuery = append(additionalQuery, "filter[workitemtype]="+listCtx.WorkitemTypeFilter.String())
		}
		if listCtx.AreaFilter != nil {
			exp = criteria.And(exp, criteria.Equals(criteria.Field(workitem.SystemArea), criteria.Literal(string(*listCtx.AreaFilter))))
			additionalQuery = append(additionalQuery, "filter[area]="+*listCtx.AreaFilter)
		}
		if listCtx.WorkitemStateFilter != nil {
			exp = criteria.And(exp, criteria.Equals(criteria.Field(workitem.SystemState), criteria.Literal(string(*listCtx.WorkitemStateFilter))))
			additionalQuery = append(additionalQuery, "filter[workitemstate]="+*listCtx.WorkitemStateFilter)
		}
		if listCtx.ParentExistsFilter != nil {
			// no need to build expression: it is taken care in wi.List call
			// we need additionalQuery to make sticky filters in URL links
			additionalQuery = append(additionalQuery, "filter[parentexists]="+strconv.FormatBool(*listCtx.ParentExistsFilter))
		}
	*/
	workitems, _, err := res.db.WorkItems().List(ctx, listCtx.SpaceID, exp, listCtx.ParentExistsFilter, nil, nil)
	if err != nil {
		return Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusBadRequest)
	}
	return Response{Res: workitems}, nil
}
