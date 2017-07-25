package design

import (
	d "github.com/goadesign/goa/design"
	a "github.com/goadesign/goa/design/apidsl"
)

var workItemTypeGroupSigle = JSONSingle(
	"workItemTypeGroupSigle",
	`Group of the work-item-types`,
	workItemTypeGroupData,
	workItemTypeGroupLinks,
)

var workItemTypeGroups = JSONList(
	"workItemTypeGroups",
	"...",
	workItemTypeGroupData,
	nil,
	nil,
)

var workItemTypeGroupLinks = a.Type("workItemTypeGroupLinks", func() {
	a.Attribute("self", d.String, func() {
		a.Example("http://api.openshift.io/api/spacetemplates/2d98c73d-6969-4ea6-958a-812c832b6c18/workitemtypegroups")
	})
	a.Required("self")
})

var workItemTypeGroup = a.Type("WorkItemTypeGroup", func() {
	a.Attribute("level", d.Integer, "denotes the hierarchical rank within the group")
	a.Attribute("sublevel", d.Integer, "denotes the hierarchical rank within the node")
	a.Attribute("group", d.String, "Name of the group this node belongs to")
	a.Attribute("name", d.String)
	a.Attribute("wit_collection", a.ArrayOf(d.UUID), "Slice of UUIDs of work item type")
	a.Required("level", "sublevel", "group", "name", "wit_collection")
})

var workItemTypeGroupAttributes = a.Type("WorkItemTypeGroupAttributes", func() {
	a.Attribute("hierarchy", a.ArrayOf(workItemTypeGroup))
	a.Required("hierarchy")
})

// workItemTypeGroupData is the JSONAPI store for the data of a work item link type.
var workItemTypeGroupData = a.Type("WorkItemTypeGroupData", func() {
	a.Description(`JSONAPI store for the data of a work item link type.
See also http://jsonapi.org/format/#document-resource-object`)
	a.Attribute("attributes", workItemTypeGroupAttributes)
	a.Attribute("included", a.ArrayOf(workItemTypeData), "An array of work item types")
	a.Required("attributes", "included")
})

var _ = a.Resource("work_item_type_group", func() {
	a.BasePath("/workitemtypegroups")
	a.Parent("space_template")

	a.Action("list", func() {
		a.Routing(
			a.GET(""),
		)
		a.Description("List of work item type groups")
		a.Response(d.OK, workItemTypeGroupSigle)
		a.Response(d.NotModified)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})
})