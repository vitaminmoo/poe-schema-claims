package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	jet "github.com/go-jet/jet/v2/sqlite"
	"github.com/ncruces/go-sqlite3"
	"github.com/vitaminmoo/poe-schema-claims/log"
	"github.com/vitaminmoo/poe-schema-claims/storage/table"
	"go.uber.org/zap"
)

// PUT /v1/enums - Create a new Enum
func (s server) PutEnums(ctx context.Context, request PutEnumsRequestObject) (PutEnumsResponseObject, error) {
	enum := Enum(*request.Body)
	source := "FIXME"
	enum.Source = &source
	enum.ServerLabels = &map[string]string{}

	stmt := table.Enum.INSERT(
		table.Enum.Source,
		table.Enum.Name,
		table.Enum.ClientLabels,
		table.Enum.ServerLabels,
		table.Enum.Vals,
		table.Enum.ZeroIndexed,
	).VALUES(
		source,
		enum.Name,
		sqlite3.JSON(enum.ClientLabels),
		sqlite3.JSON(enum.ServerLabels),
		sqlite3.JSON(enum.Values),
		enum.ZeroIndexed,
	)
	query, args := stmt.Sql()

	ret, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	id, err := ret.LastInsertId()
	if err != nil {
		return nil, err
	}
	enum.Id = &id

	return PutEnums201JSONResponse{
		Body: enum,
		Headers: PutEnums201ResponseHeaders{
			Location: fmt.Sprintf("/enums/%d", *enum.Id),
		},
	}, nil
}

// PUT /v1/enums/{id} - Update an Enum
func (s server) PutEnumsId(ctx context.Context, request PutEnumsIdRequestObject) (PutEnumsIdResponseObject, error) {
	cols := jet.ColumnList{}
	vals := []any{}

	if request.Body.ClientLabels != nil {
		cols = append(cols, table.Enum.ClientLabels)
		vals = append(vals, sqlite3.JSON(*request.Body.ClientLabels))
	}
	if request.Body.Name != nil {
		cols = append(cols, table.Enum.Name)
		vals = append(vals, *request.Body.Name)
	}
	if request.Body.Values != nil {
		cols = append(cols, table.Enum.Vals)
		vals = append(vals, sqlite3.JSON(*request.Body.Values))
	}
	if request.Body.ZeroIndexed != nil {
		cols = append(cols, table.Enum.ZeroIndexed)
		vals = append(vals, *request.Body.ZeroIndexed)
	}

	if len(cols) == 0 || len(vals) == 0 {
		return PutEnumsId304JSONResponse{}, nil
	}
	z := log.Load(ctx)
	z.Info("Updating Enum", zap.Int("cols", len(cols)), zap.Int("vals", len(vals)))

	stmt := table.Enum.UPDATE(cols).SET(vals[0], vals[1:]...).WHERE(table.Enum.ID.EQ(jet.Int64(request.Id)))
	query, args := stmt.Sql()
	ret, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("updating Enum: %w", err)
	}

	num, err := ret.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("getting number of affected rows: %w", err)
	}
	if num == 0 {
		return PutEnumsId404JSONResponse{
			NotFoundJSONResponse: NotFoundJSONResponse{
				Message: "Enum not found",
				Code:    http.StatusNotFound,
			},
		}, nil
	}
	return PutEnumsId204Response{}, nil
}

// GET /v1/enums - List all Enums
func (s server) GetEnums(ctx context.Context, request GetEnumsRequestObject) (GetEnumsResponseObject, error) {
	enums := GetEnums200JSONResponse{}

	stmt := jet.SELECT(table.Enum.AllColumns).FROM(table.Enum)
	query, args := stmt.Sql()
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing Enums: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var enum Enum
		if err := rows.Scan(&enum.Id, &enum.Source, &enum.Name, sqlite3.JSON(&enum.ClientLabels), sqlite3.JSON(&enum.ServerLabels), sqlite3.JSON(&enum.Values), &enum.ZeroIndexed); err != nil {
			return nil, fmt.Errorf("iterating Enums: %w", err)
		}
		z := log.Load(ctx)
		z.Info("Loaded enum", zap.Int64p("id", enum.Id), zap.Any("values", enum.Values))
		enums = append(enums, enum)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("querying Enums: %w", err)
	}

	return enums, nil
}

// GET /v1/enums/{id} - Retrieve an Enum by ID
func (s server) GetEnumsId(ctx context.Context, request GetEnumsIdRequestObject) (GetEnumsIdResponseObject, error) {
	source := ""
	enum := Enum{
		Id:           &request.Id,
		Source:       &source,
		ClientLabels: map[string]string{},
		ServerLabels: &map[string]string{},
		Values:       []string{},
		ZeroIndexed:  false,
	}

	stmt := jet.SELECT(table.Enum.AllColumns).FROM(table.Enum).WHERE(table.Enum.ID.EQ(jet.Int64(request.Id)))
	query, args := stmt.Sql()

	err := s.db.QueryRowContext(ctx, query, args...).Scan(enum.Id, enum.Source, &enum.Name, sqlite3.JSON(&enum.ClientLabels), sqlite3.JSON(&enum.ServerLabels), sqlite3.JSON(&enum.Values), &enum.ZeroIndexed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return GetEnumsId404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{
					Message: "Enum not found",
					Code:    http.StatusNotFound,
				},
			}, nil
		}
		return nil, fmt.Errorf("querying Enum by ID: %w", err)
	}
	return GetEnumsId200JSONResponse(enum), nil
}

// DELETE /v1/enums/{id} - Delete an Enum by ID
func (s server) DeleteEnumsId(ctx context.Context, request DeleteEnumsIdRequestObject) (DeleteEnumsIdResponseObject, error) {
	stmt := table.Enum.DELETE().WHERE(table.Enum.ID.EQ(jet.Int64(request.Id)))
	query, args := stmt.Sql()
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("deleting Enum by ID: %w", err)
	}
	num, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("getting number of affected rows: %w", err)
	}
	if num == 0 {
		return DeleteEnumsId404JSONResponse{
			NotFoundJSONResponse: NotFoundJSONResponse{
				Message: "Enum not found",
				Code:    http.StatusNotFound,
			},
		}, nil
	}

	return DeleteEnumsId204Response{}, nil
}
