package hw_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ogiusek/hw/src/hw"
)

type wrapStructCommand struct {
	Count int
}

func wrapStructCommandHandler(c wrapStructCommand) any {
	c.Count += 1
	return c
}

func Test_POST_wrap_should_return_command(t *testing.T) {
	// arrange
	reqCommand := wrapStructCommand{
		Count: 2,
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	// Encode struct and write to buffer
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
	handler := http.HandlerFunc(hw.Wrap(hw.Run(wrapStructCommandHandler), nil))
	handler.ServeHTTP(rr, req)

	// assert
	if rr.Code != http.StatusOK {
		t.Errorf("wrap response code should be %v instead of %v", http.StatusOK, rr.Code)
	}

	var args wrapStructCommand
	decoder := json.NewDecoder(rr.Body)
	if err := decoder.Decode(&args); err != nil {
		t.Errorf("wrap response cannot be decoded")
	} else {
		reqCommand.Count += 1
		if args.Count != reqCommand.Count {
			t.Errorf("wrap do not get executed command handler properly. command handler do not changed command value\ngot : %v\nwant: %v", args, reqCommand)
		}
	}
}

func Test_GET_wrap_should_return_command(t *testing.T) {
	// arrange
	reqCommand := wrapStructCommand{
		Count: 2,
	}

	req, err := http.NewRequest("GET", "/health-check?count=2", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// act
	handler := http.HandlerFunc(hw.Wrap(hw.Run(wrapStructCommandHandler), nil))
	handler.ServeHTTP(rr, req)

	// assert
	if rr.Code != http.StatusOK {
		t.Errorf("wrap response code should be %v instead of %v", http.StatusOK, rr.Code)
	}

	var args wrapStructCommand
	decoder := json.NewDecoder(rr.Body)
	if err := decoder.Decode(&args); err != nil {
		t.Errorf("wrap response cannot be decoded")
	} else {
		reqCommand.Count += 1
		if args.Count != reqCommand.Count {
			t.Errorf("wrap do not get executed command handler properly. command handler do not changed command value\ngot : %v\nwant: %v", args, reqCommand)
		}
	}
}
