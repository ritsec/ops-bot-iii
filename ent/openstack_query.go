// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ritsec/ops-bot-iii/ent/openstack"
	"github.com/ritsec/ops-bot-iii/ent/predicate"
	"github.com/ritsec/ops-bot-iii/ent/user"
)

// OpenstackQuery is the builder for querying Openstack entities.
type OpenstackQuery struct {
	config
	ctx        *QueryContext
	order      []openstack.OrderOption
	inters     []Interceptor
	predicates []predicate.Openstack
	withUser   *UserQuery
	withFKs    bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the OpenstackQuery builder.
func (oq *OpenstackQuery) Where(ps ...predicate.Openstack) *OpenstackQuery {
	oq.predicates = append(oq.predicates, ps...)
	return oq
}

// Limit the number of records to be returned by this query.
func (oq *OpenstackQuery) Limit(limit int) *OpenstackQuery {
	oq.ctx.Limit = &limit
	return oq
}

// Offset to start from.
func (oq *OpenstackQuery) Offset(offset int) *OpenstackQuery {
	oq.ctx.Offset = &offset
	return oq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (oq *OpenstackQuery) Unique(unique bool) *OpenstackQuery {
	oq.ctx.Unique = &unique
	return oq
}

// Order specifies how the records should be ordered.
func (oq *OpenstackQuery) Order(o ...openstack.OrderOption) *OpenstackQuery {
	oq.order = append(oq.order, o...)
	return oq
}

// QueryUser chains the current query on the "user" edge.
func (oq *OpenstackQuery) QueryUser() *UserQuery {
	query := (&UserClient{config: oq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := oq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := oq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(openstack.Table, openstack.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, openstack.UserTable, openstack.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(oq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Openstack entity from the query.
// Returns a *NotFoundError when no Openstack was found.
func (oq *OpenstackQuery) First(ctx context.Context) (*Openstack, error) {
	nodes, err := oq.Limit(1).All(setContextOp(ctx, oq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{openstack.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (oq *OpenstackQuery) FirstX(ctx context.Context) *Openstack {
	node, err := oq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Openstack ID from the query.
// Returns a *NotFoundError when no Openstack ID was found.
func (oq *OpenstackQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = oq.Limit(1).IDs(setContextOp(ctx, oq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{openstack.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (oq *OpenstackQuery) FirstIDX(ctx context.Context) int {
	id, err := oq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Openstack entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Openstack entity is found.
// Returns a *NotFoundError when no Openstack entities are found.
func (oq *OpenstackQuery) Only(ctx context.Context) (*Openstack, error) {
	nodes, err := oq.Limit(2).All(setContextOp(ctx, oq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{openstack.Label}
	default:
		return nil, &NotSingularError{openstack.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (oq *OpenstackQuery) OnlyX(ctx context.Context) *Openstack {
	node, err := oq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Openstack ID in the query.
// Returns a *NotSingularError when more than one Openstack ID is found.
// Returns a *NotFoundError when no entities are found.
func (oq *OpenstackQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = oq.Limit(2).IDs(setContextOp(ctx, oq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{openstack.Label}
	default:
		err = &NotSingularError{openstack.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (oq *OpenstackQuery) OnlyIDX(ctx context.Context) int {
	id, err := oq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Openstacks.
func (oq *OpenstackQuery) All(ctx context.Context) ([]*Openstack, error) {
	ctx = setContextOp(ctx, oq.ctx, ent.OpQueryAll)
	if err := oq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Openstack, *OpenstackQuery]()
	return withInterceptors[[]*Openstack](ctx, oq, qr, oq.inters)
}

// AllX is like All, but panics if an error occurs.
func (oq *OpenstackQuery) AllX(ctx context.Context) []*Openstack {
	nodes, err := oq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Openstack IDs.
func (oq *OpenstackQuery) IDs(ctx context.Context) (ids []int, err error) {
	if oq.ctx.Unique == nil && oq.path != nil {
		oq.Unique(true)
	}
	ctx = setContextOp(ctx, oq.ctx, ent.OpQueryIDs)
	if err = oq.Select(openstack.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (oq *OpenstackQuery) IDsX(ctx context.Context) []int {
	ids, err := oq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (oq *OpenstackQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, oq.ctx, ent.OpQueryCount)
	if err := oq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, oq, querierCount[*OpenstackQuery](), oq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (oq *OpenstackQuery) CountX(ctx context.Context) int {
	count, err := oq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (oq *OpenstackQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, oq.ctx, ent.OpQueryExist)
	switch _, err := oq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (oq *OpenstackQuery) ExistX(ctx context.Context) bool {
	exist, err := oq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the OpenstackQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (oq *OpenstackQuery) Clone() *OpenstackQuery {
	if oq == nil {
		return nil
	}
	return &OpenstackQuery{
		config:     oq.config,
		ctx:        oq.ctx.Clone(),
		order:      append([]openstack.OrderOption{}, oq.order...),
		inters:     append([]Interceptor{}, oq.inters...),
		predicates: append([]predicate.Openstack{}, oq.predicates...),
		withUser:   oq.withUser.Clone(),
		// clone intermediate query.
		sql:  oq.sql.Clone(),
		path: oq.path,
	}
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (oq *OpenstackQuery) WithUser(opts ...func(*UserQuery)) *OpenstackQuery {
	query := (&UserClient{config: oq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	oq.withUser = query
	return oq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Timestamp time.Time `json:"timestamp,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Openstack.Query().
//		GroupBy(openstack.FieldTimestamp).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (oq *OpenstackQuery) GroupBy(field string, fields ...string) *OpenstackGroupBy {
	oq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &OpenstackGroupBy{build: oq}
	grbuild.flds = &oq.ctx.Fields
	grbuild.label = openstack.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Timestamp time.Time `json:"timestamp,omitempty"`
//	}
//
//	client.Openstack.Query().
//		Select(openstack.FieldTimestamp).
//		Scan(ctx, &v)
func (oq *OpenstackQuery) Select(fields ...string) *OpenstackSelect {
	oq.ctx.Fields = append(oq.ctx.Fields, fields...)
	sbuild := &OpenstackSelect{OpenstackQuery: oq}
	sbuild.label = openstack.Label
	sbuild.flds, sbuild.scan = &oq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a OpenstackSelect configured with the given aggregations.
func (oq *OpenstackQuery) Aggregate(fns ...AggregateFunc) *OpenstackSelect {
	return oq.Select().Aggregate(fns...)
}

func (oq *OpenstackQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range oq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, oq); err != nil {
				return err
			}
		}
	}
	for _, f := range oq.ctx.Fields {
		if !openstack.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if oq.path != nil {
		prev, err := oq.path(ctx)
		if err != nil {
			return err
		}
		oq.sql = prev
	}
	return nil
}

func (oq *OpenstackQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Openstack, error) {
	var (
		nodes       = []*Openstack{}
		withFKs     = oq.withFKs
		_spec       = oq.querySpec()
		loadedTypes = [1]bool{
			oq.withUser != nil,
		}
	)
	if oq.withUser != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, openstack.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Openstack).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Openstack{config: oq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, oq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := oq.withUser; query != nil {
		if err := oq.loadUser(ctx, query, nodes, nil,
			func(n *Openstack, e *User) { n.Edges.User = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (oq *OpenstackQuery) loadUser(ctx context.Context, query *UserQuery, nodes []*Openstack, init func(*Openstack), assign func(*Openstack, *User)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*Openstack)
	for i := range nodes {
		if nodes[i].user_openstack == nil {
			continue
		}
		fk := *nodes[i].user_openstack
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(user.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "user_openstack" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (oq *OpenstackQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := oq.querySpec()
	_spec.Node.Columns = oq.ctx.Fields
	if len(oq.ctx.Fields) > 0 {
		_spec.Unique = oq.ctx.Unique != nil && *oq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, oq.driver, _spec)
}

func (oq *OpenstackQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(openstack.Table, openstack.Columns, sqlgraph.NewFieldSpec(openstack.FieldID, field.TypeInt))
	_spec.From = oq.sql
	if unique := oq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if oq.path != nil {
		_spec.Unique = true
	}
	if fields := oq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, openstack.FieldID)
		for i := range fields {
			if fields[i] != openstack.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := oq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := oq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := oq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := oq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (oq *OpenstackQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(oq.driver.Dialect())
	t1 := builder.Table(openstack.Table)
	columns := oq.ctx.Fields
	if len(columns) == 0 {
		columns = openstack.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if oq.sql != nil {
		selector = oq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if oq.ctx.Unique != nil && *oq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range oq.predicates {
		p(selector)
	}
	for _, p := range oq.order {
		p(selector)
	}
	if offset := oq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := oq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// OpenstackGroupBy is the group-by builder for Openstack entities.
type OpenstackGroupBy struct {
	selector
	build *OpenstackQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ogb *OpenstackGroupBy) Aggregate(fns ...AggregateFunc) *OpenstackGroupBy {
	ogb.fns = append(ogb.fns, fns...)
	return ogb
}

// Scan applies the selector query and scans the result into the given value.
func (ogb *OpenstackGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ogb.build.ctx, ent.OpQueryGroupBy)
	if err := ogb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*OpenstackQuery, *OpenstackGroupBy](ctx, ogb.build, ogb, ogb.build.inters, v)
}

func (ogb *OpenstackGroupBy) sqlScan(ctx context.Context, root *OpenstackQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(ogb.fns))
	for _, fn := range ogb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*ogb.flds)+len(ogb.fns))
		for _, f := range *ogb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*ogb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ogb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// OpenstackSelect is the builder for selecting fields of Openstack entities.
type OpenstackSelect struct {
	*OpenstackQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (os *OpenstackSelect) Aggregate(fns ...AggregateFunc) *OpenstackSelect {
	os.fns = append(os.fns, fns...)
	return os
}

// Scan applies the selector query and scans the result into the given value.
func (os *OpenstackSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, os.ctx, ent.OpQuerySelect)
	if err := os.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*OpenstackQuery, *OpenstackSelect](ctx, os.OpenstackQuery, os, os.inters, v)
}

func (os *OpenstackSelect) sqlScan(ctx context.Context, root *OpenstackQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(os.fns))
	for _, fn := range os.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*os.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := os.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
