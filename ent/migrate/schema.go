// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// BirthdaysColumns holds the columns for the "birthdays" table.
	BirthdaysColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "day", Type: field.TypeInt},
		{Name: "month", Type: field.TypeInt},
		{Name: "user_birthday", Type: field.TypeString, Unique: true, Nullable: true},
	}
	// BirthdaysTable holds the schema information for the "birthdays" table.
	BirthdaysTable = &schema.Table{
		Name:       "birthdays",
		Columns:    BirthdaysColumns,
		PrimaryKey: []*schema.Column{BirthdaysColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "birthdays_users_birthday",
				Columns:    []*schema.Column{BirthdaysColumns[3]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// ShitpostsColumns holds the columns for the "shitposts" table.
	ShitpostsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "channel_id", Type: field.TypeString},
		{Name: "count", Type: field.TypeInt},
		{Name: "user_shitposts", Type: field.TypeString, Nullable: true},
	}
	// ShitpostsTable holds the schema information for the "shitposts" table.
	ShitpostsTable = &schema.Table{
		Name:       "shitposts",
		Columns:    ShitpostsColumns,
		PrimaryKey: []*schema.Column{ShitpostsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "shitposts_users_shitposts",
				Columns:    []*schema.Column{ShitpostsColumns[3]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// SigninsColumns holds the columns for the "signins" table.
	SigninsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "timestamp", Type: field.TypeTime},
		{Name: "type", Type: field.TypeEnum, Enums: []string{"General Meeting", "Red Team", "Red Team Recruiting", "Reversing", "RVAPT", "Contagion", "Physical", "Wireless", "IR", "WiCyS", "Ops", "Ops IG", "Vulnerability Research", "Mentorship", "Other"}},
		{Name: "user_signins", Type: field.TypeString},
	}
	// SigninsTable holds the schema information for the "signins" table.
	SigninsTable = &schema.Table{
		Name:       "signins",
		Columns:    SigninsColumns,
		PrimaryKey: []*schema.Column{SigninsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "signins_users_signins",
				Columns:    []*schema.Column{SigninsColumns[3]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "email", Type: field.TypeString, Default: ""},
		{Name: "verification_attempts", Type: field.TypeInt8, Default: 0},
		{Name: "verified", Type: field.TypeBool, Default: false},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// VotesColumns holds the columns for the "votes" table.
	VotesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "selection", Type: field.TypeString},
		{Name: "rank", Type: field.TypeInt},
		{Name: "vote_id", Type: field.TypeString},
		{Name: "user_votes", Type: field.TypeString},
	}
	// VotesTable holds the schema information for the "votes" table.
	VotesTable = &schema.Table{
		Name:       "votes",
		Columns:    VotesColumns,
		PrimaryKey: []*schema.Column{VotesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "votes_users_votes",
				Columns:    []*schema.Column{VotesColumns[4]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// VoteResultsColumns holds the columns for the "vote_results" table.
	VoteResultsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "html", Type: field.TypeString},
		{Name: "plain", Type: field.TypeString},
		{Name: "vote_id", Type: field.TypeString, Unique: true},
	}
	// VoteResultsTable holds the schema information for the "vote_results" table.
	VoteResultsTable = &schema.Table{
		Name:       "vote_results",
		Columns:    VoteResultsColumns,
		PrimaryKey: []*schema.Column{VoteResultsColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		BirthdaysTable,
		ShitpostsTable,
		SigninsTable,
		UsersTable,
		VotesTable,
		VoteResultsTable,
	}
)

func init() {
	BirthdaysTable.ForeignKeys[0].RefTable = UsersTable
	ShitpostsTable.ForeignKeys[0].RefTable = UsersTable
	SigninsTable.ForeignKeys[0].RefTable = UsersTable
	VotesTable.ForeignKeys[0].RefTable = UsersTable
}
