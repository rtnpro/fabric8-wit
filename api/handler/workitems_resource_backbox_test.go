package handler_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/fabric8-services/fabric8-wit/account"
	"github.com/fabric8-services/fabric8-wit/api/model"
	"github.com/fabric8-services/fabric8-wit/gormtestsupport"
	"github.com/fabric8-services/fabric8-wit/workitem"
	"github.com/google/jsonapi"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type WorkItemsResourceTestSuite struct {
	gormtestsupport.GinkgoTestSuite
	clean func()
	repo  workitem.WorkItemTypeRepository
	ctx   context.Context
}

var _ = Describe("WorkItems", func() {

	// FIXME: need to create a test space for all those tests
	testSpaceID := "2e0698d8-753e-4cef-bb7c-f027634824a2"

	s := WorkItemsResourceTestSuite{GinkgoTestSuite: gormtestsupport.NewGinkgoTestSuite("../../config.yaml")}

	BeforeEach(func() {
		s.Setup()
	})

	AfterEach(func() {
		s.TearDown()
	})

	Describe("Test WorkItems", func() {
		Context("Create WorkItem", func() {

			It("Create WorkItem OK", func() {
				// given
				testIdentity := createOneRandomUserIdentity(context.Background(), s.DB)
				title := "A title"
				description := "A description"
				state := workitem.SystemStateNew
				wi := model.WorkItem{
					Title:       &title,
					Description: &description,
					State:       &state,
					Type: &model.WorkItemType{
						ID: "26787039-b68f-4e28-8814-c2f93be1ef4e",
					},
				}
				payload := bytes.NewBuffer(make([]byte, 0))
				err := jsonapi.MarshalPayload(payload, &wi)
				require.Nil(GinkgoT(), err)
				r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/spaces/%[1]s/workitems", testSpaceID), payload)
				r.Header.Set("Authorization", "Bearer "+makeTokenString("HS256", testIdentity.ID.String()))
				// when
				rr := Execute(s.GinkgoTestSuite, r)
				// then
				assert.Equal(GinkgoT(), http.StatusCreated, rr.Code)
				responseItem := model.WorkItem{}
				GinkgoT().Logf("Response body:\n%s", rr.Body.String())
				err = jsonapi.UnmarshalPayload(rr.Body, &responseItem)
				require.Nil(GinkgoT(), err)
				assert.NotNil(GinkgoT(), responseItem.ID)
				assert.Equal(GinkgoT(), "A description", *responseItem.Description)
			})

			It("Create WorkItem KO - missing JWT", func() {
				// given
				title := "A title"
				description := "A description"
				state := workitem.SystemStateNew
				wi := model.WorkItem{
					Title:       &title,
					Description: &description,
					State:       &state,
					Type: &model.WorkItemType{
						ID: "26787039-b68f-4e28-8814-c2f93be1ef4e",
					},
				}
				payload := bytes.NewBuffer(make([]byte, 0))
				err := jsonapi.MarshalPayload(payload, &wi)
				require.Nil(GinkgoT(), err)
				r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/spaces/%[1]s/workitems", testSpaceID), payload)
				// when
				rr := Execute(s.GinkgoTestSuite, r)
				// then
				assert.Equal(GinkgoT(), http.StatusUnauthorized, rr.Code)
			})

			It("Create WorkItem KO - invalid credentials", func() {
				// given
				title := "A title"
				description := "A description"
				state := workitem.SystemStateNew
				wi := model.WorkItem{
					Title:       &title,
					Description: &description,
					State:       &state,
					Type: &model.WorkItemType{
						ID: "26787039-b68f-4e28-8814-c2f93be1ef4e",
					},
				}
				payload := bytes.NewBuffer(make([]byte, 0))
				err := jsonapi.MarshalPayload(payload, &wi)
				require.Nil(GinkgoT(), err)
				r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/spaces/%[1]s/workitems", testSpaceID), payload)
				// generate/sign an auth token
				r.Header.Set("Authorization", "Bearer "+makeTokenString("HS256", "foo"))
				// when
				rr := Execute(s.GinkgoTestSuite, r)
				// then
				assert.Equal(GinkgoT(), http.StatusForbidden, rr.Code)
			})
		})

		Context("Update WorkItem", func() {

			var testIdentity *account.Identity
			var createdWI model.WorkItem
			var payload *bytes.Buffer

			BeforeEach(func() {
				GinkgoT().Log("creating a work item to test the updates...")
				// given
				testIdentity = createOneRandomUserIdentity(context.Background(), s.DB)
				title := "A title"
				description := "A description"
				state := workitem.SystemStateNew
				wi := model.WorkItem{
					Title:       &title,
					Description: &description,
					State:       &state,
					Type: &model.WorkItemType{
						ID: "26787039-b68f-4e28-8814-c2f93be1ef4e",
					},
				}
				payload := bytes.NewBuffer(make([]byte, 0))
				err := jsonapi.MarshalPayload(payload, &wi)
				require.Nil(GinkgoT(), err)
				r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/spaces/%[1]s/workitems", testSpaceID), payload)
				r.Header.Set("Authorization", "Bearer "+makeTokenString("HS256", testIdentity.ID.String()))
				// when
				rr := Execute(s.GinkgoTestSuite, r)
				// then
				require.Equal(GinkgoT(), http.StatusCreated, rr.Code)
				err = jsonapi.UnmarshalPayload(rr.Body, &createdWI)
				require.Nil(GinkgoT(), err)
			})

			It("Update WorkItem OK", func() {
				// given
				updatedTitle := "Updated title"
				createdWI.Title = &updatedTitle
				payload = bytes.NewBuffer(make([]byte, 0))
				err := jsonapi.MarshalPayload(payload, &createdWI)
				require.Nil(GinkgoT(), err)
				r, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/workitems/%[1]s", createdWI.ID), payload)
				r.Header.Set("Authorization", "Bearer "+makeTokenString("HS256", testIdentity.ID.String()))
				// when
				rr := Execute(s.GinkgoTestSuite, r)
				// then
				assert.Equal(GinkgoT(), http.StatusOK, rr.Code)
				responseItem := model.WorkItem{}
				GinkgoT().Logf("Response body:\n%s", rr.Body.String())
				err = jsonapi.UnmarshalPayload(rr.Body, &responseItem)
				require.Nil(GinkgoT(), err)
				assert.NotNil(GinkgoT(), responseItem.ID)
				assert.Equal(GinkgoT(), "Updated title", *responseItem.Title)
			})

			It("Update WorkItem KO - invalid credentials", func() {
				// given
				updatedTitle := "Updated title"
				createdWI.Title = &updatedTitle
				payload = bytes.NewBuffer(make([]byte, 0))
				err := jsonapi.MarshalPayload(payload, &createdWI)
				require.Nil(GinkgoT(), err)
				r, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/workitems/%[1]s", createdWI.ID), payload)
				// generate an invalid auth token
				r.Header.Set("Authorization", "Bearer "+makeTokenString("HS256", "foo"))
				// when
				rr := Execute(s.GinkgoTestSuite, r)
				// then
				assert.Equal(GinkgoT(), http.StatusForbidden, rr.Code)
			})
		})

		Context("List WorkItems", func() {
			It("List WorkItems OK", func() {
				// given
				r, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/spaces/%[1]s/workitems", testSpaceID), nil)
				// when
				rr := Execute(s.GinkgoTestSuite, r)
				// then
				assert.Equal(GinkgoT(), http.StatusOK, rr.Code)
			})

			It("Show WorkItem", func() {
				// given
				r, _ := http.NewRequest(http.MethodGet, "/api/workitems/c870914b-7942-4b87-8271-3afda49004e0", nil)
				// when
				rr := Execute(s.GinkgoTestSuite, r)
				// then
				assert.Equal(GinkgoT(), http.StatusOK, rr.Code)
			})
		})

	})

})
