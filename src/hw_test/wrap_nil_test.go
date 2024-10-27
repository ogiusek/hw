package hw_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ogiusek/hw/src/hw"
)

type wrapNilCommand struct {
}

func wrapNilCommandHandler(wrapNilCommand) any {
	return nil
}

func Test_should_return_no_content(t *testing.T) {
	// arrange
	reqCommand := wrapNilCommand{}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	err := encoder.Encode(reqCommand)
	if err != nil {
		t.Errorf("unexpected test error: %v", err)
		return
	}

	req, err := http.NewRequest("GET", "/health-check", &buf)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// act
	handler := http.HandlerFunc(hw.Wrap(hw.Run(wrapNilCommandHandler), nil))
	handler.ServeHTTP(rr, req)

	// assert
	if rr.Code != http.StatusNoContent {
		t.Errorf("wrap response code should be %v instead of %v", http.StatusNoContent, rr.Code)
	}
}
