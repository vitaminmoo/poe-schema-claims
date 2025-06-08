package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

var enumStore = sync.Map{}

// PUT /v1/enums - Create a new enum
func (s server) PutEnums(ctx context.Context, request PutEnumsRequestObject) (PutEnumsResponseObject, error) {
	enum := Enum(*request.Body)
	id, err := ID()
	if err != nil {
		return nil, err
	}
	enum.Id = &id
	source := "FIXME"
	enum.Source = &source

	enumStore.Store(id, enum)
	return PutEnums201JSONResponse{
		Body: enum,
		Headers: PutEnums201ResponseHeaders{
			Location: fmt.Sprintf("/enums/%s", id),
		},
	}, nil
}

// PUT /v1/enums/{id} - Update an enum
func (s server) PutEnumsId(ctx context.Context, request PutEnumsIdRequestObject) (PutEnumsIdResponseObject, error) {
	id := request.Id
	enumAny, ok := enumStore.Load(id)
	enum := enumAny.(Enum)
	if !ok {
		return PutEnumsId404JSONResponse{
			NotFoundJSONResponse: NotFoundJSONResponse{
				Message: "Enum not found",
				Code:    http.StatusNotFound,
			},
		}, nil
	}
	if request.Body.Name != nil {
		enum.Name = *request.Body.Name
	}
	if request.Body.Values != nil {
		enum.Values = *request.Body.Values
	}
	if request.Body.ZeroIndexed != nil {
		enum.ZeroIndexed = *request.Body.ZeroIndexed
	}
	if request.Body.ClientLabels != nil {
		enum.ClientLabels = *request.Body.ClientLabels
	}
	enumStore.Store(id, enum)
	return PutEnumsId200JSONResponse(enum), nil
}

// GET /v1/enums - List all enums
func (s server) GetEnums(ctx context.Context, request GetEnumsRequestObject) (GetEnumsResponseObject, error) {
	enums := GetEnums200JSONResponse{}
	enumStore.Range(func(key any, value any) bool {
		enums = append(enums, value.(Enum))
		return true
	})
	return enums, nil
}

// GET /v1/enums/{id} - Retrieve an enum by ID
func (s server) GetEnumsId(ctx context.Context, request GetEnumsIdRequestObject) (GetEnumsIdResponseObject, error) {
	if enumAny, ok := enumStore.Load(request.Id); ok {
		enum := enumAny.(Enum)
		return GetEnumsId200JSONResponse(enum), nil
	}
	return GetEnumsId404JSONResponse{
		NotFoundJSONResponse: NotFoundJSONResponse{
			Message: "Enum not found",
			Code:    http.StatusNotFound,
		},
	}, nil
}

// DELETE /v1/enums/{id} - Delete an enum by ID
func (s server) DeleteEnumsId(ctx context.Context, request DeleteEnumsIdRequestObject) (DeleteEnumsIdResponseObject, error) {
	if _, ok := enumStore.LoadAndDelete(request.Id); !ok {
		return DeleteEnumsId404JSONResponse{
			NotFoundJSONResponse: NotFoundJSONResponse{
				Message: "Enum not found",
				Code:    http.StatusNotFound,
			},
		}, nil
	}
	return DeleteEnumsId204Response{}, nil
}
