// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ritsec/ops-bot-iii/ent/openstack"
	"github.com/ritsec/ops-bot-iii/ent/predicate"
	"github.com/ritsec/ops-bot-iii/ent/user"
)

// OpenstackUpdate is the builder for updating Openstack entities.
type OpenstackUpdate struct {
	config
	hooks    []Hook
	mutation *OpenstackMutation
}

// Where appends a list predicates to the OpenstackUpdate builder.
func (ou *OpenstackUpdate) Where(ps ...predicate.Openstack) *OpenstackUpdate {
	ou.mutation.Where(ps...)
	return ou
}

// SetTimestamp sets the "timestamp" field.
func (ou *OpenstackUpdate) SetTimestamp(t time.Time) *OpenstackUpdate {
	ou.mutation.SetTimestamp(t)
	return ou
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (ou *OpenstackUpdate) SetNillableTimestamp(t *time.Time) *OpenstackUpdate {
	if t != nil {
		ou.SetTimestamp(*t)
	}
	return ou
}

// SetUserID sets the "user" edge to the User entity by ID.
func (ou *OpenstackUpdate) SetUserID(id string) *OpenstackUpdate {
	ou.mutation.SetUserID(id)
	return ou
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (ou *OpenstackUpdate) SetNillableUserID(id *string) *OpenstackUpdate {
	if id != nil {
		ou = ou.SetUserID(*id)
	}
	return ou
}

// SetUser sets the "user" edge to the User entity.
func (ou *OpenstackUpdate) SetUser(u *User) *OpenstackUpdate {
	return ou.SetUserID(u.ID)
}

// Mutation returns the OpenstackMutation object of the builder.
func (ou *OpenstackUpdate) Mutation() *OpenstackMutation {
	return ou.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (ou *OpenstackUpdate) ClearUser() *OpenstackUpdate {
	ou.mutation.ClearUser()
	return ou
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ou *OpenstackUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, ou.sqlSave, ou.mutation, ou.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ou *OpenstackUpdate) SaveX(ctx context.Context) int {
	affected, err := ou.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ou *OpenstackUpdate) Exec(ctx context.Context) error {
	_, err := ou.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ou *OpenstackUpdate) ExecX(ctx context.Context) {
	if err := ou.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ou *OpenstackUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(openstack.Table, openstack.Columns, sqlgraph.NewFieldSpec(openstack.FieldID, field.TypeInt))
	if ps := ou.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ou.mutation.Timestamp(); ok {
		_spec.SetField(openstack.FieldTimestamp, field.TypeTime, value)
	}
	if ou.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   openstack.UserTable,
			Columns: []string{openstack.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ou.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   openstack.UserTable,
			Columns: []string{openstack.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ou.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{openstack.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ou.mutation.done = true
	return n, nil
}

// OpenstackUpdateOne is the builder for updating a single Openstack entity.
type OpenstackUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *OpenstackMutation
}

// SetTimestamp sets the "timestamp" field.
func (ouo *OpenstackUpdateOne) SetTimestamp(t time.Time) *OpenstackUpdateOne {
	ouo.mutation.SetTimestamp(t)
	return ouo
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (ouo *OpenstackUpdateOne) SetNillableTimestamp(t *time.Time) *OpenstackUpdateOne {
	if t != nil {
		ouo.SetTimestamp(*t)
	}
	return ouo
}

// SetUserID sets the "user" edge to the User entity by ID.
func (ouo *OpenstackUpdateOne) SetUserID(id string) *OpenstackUpdateOne {
	ouo.mutation.SetUserID(id)
	return ouo
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (ouo *OpenstackUpdateOne) SetNillableUserID(id *string) *OpenstackUpdateOne {
	if id != nil {
		ouo = ouo.SetUserID(*id)
	}
	return ouo
}

// SetUser sets the "user" edge to the User entity.
func (ouo *OpenstackUpdateOne) SetUser(u *User) *OpenstackUpdateOne {
	return ouo.SetUserID(u.ID)
}

// Mutation returns the OpenstackMutation object of the builder.
func (ouo *OpenstackUpdateOne) Mutation() *OpenstackMutation {
	return ouo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (ouo *OpenstackUpdateOne) ClearUser() *OpenstackUpdateOne {
	ouo.mutation.ClearUser()
	return ouo
}

// Where appends a list predicates to the OpenstackUpdate builder.
func (ouo *OpenstackUpdateOne) Where(ps ...predicate.Openstack) *OpenstackUpdateOne {
	ouo.mutation.Where(ps...)
	return ouo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ouo *OpenstackUpdateOne) Select(field string, fields ...string) *OpenstackUpdateOne {
	ouo.fields = append([]string{field}, fields...)
	return ouo
}

// Save executes the query and returns the updated Openstack entity.
func (ouo *OpenstackUpdateOne) Save(ctx context.Context) (*Openstack, error) {
	return withHooks(ctx, ouo.sqlSave, ouo.mutation, ouo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ouo *OpenstackUpdateOne) SaveX(ctx context.Context) *Openstack {
	node, err := ouo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ouo *OpenstackUpdateOne) Exec(ctx context.Context) error {
	_, err := ouo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ouo *OpenstackUpdateOne) ExecX(ctx context.Context) {
	if err := ouo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ouo *OpenstackUpdateOne) sqlSave(ctx context.Context) (_node *Openstack, err error) {
	_spec := sqlgraph.NewUpdateSpec(openstack.Table, openstack.Columns, sqlgraph.NewFieldSpec(openstack.FieldID, field.TypeInt))
	id, ok := ouo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Openstack.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ouo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, openstack.FieldID)
		for _, f := range fields {
			if !openstack.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != openstack.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ouo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ouo.mutation.Timestamp(); ok {
		_spec.SetField(openstack.FieldTimestamp, field.TypeTime, value)
	}
	if ouo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   openstack.UserTable,
			Columns: []string{openstack.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ouo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   openstack.UserTable,
			Columns: []string{openstack.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Openstack{config: ouo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ouo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{openstack.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ouo.mutation.done = true
	return _node, nil
}
