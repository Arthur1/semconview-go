package semconview

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "go.opentelemetry.io/otel/semconv/v1.10.0"
	_ "go.opentelemetry.io/otel/semconv/v1.20.0"
	_ "go.opentelemetry.io/otel/semconv/v1.30.0"
)

func TestSemconvPackageResolver_resolveAttrKeyFromConst(t *testing.T) {
	r := newSemconvPackageResolver()
	ctx := context.Background()

	cases := []struct {
		inPkgPath   string
		inConstName string
		wantAttrKey string
		wantOK      bool
		wantErr     bool
	}{
		{
			inPkgPath:   "go.opentelemetry.io/otel/semconv/v1.10.0",
			inConstName: "HTTPStatusCodeKey",
			wantAttrKey: "http.status_code",
			wantOK:      true,
		},
		{
			inPkgPath:   "go.opentelemetry.io/otel/semconv/v1.20.0",
			inConstName: "HTTPStatusCodeKey",
			wantAttrKey: "http.status_code",
			wantOK:      true,
		},
		{
			inPkgPath:   "go.opentelemetry.io/otel/semconv/v1.30.0",
			inConstName: "HTTPStatusCodeKey",
		},
		{
			inPkgPath:   "go.opentelemetry.io/otel/semconv/v1.10.0",
			inConstName: "HTTPResponseStatusCodeKey",
		},
		{
			inPkgPath:   "go.opentelemetry.io/otel/semconv/v1.20.0",
			inConstName: "HTTPResponseStatusCodeKey",
		},
		{
			inPkgPath:   "go.opentelemetry.io/otel/semconv/v1.30.0",
			inConstName: "HTTPResponseStatusCodeKey",
			wantAttrKey: "http.response.status_code",
			wantOK:      true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.inConstName+" in "+filepath.Base(tt.inPkgPath), func(t *testing.T) {
			gotAttrKey, gotOK, gotErr := r.resolveAttrKeyFromConst(ctx, tt.inPkgPath, tt.inConstName)
			assert.Equal(t, tt.wantAttrKey, gotAttrKey)
			assert.Equal(t, tt.wantOK, gotOK)
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}

func TestSemconvPackageResolver_resolveAttrKeyFromFunc(t *testing.T) {
	r := newSemconvPackageResolver()
	ctx := context.Background()

	cases := []struct {
		inPkgPath   string
		inFuncName  string
		wantAttrKey string
		wantOK      bool
		wantErr     bool
	}{
		{
			inPkgPath:  "go.opentelemetry.io/otel/semconv/v1.10.0",
			inFuncName: "HTTPStatusCode",
		},
		{
			inPkgPath:   "go.opentelemetry.io/otel/semconv/v1.20.0",
			inFuncName:  "HTTPStatusCode",
			wantAttrKey: "http.status_code",
			wantOK:      true,
		},
		{
			inPkgPath:  "go.opentelemetry.io/otel/semconv/v1.30.0",
			inFuncName: "HTTPStatusCode",
		},
		{
			inPkgPath:  "go.opentelemetry.io/otel/semconv/v1.10.0",
			inFuncName: "HTTPResponseStatusCode",
		},
		{
			inPkgPath:  "go.opentelemetry.io/otel/semconv/v1.20.0",
			inFuncName: "HTTPResponseStatusCode",
		},
		{
			inPkgPath:   "go.opentelemetry.io/otel/semconv/v1.30.0",
			inFuncName:  "HTTPResponseStatusCode",
			wantAttrKey: "http.response.status_code",
			wantOK:      true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.inFuncName+"() in "+filepath.Base(tt.inPkgPath), func(t *testing.T) {
			gotAttrKey, gotOK, gotErr := r.resolveAttrKeyFromFunc(ctx, tt.inPkgPath, tt.inFuncName)
			assert.Equal(t, tt.wantAttrKey, gotAttrKey)
			assert.Equal(t, tt.wantOK, gotOK)
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}
