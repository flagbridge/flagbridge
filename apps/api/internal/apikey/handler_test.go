//go:build integration

package apikey_test

import (
	"net/http"
	"testing"

	"github.com/flagbridge/flagbridge/apps/api/internal/testutil"
)

func TestAPIKeyHandler_Create(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)

	srv := testutil.CreateTestServer(t, pool)
	defer srv.Close()

	token := testutil.MustLogin(t, srv)
	project := testutil.MustCreateProject(t, srv, token, "API Key Project", "apikey-project")
	projectID := project["id"].(string)

	t.Run("creates eval API key", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/api-keys", token, map[string]any{
			"name":       "My Eval Key",
			"scope":      "eval",
			"project_id": projectID,
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
				ID        string `json:"id"`
				Name      string `json:"name"`
				Scope     string `json:"scope"`
				KeyPrefix string `json:"key_prefix"`
				Key       string `json:"key"`
				ProjectID string `json:"project_id"`
			} `json:"data"`
		}
		testutil.DecodeBody(t, resp, &out)

		if out.Data.ID == "" {
			t.Error("expected non-empty id")
		}
		if out.Data.Name != "My Eval Key" {
			t.Errorf("name: got %q, want %q", out.Data.Name, "My Eval Key")
		}
		if out.Data.Scope != "eval" {
			t.Errorf("scope: got %q, want %q", out.Data.Scope, "eval")
		}
		if out.Data.ProjectID != projectID {
			t.Errorf("project_id: got %q, want %q", out.Data.ProjectID, projectID)
		}
		// The full key is only returned at creation time.
		if out.Data.Key == "" {
			t.Error("expected full key in creation response")
		}
		// Key must match the expected prefix format.
		if len(out.Data.Key) < 6 || out.Data.Key[:6] != "fb_sk_" {
			t.Errorf("key has wrong prefix: %q", out.Data.Key)
		}
	})

	t.Run("creates test scope API key", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/api-keys", token, map[string]any{
			"name":       "My Test Key",
			"scope":      "test",
			"project_id": projectID,
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
				Scope string `json:"scope"`
			} `json:"data"`
		}
		testutil.DecodeBody(t, resp, &out)

		if out.Data.Scope != "test" {
			t.Errorf("scope: got %q, want %q", out.Data.Scope, "test")
		}
	})

	t.Run("rejects invalid scope", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/api-keys", token, map[string]any{
			"name":       "Bad Scope Key",
			"scope":      "superadmin",
			"project_id": projectID,
		})

		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400 for invalid scope, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects request with missing name", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/api-keys", token, map[string]any{
			"scope":      "eval",
			"project_id": projectID,
		})

		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400 for missing name, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects request with missing project_id", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/api-keys", token, map[string]any{
			"name":  "No Project Key",
			"scope": "eval",
		})

		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected 400 for missing project_id, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects unauthenticated request", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/api-keys", "", map[string]any{
			"name":       "Anon Key",
			"scope":      "eval",
			"project_id": projectID,
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

func TestAPIKeyHandler_List(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)

	srv := testutil.CreateTestServer(t, pool)
	defer srv.Close()

	token := testutil.MustLogin(t, srv)
	project := testutil.MustCreateProject(t, srv, token, "List Keys Project", "list-keys-project")
	projectID := project["id"].(string)

	t.Run("returns empty list when no keys exist", func(t *testing.T) {
		req := testutil.AuthedRequest(t, http.MethodGet, srv.URL+"/v1/api-keys", token, nil)
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
			t.Errorf("expected 0 keys, got %d", len(out.Data))
		}
	})

	t.Run("lists created keys", func(t *testing.T) {
		// Create two keys.
		for _, name := range []string{"Key Alpha", "Key Beta"} {
			createReq := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/api-keys", token, map[string]any{
				"name":       name,
				"scope":      "eval",
				"project_id": projectID,
			})
			createResp, err := srv.Client().Do(createReq)
			if err != nil {
				t.Fatalf("create key request failed: %v", err)
			}
			createResp.Body.Close()
			if createResp.StatusCode != http.StatusCreated {
				t.Fatalf("create key returned %d for %q", createResp.StatusCode, name)
			}
		}

		req := testutil.AuthedRequest(t, http.MethodGet, srv.URL+"/v1/api-keys", token, nil)
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
			t.Errorf("expected 2 keys, got %d", len(out.Data))
		}

		// Verify the full key is NOT included in list responses (security).
		for _, k := range out.Data {
			if _, exists := k["key"]; exists {
				t.Error("list response must not include the full API key value")
			}
		}
	})
}

func TestAPIKeyHandler_Delete(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, pool)

	srv := testutil.CreateTestServer(t, pool)
	defer srv.Close()

	token := testutil.MustLogin(t, srv)
	project := testutil.MustCreateProject(t, srv, token, "Delete Keys Project", "delete-keys-project")
	projectID := project["id"].(string)

	t.Run("deletes an existing API key", func(t *testing.T) {
		// Create a key.
		createReq := testutil.AuthedRequest(t, http.MethodPost, srv.URL+"/v1/api-keys", token, map[string]any{
			"name":       "Temporary Key",
			"scope":      "eval",
			"project_id": projectID,
		})
		createResp, err := srv.Client().Do(createReq)
		if err != nil {
			t.Fatalf("create key request failed: %v", err)
		}
		defer createResp.Body.Close()

		var createOut struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		testutil.DecodeBody(t, createResp, &createOut)
		keyID := createOut.Data.ID

		// Delete it.
		delReq := testutil.AuthedRequest(t, http.MethodDelete, srv.URL+"/v1/api-keys/"+keyID, token, nil)
		delResp, err := srv.Client().Do(delReq)
		if err != nil {
			t.Fatalf("delete request failed: %v", err)
		}
		defer delResp.Body.Close()

		if delResp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", delResp.StatusCode)
		}

		// Verify it no longer appears in list.
		listReq := testutil.AuthedRequest(t, http.MethodGet, srv.URL+"/v1/api-keys", token, nil)
		listResp, err := srv.Client().Do(listReq)
		if err != nil {
			t.Fatalf("list request failed: %v", err)
		}
		defer listResp.Body.Close()

		var listOut struct {
			Data []any `json:"data"`
		}
		testutil.DecodeBody(t, listResp, &listOut)

		if len(listOut.Data) != 0 {
			t.Errorf("expected 0 keys after deletion, got %d", len(listOut.Data))
		}
	})

	t.Run("returns 404 when deleting non-existent key", func(t *testing.T) {
		fakeID := "00000000-0000-0000-0000-000000000099"
		req := testutil.AuthedRequest(t, http.MethodDelete, srv.URL+"/v1/api-keys/"+fakeID, token, nil)
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
