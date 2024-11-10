// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ritsec/ops-bot-iii/ent/predicate"
	"github.com/ritsec/ops-bot-iii/ent/shitpost"
	"github.com/ritsec/ops-bot-iii/ent/user"
)

// ShitpostUpdate is the builder for updating Shitpost entities.
type ShitpostUpdate struct {
	config
	hooks    []Hook
	mutation *ShitpostMutation
}

// Where appends a list predicates to the ShitpostUpdate builder.
func (su *ShitpostUpdate) Where(ps ...predicate.Shitpost) *ShitpostUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetChannelID sets the "channel_id" field.
func (su *ShitpostUpdate) SetChannelID(s string) *ShitpostUpdate {
	su.mutation.SetChannelID(s)
	return su
}

// SetNillableChannelID sets the "channel_id" field if the given value is not nil.
func (su *ShitpostUpdate) SetNillableChannelID(s *string) *ShitpostUpdate {
	if s != nil {
		su.SetChannelID(*s)
	}
	return su
}

// SetCount sets the "count" field.
func (su *ShitpostUpdate) SetCount(i int) *ShitpostUpdate {
	su.mutation.ResetCount()
	su.mutation.SetCount(i)
	return su
}

// SetNillableCount sets the "count" field if the given value is not nil.
func (su *ShitpostUpdate) SetNillableCount(i *int) *ShitpostUpdate {
	if i != nil {
		su.SetCount(*i)
	}
	return su
}

// AddCount adds i to the "count" field.
func (su *ShitpostUpdate) AddCount(i int) *ShitpostUpdate {
	su.mutation.AddCount(i)
	return su
}

