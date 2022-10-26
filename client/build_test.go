package client

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/cloudfoundry-community/go-cfclient/test"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestBuilds(t *testing.T) {
	g := test.NewObjectJSONGenerator(2)
	build := g.Build()
	build2 := g.Build()
	build3 := g.Build()
	build4 := g.Build()

	tests := []RouteTest{
		{
			Description: "Create build",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/builds",
				Output:   []string{build},
				Status:   http.StatusCreated,
				PostForm: `{"metadata":{"labels":{"foo":"bar"},"annotations":null},"package":{"guid":"993386e8-5f68-403c-b372-d4aba7c71dbc"}}`},
			Expected: build,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewBuildCreate("993386e8-5f68-403c-b372-d4aba7c71dbc")
				r.Metadata = &resource.Metadata{
					Labels: map[string]string{
						"foo": "bar",
					},
				}
				return c.Builds.Create(r)
			},
		},
		{
			Description: "Get build",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/builds/be9db090-ad79-41c1-9a01-6200d896f20f",
				Output:   []string{build},
				Status:   http.StatusOK,
			},
			Expected: build,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Builds.Get("be9db090-ad79-41c1-9a01-6200d896f20f")
			},
		},
		{
			Description: "Delete build",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/builds/be9db090-ad79-41c1-9a01-6200d896f20f",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Builds.Delete("be9db090-ad79-41c1-9a01-6200d896f20f")
			},
		},
		{
			Description: "Update build",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/builds/be9db090-ad79-41c1-9a01-6200d896f20f",
				Output:   []string{build},
				PostForm: `{"metadata":{"labels":{"env":"dev"},"annotations":{"foo": "bar"}}}`,
				Status:   http.StatusOK,
			},
			Expected: build,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewBuildUpdate()
				r.Metadata.Annotations["foo"] = "bar"
				r.Metadata.Labels["env"] = "dev"
				return c.Builds.Update("be9db090-ad79-41c1-9a01-6200d896f20f", r)
			},
		},
		{
			Description: "List first page of builds",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/builds",
				Output:   g.Paged("builds", []string{build}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(build),
			Action: func(c *Client, t *testing.T) (any, error) {
				builds, _, err := c.Builds.List(NewBuildListOptions())
				return builds, err
			},
		},
		{
			Description: "List all builds",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/builds",
				Output:   g.Paged("apps", []string{build, build2}, []string{build3, build4}),
				Status:   http.StatusOK},
			Expected: g.Array(build, build2, build3, build4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Builds.ListAll()
			},
		},
		{
			Description: "List first page of builds for app",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/builds",
				Output:   g.Paged("builds", []string{build}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(build),
			Action: func(c *Client, t *testing.T) (any, error) {
				opts := NewBuildAppListOptions()
				builds, _, err := c.Builds.ListForApp("1cb006ee-fb05-47e1-b541-c34179ddc446", opts)
				return builds, err
			},
		},
	}
	for _, tt := range tests {
		func() {
			setup(tt.Route, t)
			defer teardown()
			details := fmt.Sprintf("%s %s", tt.Route.Method, tt.Route.Endpoint)
			if tt.Description != "" {
				details = tt.Description + ": " + details
			}

			c, _ := NewTokenConfig(server.URL, "foobar")
			cl, err := New(c)
			require.NoError(t, err, details)

			obj, err := tt.Action(cl, t)
			require.NoError(t, err, details)
			if tt.Expected != "" {
				actual, err := json.Marshal(obj)
				require.NoError(t, err, details)
				require.JSONEq(t, tt.Expected, string(actual), details)
			}
		}()
	}
}
