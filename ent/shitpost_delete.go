// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ritsec/ops-bot-iii/ent/predicate"
	"github.com/ritsec/ops-bot-iii/ent/shitpost"
)

// ShitpostDelete is the builder for deleting a Shitpost entity.
type ShitpostDelete struct {
	config
	hooks    []Hook
	mutation *ShitpostMutation
}

// Where appends a list predicates to the ShitpostDelete builder.
func (sd *ShitpostDelete) Where(ps ...predicate.Shitpost) *ShitpostDelete {
	sd.mutation.Where(ps...)
	return sd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (sd *ShitpostDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, ShitpostMutation](ctx, sd.sqlExec, sd.mutation, sd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (sd *ShitpostDelete) ExecX(ctx context.Context) int {
	n, err := sd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (sd *ShitpostDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(shitpost.Table, sqlgraph.NewFieldSpec(shitpost.FieldID, field.TypeString))
	if ps := sd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, sd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	sd.mutation.done = true
	return affected, err
}

// ShitpostDeleteOne is the builder for deleting a single Shitpost entity.
type ShitpostDeleteOne struct {
	sd *ShitpostDelete
}

// Where appends a list predicates to the ShitpostDelete builder.
func (sdo *ShitpostDeleteOne) Where(ps ...predicate.Shitpost) *ShitpostDeleteOne {
	sdo.sd.mutation.Where(ps...)
	return sdo
}

// Exec executes the deletion query.
func (sdo *ShitpostDeleteOne) Exec(ctx context.Context) error {
	n, err := sdo.sd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{shitpost.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (sdo *ShitpostDeleteOne) ExecX(ctx context.Context) {
	if err := sdo.Exec(ctx); err != nil {
		panic(err)
	}
}