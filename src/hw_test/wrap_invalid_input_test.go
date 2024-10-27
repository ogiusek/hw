package hw

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ogiusek/hw/src/hw"
)

type invalidInputCommand struct {
}

func invalidInputCommandHandler(invalidInputCommand) any {
	return nil
}

func Test_POST_wrap_should_return_command(t *testing.T) {
	// arrange
	reqCommand := invalidInputCommand{}
	var buf bytes.Buffer
	buf.WriteString(`{"a": "}`)
	encoder := json.NewEncoder(&buf)

	err := encoder.Encode(reqCommand)
	if err != nil {
		t.Errorf("unexpected test error: %v", err)
		return
	}

	req, err := http.NewRequest("POST", "/health-check", &buf)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// act
	handler := http.HandlerFunc(hw.Wrap(hw.Run(invalidInputCommandHandler), nil))
	handler.ServeHTTP(rr, req)

	// assert
	if rr.Code != http.StatusUnsupportedMediaType {
		t.Errorf("wrap response code should be %v instead of %v", http.StatusUnsupportedMediaType, rr.Code)
	}
}
