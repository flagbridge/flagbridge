//go:build integration

package flag_test

import (
	"net/http"
	"testing"

	"github.com/flagbridge/flagbridge/internal/testutil"
)

func TestFlagHandler_Create(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)

	srv := testutil.CreateTestServer(t, pool)
	defer srv.Close()

	token := testutil.MustLogin(t, srv)
	testutil.MustCreateProject(t, srv, token, "Flag Project", "flag-project")

	t.Run("creates a boolean flag", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/projects/flag-project/flags", token, map[string]any{
			"key":           "my-flag",
			"name":          "My Flag",
			"type":          "boolean",
			"default_value": false,
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
				ID           string `json:"id"`
				Key          string `json:"key"`
				Name         string `json:"name"`
				Type         string `json:"type"`
				DefaultValue any    `json:"default_value"`
			} `json:"data"`
		}
		testutil.DecodeBody(t, resp, &out)

		if out.Data.ID == "" {
			t.Error("expected non-empty id")
		}
		if out.Data.Key != "my-flag" {
			t.Errorf("key: got %q, want %q", out.Data.Key, "my-flag")
		}
		if out.Data.Type != "boolean" {
			t.Errorf("type: got %q, want %q", out.Data.Type, "boolean")
		}
	})

	t.Run("creates a string flag", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/projects/flag-project/flags", token, map[string]any{
			"key":           "hero-variant",
			"name":          "Hero Variant",
			"type":          "string",
			"default_value": "control",
		})

		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects invalid key format", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/projects/flag-project/flags", token, map[string]any{
			"key":  "INVALID KEY!!",
			"name": "Bad Flag",
		})

		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400 for invalid key, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects duplicate key in same project", func(t *testing.T) {
		testutil.MustCreateFlag(t, srv, token, "flag-project", "unique-flag", "Unique Flag")

		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/projects/flag-project/flags", token, map[string]any{
			"key":  "unique-flag",
			"name": "Unique Flag Again",
			"type": "boolean",
		})

		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400 for duplicate key, got %d", resp.StatusCode)
		}
	})

	t.Run("returns 404 for unknown project", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/projects/ghost-project/flags", token, map[string]any{
			"key":  "some-flag",
			"name": "Some Flag",
		})

		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404 for unknown project, got %d", resp.StatusCode)
		}
	})
}

func TestFlagHandler_List(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)

	srv := testutil.CreateTestServer(t, pool)
	defer srv.Close()

	token := testutil.MustLogin(t, srv)
	testutil.MustCreateProject(t, srv, token, "List Flags Project", "list-flags-project")

	t.Run("returns empty list when no flags exist", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodGet, srv.URL+"/v1/projects/list-flags-project/flags", token, nil)
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

		if len(out.Data) != 0 {
			t.Errorf("expected 0 flags, got %d", len(out.Data))
		}
	})

	t.Run("returns flags for the project only", func(t *testing.T) {
		testutil.MustCreateFlag(t, srv, token, "list-flags-project", "flag-one", "Flag One")
		testutil.MustCreateFlag(t, srv, token, "list-flags-project", "flag-two", "Flag Two")

		// Create another project with its own flag to verify isolation.
		testutil.MustCreateProject(t, srv, token, "Other Project", "other-project")
		testutil.MustCreateFlag(t, srv, token, "other-project", "other-flag", "Other Flag")

		req := testutil.AuthedRequest(t, http.MethodGet, srv.URL+"/v1/projects/list-flags-project/flags", token, nil)
		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		var out struct {
			Data []map[string]any `json:"data"`
		}
		testutil.DecodeBody(t, resp, &out)

		if len(out.Data) != 2 {
			t.Errorf("expected 2 flags for project, got %d", len(out.Data))
		}
	})
}

