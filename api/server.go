package api

import "context"

// StrictServerInterface being defined like this causes unimplemented handler methods to be a compilation error.
var _ StrictServerInterface = (*server)(nil)

type server struct {
}

// PutColumnclaims implements StrictServerInterface.
func (s server) PutColumnclaims(ctx context.Context, request PutColumnclaimsRequestObject) (PutColumnclaimsResponseObject, error) {
	return PutColumnclaims201JSONResponse{}, nil
}

// PutColumnclaimsId implements StrictServerInterface.
func (s server) PutColumnclaimsId(ctx context.Context, request PutColumnclaimsIdRequestObject) (PutColumnclaimsIdResponseObject, error) {
	return PutColumnclaimsId200JSONResponse{}, nil
}

// PutEnums implements StrictServerInterface.
func (s server) PutEnums(ctx context.Context, request PutEnumsRequestObject) (PutEnumsResponseObject, error) {
	return PutEnums201JSONResponse{}, nil
}

// PutEnumsId implements StrictServerInterface.
func (s server) PutEnumsId(ctx context.Context, request PutEnumsIdRequestObject) (PutEnumsIdResponseObject, error) {
	return PutEnumsId200JSONResponse{}, nil
}

// DeleteColumnclaimsId implements StrictServerInterface.
func (s server) DeleteColumnclaimsId(ctx context.Context, request DeleteColumnclaimsIdRequestObject) (DeleteColumnclaimsIdResponseObject, error) {
	return DeleteColumnclaimsId204Response{}, nil
}

// DeleteEnumsId implements StrictServerInterface.
func (s server) DeleteEnumsId(ctx context.Context, request DeleteEnumsIdRequestObject) (DeleteEnumsIdResponseObject, error) {
	return DeleteEnumsId204Response{}, nil
}

// GetColumnclaims implements StrictServerInterface.
func (s server) GetColumnclaims(ctx context.Context, request GetColumnclaimsRequestObject) (GetColumnclaimsResponseObject, error) {
	return GetColumnclaims200JSONResponse{}, nil
}

// GetColumnclaimsId implements StrictServerInterface.
func (s server) GetColumnclaimsId(ctx context.Context, request GetColumnclaimsIdRequestObject) (GetColumnclaimsIdResponseObject, error) {
	return GetColumnclaimsId200JSONResponse{}, nil
}

// GetEnums implements StrictServerInterface.
func (s server) GetEnums(ctx context.Context, request GetEnumsRequestObject) (GetEnumsResponseObject, error) {
	return GetEnums200JSONResponse{}, nil
}

// GetEnumsId implements StrictServerInterface.
func (s server) GetEnumsId(ctx context.Context, request GetEnumsIdRequestObject) (GetEnumsIdResponseObject, error) {
	return GetEnumsId200JSONResponse{}, nil
}

func NewServer() server {
	return server{}
}
