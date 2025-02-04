package axiom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// HINT(lukasmalkmus): Most of the tests below just test against the "api"
// endpoint. However, the "ingest" and "personal" implementation is the same as
// "api" one: Under the hood, they both use the TokenService. The integration
// tests make sure this implementation works against both endpoints.

func TestTokensService_List(t *testing.T) {
	exp := []*Token{
		{
			ID:   "08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e",
			Name: "Test",
			Scopes: []string{
				"*",
			},
		},
	}

	hf := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		_, err := fmt.Fprint(w, `[
			{
				"id": "08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e",
				"name": "Test",
				"scopes": [
            		"*"
        		]
			}
		]`)
		assert.NoError(t, err)
	}

	client, teardown := setup(t, "/api/v1/tokens/api", hf)
	defer teardown()

	res, err := client.Tokens.API.List(context.Background())
	require.NoError(t, err)

	assert.Equal(t, exp, res)
}

func TestTokensService_Get(t *testing.T) {
	exp := &Token{
		ID:   "08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e",
		Name: "Test",
		Scopes: []string{
			"*",
		},
	}

	hf := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		_, err := fmt.Fprint(w, `{
			"id": "08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e",
			"name": "Test",
			"scopes": [
				"*"
			]
		}`)
		assert.NoError(t, err)
	}

	client, teardown := setup(t, "/api/v1/tokens/api/08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e", hf)
	defer teardown()

	res, err := client.Tokens.API.Get(context.Background(), "08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e")
	require.NoError(t, err)

	assert.Equal(t, exp, res)
}

func TestTokensService_View(t *testing.T) {
	exp := &RawToken{
		Token: "ae51e8d9-5fa2-4957-9847-3c1ccfa5ffe9",
		Scopes: []string{
			"*",
		},
	}

	hf := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		_, err := fmt.Fprint(w, `{
			"token": "ae51e8d9-5fa2-4957-9847-3c1ccfa5ffe9",
			"scopes": [
				"*"
			]
		}`)
		assert.NoError(t, err)
	}

	client, teardown := setup(t, "/api/v1/tokens/api/08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e/token", hf)
	defer teardown()

	res, err := client.Tokens.API.View(context.Background(), "08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e")
	require.NoError(t, err)

	assert.Equal(t, exp, res)
}

func TestTokensService_Create(t *testing.T) {
	exp := &Token{
		ID:          "08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e",
		Name:        "Test",
		Description: "A test token",
		Scopes: []string{
			"*",
		},
	}

	hf := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		_, err := fmt.Fprint(w, `{
			"id": "08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e",
			"name": "Test",
			"description": "A test token",
			"scopes": [
				"*"
			]
		}`)
		assert.NoError(t, err)
	}

	client, teardown := setup(t, "/api/v1/tokens/api", hf)
	defer teardown()

	res, err := client.Tokens.API.Create(context.Background(), TokenCreateUpdateRequest{
		Name:        "Test",
		Description: "A test token",
	})
	require.NoError(t, err)

	assert.Equal(t, exp, res)
}

func TestTokensService_Update(t *testing.T) {
	exp := &Token{
		ID:          "08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e",
		Name:        "Test",
		Description: "A very good test token",
		Scopes: []string{
			"*",
		},
	}

	hf := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		_, err := fmt.Fprint(w, `{
			"id": "08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e",
			"name": "Test",
			"description": "A very good test token",
			"scopes": [
				"*"
			]
		}`)
		assert.NoError(t, err)
	}

	client, teardown := setup(t, "/api/v1/tokens/api/08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e", hf)
	defer teardown()

	res, err := client.Tokens.API.Update(context.Background(), "08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e", TokenCreateUpdateRequest{
		Name:        "Michael Doe",
		Description: "A very good test token",
	})
	require.NoError(t, err)

	assert.Equal(t, exp, res)
}

func TestTokensService_Delete(t *testing.T) {
	hf := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	}

	client, teardown := setup(t, "/api/v1/tokens/api/08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e", hf)
	defer teardown()

	err := client.Tokens.API.Delete(context.Background(), "08fceb797a467c3c23151f3584c31cfaea962e3ca306e3af69c2dab28e8c2e6e")
	require.NoError(t, err)
}

func TestPermission_Marshal(t *testing.T) {
	exp := `{
		"permission": "CanIngest"
	}`

	b, err := json.Marshal(struct {
		Permission Permission `json:"permission"`
	}{
		Permission: CanIngest,
	})
	require.NoError(t, err)
	require.NotEmpty(t, b)

	assert.JSONEq(t, exp, string(b))
}

func TestPermission_Unmarshal(t *testing.T) {
	var act struct {
		Permission Permission `json:"permission"`
	}
	err := json.Unmarshal([]byte(`{ "permission": "CanIngest" }`), &act)
	require.NoError(t, err)

	assert.Equal(t, CanIngest, act.Permission)
}

func TestPermission_String(t *testing.T) {
	// Check outer bounds.
	assert.Empty(t, Permission(0).String())
	assert.Empty(t, emptyPermission.String())
	assert.Equal(t, emptyPermission, Permission(0))
	assert.Contains(t, (CanQuery + 1).String(), "Permission(")

	for c := CanIngest; c <= CanQuery; c++ {
		s := c.String()
		assert.NotEmpty(t, s)
		assert.NotContains(t, s, "Permission(")
	}
}

func TestPermissionFromString(t *testing.T) {
	for permission := CanIngest; permission <= CanQuery; permission++ {
		s := permission.String()

		parsedPermission, err := permissionFromString(s)
		assert.NoError(t, err)

		assert.NotEmpty(t, s)
		assert.Equal(t, permission, parsedPermission)
	}
}