func TestFlagHandler_Get(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)

	srv := testutil.CreateTestServer(t, pool)
	defer srv.Close()

	token := testutil.MustLogin(t, srv)
	testutil.MustCreateProject(t, srv, token, "Get Flag Project", "get-flag-project")

	t.Run("returns flag by key", func(t *testing.T) {
		testutil.MustCreateFlag(t, srv, token, "get-flag-project", "target-flag", "Target Flag")

		req := testutil.AuthedRequest(t, http.MethodGet, srv.URL+"/v1/projects/get-flag-project/flags/target-flag", token, nil)
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
				Key  string `json:"key"`
				Name string `json:"name"`
			} `json:"data"`
		}
		testutil.DecodeBody(t, resp, &out)

		if out.Data.Key != "target-flag" {
			t.Errorf("key: got %q, want %q", out.Data.Key, "target-flag")
		}
	})

	t.Run("returns 404 for unknown flag key", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodGet, srv.URL+"/v1/projects/get-flag-project/flags/ghost-flag", token, nil)
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

func TestFlagHandler_Delete(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)

	srv := testutil.CreateTestServer(t, pool)
	defer srv.Close()

	token := testutil.MustLogin(t, srv)
	testutil.MustCreateProject(t, srv, token, "Delete Flag Project", "delete-flag-project")

	t.Run("deletes existing flag", func(t *testing.T) {
		testutil.MustCreateFlag(t, srv, token, "delete-flag-project", "doomed-flag", "Doomed Flag")

		req := testutil.AuthedRequest(t, http.MethodDelete, srv.URL+"/v1/projects/delete-flag-project/flags/doomed-flag", token, nil)
		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}

		// Verify it's gone.
		getReq := testutil.AuthedRequest(t, http.MethodGet, srv.URL+"/v1/projects/delete-flag-project/flags/doomed-flag", token, nil)
		getResp, err := srv.Client().Do(getReq)
		if err != nil {
			t.Fatalf("get after delete failed: %v", err)
		}
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404 after delete, got %d", getResp.StatusCode)
		}
	})

	t.Run("returns 404 when deleting non-existent flag", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodDelete, srv.URL+"/v1/projects/delete-flag-project/flags/does-not-exist", token, nil)
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

func TestFlagHandler_SetState(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)

	srv := testutil.CreateTestServer(t, pool)
	defer srv.Close()

	token := testutil.MustLogin(t, srv)
	testutil.MustCreateProject(t, srv, token, "State Project", "state-project")

	// Create an environment for state tests.
	envReq := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/projects/state-project/environments", token, map[string]any{
		"name":       "Production",
		"slug":       "production",
		"sort_order": 0,
	})
	envResp, err := srv.Client().Do(envReq)
	if err != nil {
		t.Fatalf("create environment request failed: %v", err)
	}
	envResp.Body.Close()
	if envResp.StatusCode != http.StatusCreated {
		t.Fatalf("create environment returned %d", envResp.StatusCode)
	}

	testutil.MustCreateFlag(t, srv, token, "state-project", "kill-switch", "Kill Switch")

	t.Run("enables a flag for an environment", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPut,
			srv.URL+"/v1/projects/state-project/flags/kill-switch/states/production",
			token, map[string]any{"enabled": true},
		)

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
				Enabled bool `json:"enabled"`
			} `json:"data"`
		}
		testutil.DecodeBody(t, resp, &out)

		if !out.Data.Enabled {
			t.Error("expected flag to be enabled")
		}
	})

	t.Run("disables a flag for an environment", func(t *testing.T) {
		// First enable it.
		enableReq := testutil.AuthedRequest(t, http.MethodPut,
			srv.URL+"/v1/projects/state-project/flags/kill-switch/states/production",
			token, map[string]any{"enabled": true},
		)
		enableResp, _ := srv.Client().Do(enableReq)
		enableResp.Body.Close()

		// Then disable it.
		req := testutil.AuthedRequest(t, http.MethodPut,
			srv.URL+"/v1/projects/state-project/flags/kill-switch/states/production",
			token, map[string]any{"enabled": false},
		)
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
				Enabled bool `json:"enabled"`
			} `json:"data"`
		}
		testutil.DecodeBody(t, resp, &out)

		if out.Data.Enabled {
			t.Error("expected flag to be disabled")
		}
	})

	t.Run("returns 404 for unknown environment", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPut,
			srv.URL+"/v1/projects/state-project/flags/kill-switch/states/ghost-env",
			token, map[string]any{"enabled": true},
		)
		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404 for unknown env, got %d", resp.StatusCode)
		}
	})
}
