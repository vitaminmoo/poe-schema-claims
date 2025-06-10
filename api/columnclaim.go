package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

var columnclaimStore = sync.Map{}

// PUT /v1/columnclaims - Create a new columnclaim
func (s server) PutColumnClaims(ctx context.Context, request PutColumnClaimsRequestObject) (PutColumnClaimsResponseObject, error) {
	columnClaim := ColumnClaim(*request.Body)
	id := int64(1)
	columnClaim.Id = &id
	source := "FIXME"
	columnClaim.Source = &source
	columnclaimStore.Store(id, columnClaim)
	return PutColumnClaims201JSONResponse{
		Body: columnClaim,
		Headers: PutColumnClaims201ResponseHeaders{
			Location: fmt.Sprintf("/column_claims/%d", id),
		},
	}, nil
}

// PUT /v1/column_claims/{id} - Update an columnclaim
func (s server) PutColumnClaimsId(ctx context.Context, request PutColumnClaimsIdRequestObject) (PutColumnClaimsIdResponseObject, error) {
	id := request.Id
	columnclaimAny, ok := columnclaimStore.Load(id)
	columnclaim := columnclaimAny.(ColumnClaim)
	if !ok {
		return PutColumnClaimsId404JSONResponse{
			NotFoundJSONResponse: NotFoundJSONResponse{
				Message: "ColumnClaim not found",
				Code:    http.StatusNotFound,
			},
		}, nil
	}
	if request.Body.Name != nil {
		columnclaim.Name = request.Body.Name
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
	if request.Body.ClientLabels != nil {
		columnclaim.ClientLabels = *request.Body.ClientLabels
	}

	columnclaimStore.Store(id, columnclaim)
	return PutColumnClaimsId204Response{}, nil
}

// GET /v1/columnclaims - List all columnclaims
func (s server) GetColumnClaims(ctx context.Context, request GetColumnClaimsRequestObject) (GetColumnClaimsResponseObject, error) {
	columnclaims := GetColumnClaims200JSONResponse{}
	columnclaimStore.Range(func(key any, value any) bool {
		columnclaims = append(columnclaims, value.(ColumnClaim))
		return true
	})
	return columnclaims, nil
}

// GET /v1/column_claims/{id} - Retrieve an columnclaim by ID
func (s server) GetColumnClaimsId(ctx context.Context, request GetColumnClaimsIdRequestObject) (GetColumnClaimsIdResponseObject, error) {
	if columnclaimAny, ok := columnclaimStore.Load(request.Id); ok {
		columnclaim := columnclaimAny.(ColumnClaim)
		return GetColumnClaimsId200JSONResponse(columnclaim), nil
	}
	return GetColumnClaimsId404JSONResponse{
		NotFoundJSONResponse: NotFoundJSONResponse{
			Message: "ColumnClaim not found",
			Code:    http.StatusNotFound,
		},
	}, nil
}

// DELETE /v1/column_claims/{id} - Delete an columnclaim by ID
func (s server) DeleteColumnClaimsId(ctx context.Context, request DeleteColumnClaimsIdRequestObject) (DeleteColumnClaimsIdResponseObject, error) {
	if _, ok := columnclaimStore.LoadAndDelete(request.Id); !ok {
		return DeleteColumnClaimsId404JSONResponse{
			NotFoundJSONResponse: NotFoundJSONResponse{
				Message: "ColulmnClaim not found",
				Code:    http.StatusNotFound,
			},
		}, nil
	}
	return DeleteColumnClaimsId204Response{}, nil
}
