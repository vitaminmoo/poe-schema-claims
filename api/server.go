package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

// StrictServerInterface being defined like this causes unimplemented handler methods to be a compilation error.
var _ StrictServerInterface = (*server)(nil)

type server struct{}

func NewServer() server {
	return server{}
}

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
	if request.Body.Labels != nil {
		enum.Labels = *request.Body.Labels
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

var columnclaimStore = sync.Map{}

// PUT /v1/columnclaims - Create a new columnclaim
func (s server) PutColumnclaims(ctx context.Context, request PutColumnclaimsRequestObject) (PutColumnclaimsResponseObject, error) {
	columnClaim := ColumnClaim(*request.Body)
	id, err := ID()
	if err != nil {
		return nil, err
	}
	columnClaim.Id = &id
	source := "FIXME"
	columnClaim.Source = &source
	columnclaimStore.Store(id, columnClaim)
	return PutColumnclaims201JSONResponse{
		Body: columnClaim,
		Headers: PutColumnclaims201ResponseHeaders{
			Location: fmt.Sprintf("/columnclaims/%s", id),
		},
	}, nil
}

// PUT /v1/columnclaims/{id} - Update an columnclaim
func (s server) PutColumnclaimsId(ctx context.Context, request PutColumnclaimsIdRequestObject) (PutColumnclaimsIdResponseObject, error) {
	id := request.Id
	columnclaimAny, ok := columnclaimStore.Load(id)
	columnclaim := columnclaimAny.(ColumnClaim)
	if !ok {
		return PutColumnclaimsId404JSONResponse{
			NotFoundJSONResponse: NotFoundJSONResponse{
				Message: "ColumnClaim not found",
				Code:    http.StatusNotFound,
			},
		}, nil
	}
	if request.Body.Name != nil {
		columnclaim.Name = *request.Body.Name
	}
	if request.Body.Offset != nil {
		columnclaim.Offset = *request.Body.Offset
	}
	if request.Body.Bytes != nil {
		columnclaim.Bytes = *request.Body.Bytes
	}
	if request.Body.IsArray != nil {
		columnclaim.IsArray = request.Body.IsArray
	}
	if request.Body.Column != nil {
		columnclaim.Column = *request.Body.Column
	}
	if request.Body.Datfile != nil {
		columnclaim.Datfile = *request.Body.Datfile
	}
	if request.Body.Labels != nil {
		columnclaim.Labels = *request.Body.Labels
	}

	columnclaimStore.Store(id, columnclaim)
	return PutColumnclaimsId200JSONResponse(columnclaim), nil
}

// GET /v1/columnclaims - List all columnclaims
func (s server) GetColumnclaims(ctx context.Context, request GetColumnclaimsRequestObject) (GetColumnclaimsResponseObject, error) {
	columnclaims := GetColumnclaims200JSONResponse{}
	columnclaimStore.Range(func(key any, value any) bool {
		columnclaims = append(columnclaims, value.(ColumnClaim))
		return true
	})
	return columnclaims, nil
}

// GET /v1/columnclaims/{id} - Retrieve an columnclaim by ID
func (s server) GetColumnclaimsId(ctx context.Context, request GetColumnclaimsIdRequestObject) (GetColumnclaimsIdResponseObject, error) {
	if columnclaimAny, ok := columnclaimStore.Load(request.Id); ok {
		columnclaim := columnclaimAny.(ColumnClaim)
		return GetColumnclaimsId200JSONResponse(columnclaim), nil
	}
	return GetColumnclaimsId404JSONResponse{
		NotFoundJSONResponse: NotFoundJSONResponse{
			Message: "ColumnClaim not found",
			Code:    http.StatusNotFound,
		},
	}, nil
}

// DELETE /v1/columnclaims/{id} - Delete an columnclaim by ID
func (s server) DeleteColumnclaimsId(ctx context.Context, request DeleteColumnclaimsIdRequestObject) (DeleteColumnclaimsIdResponseObject, error) {
	if _, ok := columnclaimStore.LoadAndDelete(request.Id); !ok {
		return DeleteColumnclaimsId404JSONResponse{
			NotFoundJSONResponse: NotFoundJSONResponse{
				Message: "ColulmnClaim not found",
				Code:    http.StatusNotFound,
			},
		}, nil
	}
	return DeleteColumnclaimsId204Response{}, nil
}
