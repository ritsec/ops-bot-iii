// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/ritsec/ops-bot-iii/ent/voteresult"
)

// VoteResult is the model entity for the VoteResult schema.
type VoteResult struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// The vote's HTML Results
	HTML string `json:"html,omitempty"`
	// The vote's plaintext results
	Plain string `json:"plain,omitempty"`
	// The vote's ID
	VoteID       string `json:"vote_id,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*VoteResult) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case voteresult.FieldID:
			values[i] = new(sql.NullInt64)
		case voteresult.FieldHTML, voteresult.FieldPlain, voteresult.FieldVoteID:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the VoteResult fields.
func (vr *VoteResult) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case voteresult.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			vr.ID = int(value.Int64)
		case voteresult.FieldHTML:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field html", values[i])
			} else if value.Valid {
				vr.HTML = value.String
			}
		case voteresult.FieldPlain:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field plain", values[i])
			} else if value.Valid {
				vr.Plain = value.String
			}
		case voteresult.FieldVoteID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field vote_id", values[i])
			} else if value.Valid {
				vr.VoteID = value.String
			}
		default:
			vr.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the VoteResult.
// This includes values selected through modifiers, order, etc.
func (vr *VoteResult) Value(name string) (ent.Value, error) {
	return vr.selectValues.Get(name)
}

// Update returns a builder for updating this VoteResult.
// Note that you need to call VoteResult.Unwrap() before calling this method if this VoteResult
// was returned from a transaction, and the transaction was committed or rolled back.
func (vr *VoteResult) Update() *VoteResultUpdateOne {
	return NewVoteResultClient(vr.config).UpdateOne(vr)
}

// Unwrap unwraps the VoteResult entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (vr *VoteResult) Unwrap() *VoteResult {
	_tx, ok := vr.config.driver.(*txDriver)
	if !ok {
		panic("ent: VoteResult is not a transactional entity")
	}
	vr.config.driver = _tx.drv
	return vr
}

// String implements the fmt.Stringer.
func (vr *VoteResult) String() string {
	var builder strings.Builder
	builder.WriteString("VoteResult(")
	builder.WriteString(fmt.Sprintf("id=%v, ", vr.ID))
	builder.WriteString("html=")
	builder.WriteString(vr.HTML)
	builder.WriteString(", ")
	builder.WriteString("plain=")
	builder.WriteString(vr.Plain)
	builder.WriteString(", ")
	builder.WriteString("vote_id=")
	builder.WriteString(vr.VoteID)
	builder.WriteByte(')')
	return builder.String()
}

// VoteResults is a parsable slice of VoteResult.
type VoteResults []*VoteResult
