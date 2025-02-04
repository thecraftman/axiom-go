package axiom

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	endpoint = "http://axiom.local"
	// apiToken is a placeholder API token.
	apiToken = "xaat-XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
	// personalToken is a placeholder personal token.
	personalToken = "xapt-XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX" //nolint:gosec // Chill, it's just testing.
	// orgID is a placeholder organization id.
	orgID = "awkward-identifier-c3po"
)

var tokenRe = regexp.MustCompile("xa(a|p|)t-[a-zA-z0-9]{8}-[a-zA-z0-9]{4}-[a-zA-z0-9]{4}-[a-zA-z0-9]{4}-[a-zA-z0-9]{12}")

// SetStrictDecoding is a special testing-only client option that - when set to
// 'true' - failes JSON response decoding if fields not present in the
// destination struct are encountered.
func SetStrictDecoding(b bool) Option {
	return func(c *Client) error {
		c.strictDecoding = b
		return nil
	}
}

func TestNewClient(t *testing.T) {
	defer os.Clearenv()

	tests := []struct {
		name        string
		environment map[string]string
		options     []Option
		err         error
	}{
		{
			name: "no environment no options",
			err:  ErrMissingAccessToken,
		},
		{
			name: "no environment accessToken option",
			options: []Option{
				SetAccessToken(personalToken),
			},
			err: ErrMissingOrganizationID,
		},
		{
			name: "no environment accessToken option with API token",
			options: []Option{
				SetAccessToken(apiToken),
			},
		},
		{
			name: "orgID environment no options",
			environment: map[string]string{
				"AXIOM_TOKEN": personalToken,
			},
			err: ErrMissingOrganizationID,
		},
		{
			name: "orgID environment no options with API token",
			environment: map[string]string{
				"AXIOM_TOKEN": apiToken,
			},
		},
		{
			name: "no environment accessToken and orgID option",
			options: []Option{
				SetAccessToken(personalToken),
				SetOrgID(orgID),
			},
		},
		{
			name: "accessToken and orgID environment no options",
			environment: map[string]string{
				"AXIOM_TOKEN":  personalToken,
				"AXIOM_ORG_ID": orgID,
			},
		},
		{
			name: "accessToken environment orgID option",
			environment: map[string]string{
				"AXIOM_TOKEN": personalToken,
			},
			options: []Option{
				SetOrgID(orgID),
			},
		},
		{
			name: "orgID environment accessToken option",
			environment: map[string]string{
				"AXIOM_ORG_ID": orgID,
			},
			options: []Option{
				SetAccessToken(personalToken),
			},
		},
		{
			name: "no environment url and accessToken option",
			options: []Option{
				SetURL(endpoint),
				SetAccessToken(personalToken),
			},
		},
		{
			name: "url and accessToken environment no options",
			environment: map[string]string{
				"AXIOM_URL":   endpoint,
				"AXIOM_TOKEN": personalToken,
			},
		},
		{
			name: "accessToken and orgID environment cloudUrl option",
			environment: map[string]string{
				"AXIOM_TOKEN":  personalToken,
				"AXIOM_ORG_ID": orgID,
			},
			options: []Option{
				SetURL(CloudURL),
			},
		},
		{
			name: "accessToken and orgID environment enhanced cloudUrl option",
			environment: map[string]string{
				"AXIOM_TOKEN":  personalToken,
				"AXIOM_ORG_ID": orgID,
			},
			options: []Option{
				SetURL(CloudURL + "/"),
			},
		},
		{
			name: "cloudUrl accessToken and orgID environment no options",
			environment: map[string]string{
				"AXIOM_URL":    CloudURL,
				"AXIOM_TOKEN":  personalToken,
				"AXIOM_ORG_ID": orgID,
			},
		},
		{
			name: "enhanced cloudUrl, accessToken and orgID environment no options",
			environment: map[string]string{
				"AXIOM_URL":    CloudURL + "/",
				"AXIOM_TOKEN":  personalToken,
				"AXIOM_ORG_ID": orgID,
			},
		},
		{
			name: "dev url, accessToken and orgID environment no options",
			environment: map[string]string{
				"AXIOM_URL":    "https://dev.axiom.co",
				"AXIOM_TOKEN":  personalToken,
				"AXIOM_ORG_ID": orgID,
			},
		},
		{
			name: "cloudUrl and accessToken environment orgID option",
			environment: map[string]string{
				"AXIOM_URL":   CloudURL,
				"AXIOM_TOKEN": personalToken,
			},
			options: []Option{
				SetOrgID(orgID),
			},
		},
		{
			name: "accessToken and orgID environment noEnv option",
			environment: map[string]string{
				"AXIOM_TOKEN":  personalToken,
				"AXIOM_ORG_ID": orgID,
			},
			options: []Option{
				SetNoEnv(),
			},
			err: ErrMissingAccessToken,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for k, v := range tt.environment {
				os.Setenv(k, v)
			}

			client, err := NewClient(tt.options...)
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
			} else {
				if assert.NoError(t, err) {
					assert.Regexp(t, tokenRe, client.accessToken)
					assert.NotEmpty(t, client.baseURL)
				}
			}
		})
	}
}

