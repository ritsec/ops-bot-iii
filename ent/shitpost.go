// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/ritsec/ops-bot-iii/ent/shitpost"
	"github.com/ritsec/ops-bot-iii/ent/user"
)

// Shitpost is the model entity for the Shitpost schema.
type Shitpost struct {
	config `json:"-"`
	// ID of the ent.
	// Message ID
	ID string `json:"id,omitempty"`
	// Shitpost Count
	Count int `json:"count,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ShitpostQuery when eager-loading is set.
	Edges          ShitpostEdges `json:"edges"`
	user_shitposts *string
	selectValues   sql.SelectValues
}

// ShitpostEdges holds the relations/edges for other nodes in the graph.
type ShitpostEdges struct {
	// Shitpost Author
	User *User `json:"user,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ShitpostEdges) UserOrErr() (*User, error) {
	if e.loadedTypes[0] {
		if e.User == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.User, nil
	}
	return nil, &NotLoadedError{edge: "user"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Shitpost) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case shitpost.FieldCount:
			values[i] = new(sql.NullInt64)
		case shitpost.FieldID:
			values[i] = new(sql.NullString)
		case shitpost.ForeignKeys[0]: // user_shitposts
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Shitpost fields.
func (s *Shitpost) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case shitpost.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				s.ID = value.String
			}
		case shitpost.FieldCount:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field count", values[i])
			} else if value.Valid {
				s.Count = int(value.Int64)
			}
		case shitpost.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field user_shitposts", values[i])
			} else if value.Valid {
				s.user_shitposts = new(string)
				*s.user_shitposts = value.String
			}
		default:
			s.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Shitpost.
// This includes values selected through modifiers, order, etc.
func (s *Shitpost) Value(name string) (ent.Value, error) {
	return s.selectValues.Get(name)
}

// QueryUser queries the "user" edge of the Shitpost entity.
func (s *Shitpost) QueryUser() *UserQuery {
	return NewShitpostClient(s.config).QueryUser(s)
}

// Update returns a builder for updating this Shitpost.
// Note that you need to call Shitpost.Unwrap() before calling this method if this Shitpost
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Shitpost) Update() *ShitpostUpdateOne {
	return NewShitpostClient(s.config).UpdateOne(s)
}

// Unwrap unwraps the Shitpost entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (s *Shitpost) Unwrap() *Shitpost {
	_tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("ent: Shitpost is not a transactional entity")
	}
	s.config.driver = _tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *Shitpost) String() string {
	var builder strings.Builder
	builder.WriteString("Shitpost(")
	builder.WriteString(fmt.Sprintf("id=%v, ", s.ID))
	builder.WriteString("count=")
	builder.WriteString(fmt.Sprintf("%v", s.Count))
	builder.WriteByte(')')
	return builder.String()
}

// Shitposts is a parsable slice of Shitpost.
type Shitposts []*Shitpost