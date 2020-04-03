package metakube

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientDo(t *testing.T) {
	token := "mySecretCredentials"
	eMethod := "GET"
	eURL := "/test"
	eToken := "Bearer " + token
	eBody := "testing"
	eRespBody := "ok"
	var rMethod, rURL, rToken, rBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rMethod = r.Method
		rURL = r.URL.Path
		rToken = r.Header.Get("authorization")
		var buf bytes.Buffer
		buf.ReadFrom(r.Body)
		rBody = buf.String()

		w.Write([]byte(eRespBody))
		return
	}))
	defer ts.Close()

	client := NewClient(ts.URL, token)
	resp, err := client.do(eMethod, eURL, strings.NewReader(eBody))
	if err != nil {
		t.Fatalf("Error not expected: %v", err)
	}
	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	rRespBody := buf.String()

	if eMethod != rMethod {
		t.Fatalf("expected method '%s', found '%s'", eMethod, rMethod)
	}
	if eURL != rURL {
		t.Fatalf("expected URL '%s', found '%s'", eURL, rURL)
	}
	if eToken != rToken {
		t.Fatalf("expected token '%s', found '%s'", eToken, rToken)
	}
	if eBody != rBody {
		t.Fatalf("expected body '%s', found '%s'", eBody, rBody)
	}
	if eRespBody != rRespBody {
		t.Fatalf("expected body '%s', found '%s'", eRespBody, rRespBody)
	}
}
