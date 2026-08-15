package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	supabase "github.com/supabase-community/supabase-go"
)

func TestPublicInfo(t *testing.T) {
	var db *supabase.Client
	r := New(db)

	req, err := http.NewRequest("GET", "/public/info", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	if resp["message"] != "Welcome stranger! This info is public." {
		t.Errorf("handler returned unexpected message: got %q want %q", resp["message"], "Welcome stranger! This info is public.")
	}
}

func TestProtectedProfile_MissingHeader(t *testing.T) {
	var db *supabase.Client
	r := New(db)

	req, err := http.NewRequest("GET", "/protected/profile", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	if resp["error"] != "Access token required" {
		t.Errorf("handler returned unexpected error: got %q want %q", resp["error"], "Access token required")
	}
}

func TestProtectedProfile_IncorrectFormat(t *testing.T) {
	var db *supabase.Client
	r := New(db)

	req, err := http.NewRequest("GET", "/protected/profile", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Token abc")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	if resp["error"] != "Access token required" {
		t.Errorf("handler returned unexpected error: got %q want %q", resp["error"], "Access token required")
	}
}

func TestProtectedProfile_EmptyToken(t *testing.T) {
	var db *supabase.Client
	r := New(db)

	req, err := http.NewRequest("GET", "/protected/profile", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer ")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}

	if resp["error"] != "Access token required" {
		t.Errorf("handler returned unexpected error: got %q want %q", resp["error"], "Access token required")
	}
}
