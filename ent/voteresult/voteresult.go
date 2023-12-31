// Code generated by ent, DO NOT EDIT.

package voteresult

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the voteresult type in the database.
	Label = "vote_result"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldHTML holds the string denoting the html field in the database.
	FieldHTML = "html"
	// FieldPlain holds the string denoting the plain field in the database.
	FieldPlain = "plain"
	// FieldVoteID holds the string denoting the vote_id field in the database.
	FieldVoteID = "vote_id"
	// Table holds the table name of the voteresult in the database.
	Table = "vote_results"
)

// Columns holds all SQL columns for voteresult fields.
var Columns = []string{
	FieldID,
	FieldHTML,
	FieldPlain,
	FieldVoteID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// VoteIDValidator is a validator for the "vote_id" field. It is called by the builders before save.
	VoteIDValidator func(string) error
)

// OrderOption defines the ordering options for the VoteResult queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByHTML orders the results by the html field.
func ByHTML(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHTML, opts...).ToFunc()
}

// ByPlain orders the results by the plain field.
func ByPlain(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPlain, opts...).ToFunc()
}

// ByVoteID orders the results by the vote_id field.
func ByVoteID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVoteID, opts...).ToFunc()
}
