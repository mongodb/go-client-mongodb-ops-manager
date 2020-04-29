package opsmngr

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func TestDiagnostics_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	groupID := "5c8100bcf2a30b12ff88258f"

	path := fmt.Sprintf("/groups/%s/diagnostics", groupID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, "test")
	})

	buf := new(bytes.Buffer)
	_, err := client.Diagnostics.Get(ctx, groupID, buf)
	if err != nil {
		t.Fatalf("Diagnostics.Get returned error: %v", err)
	}

	if buf.String() != "test" {
		t.Fatalf("Diagnostics.Get returned error: %v", err)
	}

}

func TestDiagnostics_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/diagnostics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, "test")
	})

	buf := new(bytes.Buffer)
	_, err := client.Diagnostics.List(ctx, buf)
	if err != nil {
		t.Fatalf("Diagnostics.List returned error: %v", err)
	}

	if buf.String() != "test" {
		t.Fatalf("Diagnostics.List returned error: %v", err)
	}

}
