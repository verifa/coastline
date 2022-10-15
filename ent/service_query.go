// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/verifa/coastline/ent/predicate"
	"github.com/verifa/coastline/ent/request"
	"github.com/verifa/coastline/ent/service"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ServiceQuery is the builder for querying Service entities.
type ServiceQuery struct {
	config
	limit        *int
	offset       *int
	unique       *bool
	order        []OrderFunc
	fields       []string
	predicates   []predicate.Service
	withRequests *RequestQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ServiceQuery builder.
func (sq *ServiceQuery) Where(ps ...predicate.Service) *ServiceQuery {
	sq.predicates = append(sq.predicates, ps...)
	return sq
}

// Limit adds a limit step to the query.
func (sq *ServiceQuery) Limit(limit int) *ServiceQuery {
	sq.limit = &limit
	return sq
}

// Offset adds an offset step to the query.
func (sq *ServiceQuery) Offset(offset int) *ServiceQuery {
	sq.offset = &offset
	return sq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (sq *ServiceQuery) Unique(unique bool) *ServiceQuery {
	sq.unique = &unique
	return sq
}

// Order adds an order step to the query.
func (sq *ServiceQuery) Order(o ...OrderFunc) *ServiceQuery {
	sq.order = append(sq.order, o...)
	return sq
}

// QueryRequests chains the current query on the "Requests" edge.
func (sq *ServiceQuery) QueryRequests() *RequestQuery {
	query := &RequestQuery{config: sq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := sq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(service.Table, service.FieldID, selector),
			sqlgraph.To(request.Table, request.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, service.RequestsTable, service.RequestsPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(sq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Service entity from the query.
// Returns a *NotFoundError when no Service was found.
func (sq *ServiceQuery) First(ctx context.Context) (*Service, error) {
	nodes, err := sq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{service.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (sq *ServiceQuery) FirstX(ctx context.Context) *Service {
	node, err := sq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Service ID from the query.
// Returns a *NotFoundError when no Service ID was found.
func (sq *ServiceQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = sq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{service.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (sq *ServiceQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := sq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Service entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Service entity is found.
// Returns a *NotFoundError when no Service entities are found.
func (sq *ServiceQuery) Only(ctx context.Context) (*Service, error) {
	nodes, err := sq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{service.Label}
	default:
		return nil, &NotSingularError{service.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (sq *ServiceQuery) OnlyX(ctx context.Context) *Service {
	node, err := sq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Service ID in the query.
// Returns a *NotSingularError when more than one Service ID is found.
// Returns a *NotFoundError when no entities are found.
func (sq *ServiceQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = sq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{service.Label}
	default:
		err = &NotSingularError{service.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (sq *ServiceQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := sq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Services.
func (sq *ServiceQuery) All(ctx context.Context) ([]*Service, error) {
	if err := sq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return sq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (sq *ServiceQuery) AllX(ctx context.Context) []*Service {
	nodes, err := sq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Service IDs.
func (sq *ServiceQuery) IDs(ctx context.Context) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	if err := sq.Select(service.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (sq *ServiceQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := sq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (sq *ServiceQuery) Count(ctx context.Context) (int, error) {
	if err := sq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return sq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (sq *ServiceQuery) CountX(ctx context.Context) int {
	count, err := sq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (sq *ServiceQuery) Exist(ctx context.Context) (bool, error) {
	if err := sq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return sq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (sq *ServiceQuery) ExistX(ctx context.Context) bool {
	exist, err := sq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ServiceQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (sq *ServiceQuery) Clone() *ServiceQuery {
	if sq == nil {
		return nil
	}
	return &ServiceQuery{
		config:       sq.config,
		limit:        sq.limit,
		offset:       sq.offset,
		order:        append([]OrderFunc{}, sq.order...),
		predicates:   append([]predicate.Service{}, sq.predicates...),
		withRequests: sq.withRequests.Clone(),
		// clone intermediate query.
		sql:    sq.sql.Clone(),
		path:   sq.path,
		unique: sq.unique,
	}
}

// WithRequests tells the query-builder to eager-load the nodes that are connected to
// the "Requests" edge. The optional arguments are used to configure the query builder of the edge.
func (sq *ServiceQuery) WithRequests(opts ...func(*RequestQuery)) *ServiceQuery {
	query := &RequestQuery{config: sq.config}
	for _, opt := range opts {
		opt(query)
	}
	sq.withRequests = query
	return sq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Service.Query().
//		GroupBy(service.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (sq *ServiceQuery) GroupBy(field string, fields ...string) *ServiceGroupBy {
	grbuild := &ServiceGroupBy{config: sq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := sq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return sq.sqlQuery(ctx), nil
	}
	grbuild.label = service.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.Service.Query().
//		Select(service.FieldName).
//		Scan(ctx, &v)
//
func (sq *ServiceQuery) Select(fields ...string) *ServiceSelect {
	sq.fields = append(sq.fields, fields...)
	selbuild := &ServiceSelect{ServiceQuery: sq}
	selbuild.label = service.Label
	selbuild.flds, selbuild.scan = &sq.fields, selbuild.Scan
	return selbuild
}

func (sq *ServiceQuery) prepareQuery(ctx context.Context) error {
	for _, f := range sq.fields {
		if !service.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if sq.path != nil {
		prev, err := sq.path(ctx)
		if err != nil {
			return err
		}
		sq.sql = prev
	}
	return nil
}

func (sq *ServiceQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Service, error) {
	var (
		nodes       = []*Service{}
		_spec       = sq.querySpec()
		loadedTypes = [1]bool{
			sq.withRequests != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Service).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Service{config: sq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, sq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := sq.withRequests; query != nil {
		if err := sq.loadRequests(ctx, query, nodes,
			func(n *Service) { n.Edges.Requests = []*Request{} },
			func(n *Service, e *Request) { n.Edges.Requests = append(n.Edges.Requests, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (sq *ServiceQuery) loadRequests(ctx context.Context, query *RequestQuery, nodes []*Service, init func(*Service), assign func(*Service, *Request)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[uuid.UUID]*Service)
	nids := make(map[uuid.UUID]map[*Service]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(service.RequestsTable)
		s.Join(joinT).On(s.C(request.FieldID), joinT.C(service.RequestsPrimaryKey[0]))
		s.Where(sql.InValues(joinT.C(service.RequestsPrimaryKey[1]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(service.RequestsPrimaryKey[1]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	neighbors, err := query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
		assign := spec.Assign
		values := spec.ScanValues
		spec.ScanValues = func(columns []string) ([]any, error) {
			values, err := values(columns[1:])
			if err != nil {
				return nil, err
			}
			return append([]any{new(uuid.UUID)}, values...), nil
		}
		spec.Assign = func(columns []string, values []any) error {
			outValue := *values[0].(*uuid.UUID)
			inValue := *values[1].(*uuid.UUID)
			if nids[inValue] == nil {
				nids[inValue] = map[*Service]struct{}{byID[outValue]: struct{}{}}
				return assign(columns[1:], values[1:])
			}
			nids[inValue][byID[outValue]] = struct{}{}
			return nil
		}
	})
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "Requests" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}

func (sq *ServiceQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := sq.querySpec()
	_spec.Node.Columns = sq.fields
	if len(sq.fields) > 0 {
		_spec.Unique = sq.unique != nil && *sq.unique
	}
	return sqlgraph.CountNodes(ctx, sq.driver, _spec)
}

func (sq *ServiceQuery) sqlExist(ctx context.Context) (bool, error) {
	switch _, err := sq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

func (sq *ServiceQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   service.Table,
			Columns: service.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: service.FieldID,
			},
		},
		From:   sq.sql,
		Unique: true,
	}
	if unique := sq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := sq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, service.FieldID)
		for i := range fields {
			if fields[i] != service.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := sq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := sq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := sq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := sq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (sq *ServiceQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(sq.driver.Dialect())
	t1 := builder.Table(service.Table)
	columns := sq.fields
	if len(columns) == 0 {
		columns = service.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if sq.sql != nil {
		selector = sq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if sq.unique != nil && *sq.unique {
		selector.Distinct()
	}
	for _, p := range sq.predicates {
		p(selector)
	}
	for _, p := range sq.order {
		p(selector)
	}
	if offset := sq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := sq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ServiceGroupBy is the group-by builder for Service entities.
type ServiceGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (sgb *ServiceGroupBy) Aggregate(fns ...AggregateFunc) *ServiceGroupBy {
	sgb.fns = append(sgb.fns, fns...)
	return sgb
}

// Scan applies the group-by query and scans the result into the given value.
func (sgb *ServiceGroupBy) Scan(ctx context.Context, v any) error {
	query, err := sgb.path(ctx)
	if err != nil {
		return err
	}
	sgb.sql = query
	return sgb.sqlScan(ctx, v)
}

func (sgb *ServiceGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range sgb.fields {
		if !service.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := sgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := sgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (sgb *ServiceGroupBy) sqlQuery() *sql.Selector {
	selector := sgb.sql.Select()
	aggregation := make([]string, 0, len(sgb.fns))
	for _, fn := range sgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(sgb.fields)+len(sgb.fns))
		for _, f := range sgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(sgb.fields...)...)
}

// ServiceSelect is the builder for selecting fields of Service entities.
type ServiceSelect struct {
	*ServiceQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (ss *ServiceSelect) Scan(ctx context.Context, v any) error {
	if err := ss.prepareQuery(ctx); err != nil {
		return err
	}
	ss.sql = ss.ServiceQuery.sqlQuery(ctx)
	return ss.sqlScan(ctx, v)
}

func (ss *ServiceSelect) sqlScan(ctx context.Context, v any) error {
	rows := &sql.Rows{}
	query, args := ss.sql.Query()
	if err := ss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