// SetUserID sets the "user" edge to the User entity by ID.
func (su *ShitpostUpdate) SetUserID(id string) *ShitpostUpdate {
	su.mutation.SetUserID(id)
	return su
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (su *ShitpostUpdate) SetNillableUserID(id *string) *ShitpostUpdate {
	if id != nil {
		su = su.SetUserID(*id)
	}
	return su
}

// SetUser sets the "user" edge to the User entity.
func (su *ShitpostUpdate) SetUser(u *User) *ShitpostUpdate {
	return su.SetUserID(u.ID)
}

// Mutation returns the ShitpostMutation object of the builder.
func (su *ShitpostUpdate) Mutation() *ShitpostMutation {
	return su.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (su *ShitpostUpdate) ClearUser() *ShitpostUpdate {
	su.mutation.ClearUser()
	return su
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *ShitpostUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, su.sqlSave, su.mutation, su.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (su *ShitpostUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *ShitpostUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *ShitpostUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (su *ShitpostUpdate) check() error {
	if v, ok := su.mutation.ChannelID(); ok {
		if err := shitpost.ChannelIDValidator(v); err != nil {
			return &ValidationError{Name: "channel_id", err: fmt.Errorf(`ent: validator failed for field "Shitpost.channel_id": %w`, err)}
		}
	}
	return nil
}

func (su *ShitpostUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := su.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(shitpost.Table, shitpost.Columns, sqlgraph.NewFieldSpec(shitpost.FieldID, field.TypeString))
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.ChannelID(); ok {
		_spec.SetField(shitpost.FieldChannelID, field.TypeString, value)
	}
	if value, ok := su.mutation.Count(); ok {
		_spec.SetField(shitpost.FieldCount, field.TypeInt, value)
	}
	if value, ok := su.mutation.AddedCount(); ok {
		_spec.AddField(shitpost.FieldCount, field.TypeInt, value)
	}
	if su.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   shitpost.UserTable,
			Columns: []string{shitpost.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   shitpost.UserTable,
			Columns: []string{shitpost.UserColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{shitpost.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	su.mutation.done = true
	return n, nil
}

// ShitpostUpdateOne is the builder for updating a single Shitpost entity.
type ShitpostUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ShitpostMutation
}

// SetChannelID sets the "channel_id" field.
func (suo *ShitpostUpdateOne) SetChannelID(s string) *ShitpostUpdateOne {
	suo.mutation.SetChannelID(s)
	return suo
}

// SetNillableChannelID sets the "channel_id" field if the given value is not nil.
func (suo *ShitpostUpdateOne) SetNillableChannelID(s *string) *ShitpostUpdateOne {
	if s != nil {
		suo.SetChannelID(*s)
	}
	return suo
}

// SetCount sets the "count" field.
func (suo *ShitpostUpdateOne) SetCount(i int) *ShitpostUpdateOne {
	suo.mutation.ResetCount()
	suo.mutation.SetCount(i)
	return suo
}

// SetNillableCount sets the "count" field if the given value is not nil.
func (suo *ShitpostUpdateOne) SetNillableCount(i *int) *ShitpostUpdateOne {
	if i != nil {
		suo.SetCount(*i)
	}
	return suo
}

// AddCount adds i to the "count" field.
func (suo *ShitpostUpdateOne) AddCount(i int) *ShitpostUpdateOne {
	suo.mutation.AddCount(i)
	return suo
}

// SetUserID sets the "user" edge to the User entity by ID.
func (suo *ShitpostUpdateOne) SetUserID(id string) *ShitpostUpdateOne {
	suo.mutation.SetUserID(id)
	return suo
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (suo *ShitpostUpdateOne) SetNillableUserID(id *string) *ShitpostUpdateOne {
	if id != nil {
		suo = suo.SetUserID(*id)
	}
	return suo
}

// SetUser sets the "user" edge to the User entity.
func (suo *ShitpostUpdateOne) SetUser(u *User) *ShitpostUpdateOne {
	return suo.SetUserID(u.ID)
}

// Mutation returns the ShitpostMutation object of the builder.
func (suo *ShitpostUpdateOne) Mutation() *ShitpostMutation {
	return suo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (suo *ShitpostUpdateOne) ClearUser() *ShitpostUpdateOne {
	suo.mutation.ClearUser()
	return suo
}

// Where appends a list predicates to the ShitpostUpdate builder.
func (suo *ShitpostUpdateOne) Where(ps ...predicate.Shitpost) *ShitpostUpdateOne {
	suo.mutation.Where(ps...)
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *ShitpostUpdateOne) Select(field string, fields ...string) *ShitpostUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Shitpost entity.
func (suo *ShitpostUpdateOne) Save(ctx context.Context) (*Shitpost, error) {
	return withHooks(ctx, suo.sqlSave, suo.mutation, suo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (suo *ShitpostUpdateOne) SaveX(ctx context.Context) *Shitpost {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *ShitpostUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *ShitpostUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (suo *ShitpostUpdateOne) check() error {
	if v, ok := suo.mutation.ChannelID(); ok {
		if err := shitpost.ChannelIDValidator(v); err != nil {
			return &ValidationError{Name: "channel_id", err: fmt.Errorf(`ent: validator failed for field "Shitpost.channel_id": %w`, err)}
		}
	}
	return nil
}

func (suo *ShitpostUpdateOne) sqlSave(ctx context.Context) (_node *Shitpost, err error) {
	if err := suo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(shitpost.Table, shitpost.Columns, sqlgraph.NewFieldSpec(shitpost.FieldID, field.TypeString))
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Shitpost.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, shitpost.FieldID)
		for _, f := range fields {
			if !shitpost.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != shitpost.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := suo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := suo.mutation.ChannelID(); ok {
		_spec.SetField(shitpost.FieldChannelID, field.TypeString, value)
	}
	if value, ok := suo.mutation.Count(); ok {
		_spec.SetField(shitpost.FieldCount, field.TypeInt, value)
	}
	if value, ok := suo.mutation.AddedCount(); ok {
		_spec.AddField(shitpost.FieldCount, field.TypeInt, value)
	}
	if suo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   shitpost.UserTable,
			Columns: []string{shitpost.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   shitpost.UserTable,
			Columns: []string{shitpost.UserColumn},
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
	_node = &Shitpost{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{shitpost.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	suo.mutation.done = true
	return _node, nil
}
