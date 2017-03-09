package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/aabizri/goil.v0"
	"net/http"
)

// Group ID - Group Name
type groupsResponse map[goil.Group]string

func getGroups(w http.ResponseWriter, req *http.Request) {
	// Retrieve goil session
	raw := req.Context().Value("session")
	gs, ok := raw.(*goil.Session)
	if !ok {
		err := NewError(InternalErrorTypeAssertionFailed, "", "Type assertion failure while casting req.Context().Value(\"session\") to (*goil.Session)").JSONWrite(w)
		if err != nil {
			panic(err)
		}
		return
	}

	// Execute the request
	gr, err := gs.AvailableGroups()
	if err != nil {
		err := NewError(InternalError, "Error while calling gs.AvailableGroups()", "").JSONWrite(w)
		if err != nil {
			panic(err)
		}
		return
	}

	// Marhsal it
	jresp, err := json.Marshal(gr)
	if err != nil {
		err := NewError(InternalErrorJSONMarshallingFailed, "", "").JSONWrite(w)
		if err != nil {
			panic(err)
		}
		return
	}

	// Return it
	_, err = bytes.NewBuffer(jresp).WriteTo(w)
	if err != nil {
		err := NewError(InternalErrorWritingResponseWriterFailed, "", "").JSONWrite(w)
		if err != nil {
			panic(err)
		}
		return
	}
}