func TestNewClient_Valid(t *testing.T) {
	client := newClient(t)

	// Are endpoints/resources present?
	assert.NotNil(t, client.Dashboards)
	assert.NotNil(t, client.Datasets)
	assert.NotNil(t, client.Monitors)
	assert.NotNil(t, client.Notifiers)
	assert.NotNil(t, client.Organizations.Cloud)
	assert.NotNil(t, client.Organizations.Selfhost)
	assert.NotNil(t, client.StarredQueries)
	assert.NotNil(t, client.Teams)
	assert.NotNil(t, client.Tokens.API)
	assert.NotNil(t, client.Tokens.Personal)
	assert.NotNil(t, client.Users)
	assert.NotNil(t, client.Version)
	assert.NotNil(t, client.VirtualFields)

	// Is default configuration present?
	assert.Equal(t, endpoint, client.baseURL.String())
	assert.Equal(t, personalToken, client.accessToken)
	assert.Empty(t, client.orgID)
	assert.NotNil(t, client.httpClient)
	assert.NotEmpty(t, client.userAgent)
	assert.False(t, client.strictDecoding)
}

func TestClient_Options_SetAccessToken(t *testing.T) {
	client := newClient(t)

	exp := personalToken
	opt := SetAccessToken(exp)

	err := client.Options(opt)
	assert.NoError(t, err)

	assert.Equal(t, exp, client.accessToken)
}

func TestClient_Options_SetClient(t *testing.T) {
	client := newClient(t)

	exp := &http.Client{
		Timeout: time.Second,
	}
	opt := SetClient(exp)

	err := client.Options(opt)
	assert.NoError(t, err)

	assert.Equal(t, exp, client.httpClient)
}

func TestClient_Options_SetCloudConfig(t *testing.T) {
	client := newClient(t)

	opt := SetCloudConfig(personalToken, orgID)

	err := client.Options(opt)
	assert.NoError(t, err)

	assert.Equal(t, personalToken, client.accessToken)
	assert.Equal(t, orgID, client.orgID)
}

func TestClient_Options_SetOrgID(t *testing.T) {
	client := newClient(t)

	exp := orgID
	opt := SetOrgID(exp)

	err := client.Options(opt)
	assert.NoError(t, err)

	assert.Equal(t, exp, client.orgID)
}

func TestClient_Options_SetSelfhostConfig(t *testing.T) {
	client := newClient(t)

	opt := SetSelfhostConfig(endpoint, personalToken)

	err := client.Options(opt)
	assert.NoError(t, err)

	assert.Equal(t, endpoint, client.baseURL.String())
	assert.Equal(t, personalToken, client.accessToken)
}

func TestClient_Options_SetURL(t *testing.T) {
	client := newClient(t)

	exp := endpoint
	opt := SetURL(exp)

	err := client.Options(opt)
	assert.NoError(t, err)

	assert.Equal(t, exp, client.baseURL.String())
}

func TestClient_Options_SetUserAgent(t *testing.T) {
	client := newClient(t)

	exp := "axiom-go/1.0.0"
	opt := SetUserAgent(exp)

	err := client.Options(opt)
	assert.NoError(t, err)

	assert.Equal(t, exp, client.userAgent)
}

func TestClient_newRequest_BadURL(t *testing.T) {
	client := newClient(t)

	_, err := client.newRequest(context.Background(), http.MethodGet, ":", nil)
	assert.Error(t, err)

	if assert.IsType(t, err, new(url.Error)) {
		urlErr := err.(*url.Error)
		assert.Equal(t, urlErr.Op, "parse")
	}
}

// If a nil body is passed to NewRequest, make sure that nil is also passed to
// http.NewRequest. In most cases, passing an io.Reader that returns no content
// is fine, since there is no difference between an HTTP request body that is an
// empty string versus one that is not set at all. However in certain cases,
// intermediate systems may treat these differently resulting in subtle errors.
func TestClient_newRequest_EmptyBody(t *testing.T) {
	client := newClient(t)

	req, err := client.newRequest(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	assert.Empty(t, req.Body)
}

func TestClient_do(t *testing.T) {
	hf := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `{"A":"a"}`)
	}

	client, teardown := setup(t, "/", hf)
	defer teardown()

	type foo struct {
		A string
	}

	req, err := client.newRequest(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	var body foo
	_, err = client.do(req, &body)
	require.NoError(t, err)

	assert.Equal(t, foo{"a"}, body)
}

func TestClient_do_ioWriter(t *testing.T) {
	content := `{"A":"a"}`

	hf := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, content)
	}

	client, teardown := setup(t, "/", hf)
	defer teardown()

	req, err := client.newRequest(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	var buf bytes.Buffer
	_, err = client.do(req, &buf)
	require.NoError(t, err)

	assert.Equal(t, content, buf.String())
}

