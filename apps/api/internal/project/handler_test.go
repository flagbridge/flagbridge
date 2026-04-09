//go:build integration

package project_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/flagbridge/flagbridge/apps/api/internal/testutil"
)

func TestProjectHandler_Create(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)

	srv := testutil.CreateTestServer(t, pool)
	defer srv.Close()

	token := testutil.MustLogin(t, srv)

	t.Run("creates project with valid fields", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/projects", token, map[string]string{
			"name":        "My Project",
			"slug":        "my-project",
			"description": "A test project",
		})

		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201, got %d", resp.StatusCode)
		}

		var out struct {
			Data struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				Slug        string `json:"slug"`
				Description string `json:"description"`
				CreatedBy   string `json:"created_by"`
			} `json:"data"`
		}
		testutil.DecodeBody(t, resp, &out)

		if out.Data.ID == "" {
			t.Error("expected non-empty id")
		}
		if out.Data.Name != "My Project" {
			t.Errorf("name: got %q, want %q", out.Data.Name, "My Project")
		}
		if out.Data.Slug != "my-project" {
			t.Errorf("slug: got %q, want %q", out.Data.Slug, "my-project")
		}
		if out.Data.Description != "A test project" {
			t.Errorf("description: got %q, want %q", out.Data.Description, "A test project")
		}
		if out.Data.CreatedBy != testutil.SeedUserID {
			t.Errorf("created_by: got %q, want %q", out.Data.CreatedBy, testutil.SeedUserID)
		}
	})

	t.Run("rejects request with missing name", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/projects", token, map[string]string{
			"slug": "no-name",
		})

		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects invalid slug format", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/projects", token, map[string]string{
			"name": "Bad Slug",
			"slug": "Bad_Slug!!",
		})

		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects duplicate slug", func(t *testing.T) {
		// First creation should succeed.
		testutil.MustCreateProject(t, srv, token, "Dupe Project", "dupe-slug")

		// Second creation with the same slug must fail.
		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/projects", token, map[string]string{
			"name": "Dupe Project 2",
			"slug": "dupe-slug",
		})

		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400 for duplicate slug, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects unauthenticated request", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/projects", "", map[string]string{
			"name": "Auth Project",
			"slug": "auth-project",
		})

		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
	})
}

func TestProjectHandler_List(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)

	srv := testutil.CreateTestServer(t, pool)
	defer srv.Close()

	token := testutil.MustLogin(t, srv)

	t.Run("returns empty list when no projects exist", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodGet, srv.URL+"/v1/projects", token, nil)
		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}

		var out struct {
			Data []any `json:"data"`
		}
		testutil.DecodeBody(t, resp, &out)

		// nil and empty slice both mean "no projects"
		if len(out.Data) != 0 {
			t.Errorf("expected 0 projects, got %d", len(out.Data))
		}
	})

	t.Run("returns all created projects", func(t *testing.T) {
		testutil.MustCreateProject(t, srv, token, "Project Alpha", "project-alpha")
		testutil.MustCreateProject(t, srv, token, "Project Beta", "project-beta")

		req := testutil.AuthedRequest(t, http.MethodGet, srv.URL+"/v1/projects", token, nil)
		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}

		var out struct {
			Data []map[string]any `json:"data"`
		}
		testutil.DecodeBody(t, resp, &out)

		if len(out.Data) != 2 {
			t.Errorf("expected 2 projects, got %d", len(out.Data))
		}
	})
}

func TestProjectHandler_GetBySlug(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)

	srv := testutil.CreateTestServer(t, pool)
	defer srv.Close()

	token := testutil.MustLogin(t, srv)

	t.Run("returns project by slug", func(t *testing.T) {
		testutil.MustCreateProject(t, srv, token, "Lookup Project", "lookup-project")

		req := testutil.AuthedRequest(t, http.MethodGet, srv.URL+"/v1/projects/lookup-project", token, nil)
		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}

		var out struct {
			Data struct {
				Slug string `json:"slug"`
				Name string `json:"name"`
			} `json:"data"`
		}
		testutil.DecodeBody(t, resp, &out)

		if out.Data.Slug != "lookup-project" {
			t.Errorf("slug: got %q, want %q", out.Data.Slug, "lookup-project")
		}
		if out.Data.Name != "Lookup Project" {
			t.Errorf("name: got %q, want %q", out.Data.Name, "Lookup Project")
		}
	})

	t.Run("returns 404 for non-existent slug", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodGet, srv.URL+"/v1/projects/does-not-exist", token, nil)
		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404, got %d", resp.StatusCode)
		}
	})
}

func TestProjectHandler_Update(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)

	srv := testutil.CreateTestServer(t, pool)
	defer srv.Close()

	token := testutil.MustLogin(t, srv)

	t.Run("updates project name", func(t *testing.T) {
		testutil.MustCreateProject(t, srv, token, "Original Name", "update-me")

		newName := "Updated Name"
		req := testutil.AuthedRequest(t, http.MethodPatch, srv.URL+"/v1/projects/update-me", token, map[string]*string{
			"name": &newName,
		})
		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}

		var out struct {
			Data struct {
				Name string `json:"name"`
				Slug string `json:"slug"`
			} `json:"data"`
		}
		testutil.DecodeBody(t, resp, &out)

		if out.Data.Name != "Updated Name" {
			t.Errorf("name: got %q, want %q", out.Data.Name, "Updated Name")
		}
		// Slug should remain unchanged.
		if out.Data.Slug != "update-me" {
			t.Errorf("slug changed unexpectedly: got %q", out.Data.Slug)
		}
	})

	t.Run("updates project description", func(t *testing.T) {
		testutil.MustCreateProject(t, srv, token, "Desc Project", "desc-project")

		newDesc := "Brand new description"
		req := testutil.AuthedRequest(t, http.MethodPatch, srv.URL+"/v1/projects/desc-project", token, map[string]*string{
			"description": &newDesc,
		})
		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}

		var out struct {
			Data struct {
				Description string `json:"description"`
			} `json:"data"`
		}
		testutil.DecodeBody(t, resp, &out)

		if out.Data.Description != "Brand new description" {
			t.Errorf("description: got %q, want %q", out.Data.Description, "Brand new description")
		}
	})

	t.Run("rejects update with no fields", func(t *testing.T) {
		testutil.MustCreateProject(t, srv, token, "No Field Project", "no-field-project")

		req := testutil.AuthedRequest(t, http.MethodPatch, srv.URL+"/v1/projects/no-field-project", token, map[string]any{})
		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400 for empty patch, got %d", resp.StatusCode)
		}
	})
}

func TestProjectHandler_Delete(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)

	srv := testutil.CreateTestServer(t, pool)
	defer srv.Close()

	token := testutil.MustLogin(t, srv)

	t.Run("deletes an existing project", func(t *testing.T) {
		testutil.MustCreateProject(t, srv, token, "Delete Me", "delete-me")

		req := testutil.AuthedRequest(t, http.MethodDelete, srv.URL+"/v1/projects/delete-me", token, nil)
		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", resp.StatusCode)
		}

		// Verify it's gone.
		getReq := testutil.AuthedRequest(t, http.MethodGet, srv.URL+"/v1/projects/delete-me", token, nil)
		getResp, err := srv.Client().Do(getReq)
		if err != nil {
			t.Fatalf("get after delete request failed: %v", err)
		}
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404 after delete, got %d", getResp.StatusCode)
		}
	})

	t.Run("returns 404 when deleting non-existent project", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodDelete, srv.URL+"/v1/projects/ghost-project", token, nil)
		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404, got %d", resp.StatusCode)
		}
	})

	t.Run("list is empty after all projects deleted", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			slug := fmt.Sprintf("batch-project-%d", i)
			testutil.MustCreateProject(t, srv, token, "Batch Project", slug)
			delReq := testutil.AuthedRequest(t, http.MethodDelete, srv.URL+"/v1/projects/"+slug, token, nil)
			delResp, _ := srv.Client().Do(delReq)
			delResp.Body.Close()
		}

		listReq := testutil.AuthedRequest(t, http.MethodGet, srv.URL+"/v1/projects", token, nil)
		listResp, err := srv.Client().Do(listReq)
		if err != nil {
			t.Fatalf("list request failed: %v", err)
		}
		defer listResp.Body.Close()

		var out struct {
			Data []any `json:"data"`
		}
		testutil.DecodeBody(t, listResp, &out)

		if len(out.Data) != 0 {
			t.Errorf("expected 0 projects after deletion, got %d", len(out.Data))
		}
	})
}