func TestClient_do_HTTPError(t *testing.T) {
	hf := func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusBadRequest

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)

		httpErr := Error{
			Message: "This is a Bad Request error",
		}

		err := json.NewEncoder(w).Encode(httpErr)
		assert.NoError(t, err)
	}

	client, teardown := setup(t, "/", hf)
	defer teardown()

	req, err := client.newRequest(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	resp, err := client.do(req, nil)
	require.NotNil(t, resp)
	require.Error(t, err)

	if assert.IsType(t, Error{}, err) {
		assert.EqualError(t, err, "API error 400: This is a Bad Request error")
	}
}

func TestClient_do_Unauthenticated(t *testing.T) {
	hf := func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusUnauthorized

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)

		httpErr := Error{
			Message: "You are not allowed here!",
		}

		err := json.NewEncoder(w).Encode(httpErr)
		assert.NoError(t, err)
	}

	client, teardown := setup(t, "/", hf)
	defer teardown()

	req, err := client.newRequest(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	_, err = client.do(req, nil)
	require.ErrorIs(t, err, ErrUnauthenticated)
}

func TestClient_do_UnprivilegedToken(t *testing.T) {
	client, teardown := setup(t, "/", nil)
	defer teardown()

	err := client.Options(SetAccessToken("xaat-123"))
	require.NoError(t, err)

	_, err = client.newRequest(context.Background(), http.MethodGet, "/", nil)
	require.ErrorIs(t, err, ErrUnprivilegedToken)
}

func TestClient_do_RedirectLoop(t *testing.T) {
	hf := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	client, teardown := setup(t, "/", hf)
	defer teardown()

	req, err := client.newRequest(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	_, err = client.do(req, nil)
	require.Error(t, err)

	assert.IsType(t, err, new(url.Error))
}

func TestClient_do_validOnlyAPITokenPaths(t *testing.T) {
	hf := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	tests := []string{
		"/api/v1/datasets/test/query",
		"/api/v1/datasets/_apl",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			client, teardown := setup(t, tt, hf)
			defer teardown()

			err := client.Options(SetAccessToken("xaat-123"))
			require.NoError(t, err)

			req, err := client.newRequest(context.Background(), http.MethodGet, tt, nil)
			require.Nil(t, err)

			_, err = client.do(req, nil)
			require.NoError(t, err)
		})
	}
}

func TestAPITokenPathRegex(t *testing.T) {
	tests := []struct {
		input string
		match bool
	}{
		{
			input: "/api/v1/datasets/test/ingest",
			match: true,
		},
		{
			input: "/api/v1/datasets/test/ingest?timestamp-format=unix",
			match: true,
		},
		{
			input: "/api/v1/datasets/test/query",
			match: true,
		},
		{
			input: "/api/v1/datasets/_apl",
			match: true,
		},
		{
			input: "/api/v1/datasets/test/query?nocache=true",
			match: true,
		},
		{
			input: "/api/v1/datasets/_apl?nocache=true",
			match: true,
		},
		{
			input: "/api/v1/datasets//query",
			match: false,
		},
		{
			input: "/api/v1/datasets/query",
			match: false,
		},
		{
			input: "/api/v1/datasets/test/elastic",
			match: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.match, validOnlyAPITokenPaths.MatchString(tt.input))
		})
	}
}

// setup sets up a test HTTP server along with a client that is configured to
// talk to that test server. Tests should pass a handler function which provides
// the response for the API method being tested.
func setup(t *testing.T, path string, handler http.HandlerFunc) (*Client, func()) {
	t.Helper()

	r := http.NewServeMux()
	r.HandleFunc(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotEmpty(t, r.Header.Get("Authorization"), "no authorization header present on the request")
		assert.Equal(t, "application/json", r.Header.Get("Accept"), "bad accept header present on the request")
		assert.Equal(t, "axiom-go", r.Header.Get("User-Agent"), "bad user-agent header present on the request")
		if orgIDHeader := r.Header.Get("X-Axiom-Org-Id"); orgIDHeader != "" {
			assert.Equal(t, orgID, orgIDHeader, "bad x-axiom-org-id header present on the request")
		}

		if r.ContentLength > 0 {
			assert.NotEmpty(t, r.Header.Get("Content-Type"), "no Content-Type header present on the request")
		}

		handler.ServeHTTP(w, r)
	}))
	srv := httptest.NewServer(r)

	client, err := NewClient(
		SetURL(srv.URL),
		SetAccessToken(personalToken),
		SetOrgID(orgID),
		SetClient(srv.Client()),
		SetStrictDecoding(true),
		SetNoEnv(),
	)
	require.NoError(t, err)

	return client, func() { srv.Close() }
}

// newClient returns a new client with stub properties for testing methods that
// don't actually make a http call.
func newClient(t *testing.T) *Client {
	t.Helper()

	client, err := NewClient(
		SetURL(endpoint),
		SetAccessToken(personalToken),
		SetNoEnv(),
	)
	require.NoError(t, err)

	return client
}

func mustTimeParse(tb testing.TB, layout, value string) time.Time {
	ts, err := time.Parse(layout, value)
	require.NoError(tb, err)
	return ts
}
