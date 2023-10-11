// Code generated by ent, DO NOT EDIT.

package vote

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/ritsec/ops-bot-iii/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Vote {
	return predicate.Vote(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Vote {
	return predicate.Vote(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Vote {
	return predicate.Vote(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Vote {
	return predicate.Vote(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Vote {
	return predicate.Vote(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Vote {
	return predicate.Vote(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Vote {
	return predicate.Vote(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Vote {
	return predicate.Vote(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Vote {
	return predicate.Vote(sql.FieldLTE(FieldID, id))
}

// Selection applies equality check predicate on the "selection" field. It's identical to SelectionEQ.
func Selection(v string) predicate.Vote {
	return predicate.Vote(sql.FieldEQ(FieldSelection, v))
}

// Rank applies equality check predicate on the "rank" field. It's identical to RankEQ.
func Rank(v int) predicate.Vote {
	return predicate.Vote(sql.FieldEQ(FieldRank, v))
}

// VoteID applies equality check predicate on the "vote_id" field. It's identical to VoteIDEQ.
func VoteID(v string) predicate.Vote {
	return predicate.Vote(sql.FieldEQ(FieldVoteID, v))
}

// SelectionEQ applies the EQ predicate on the "selection" field.
func SelectionEQ(v string) predicate.Vote {
	return predicate.Vote(sql.FieldEQ(FieldSelection, v))
}

// SelectionNEQ applies the NEQ predicate on the "selection" field.
func SelectionNEQ(v string) predicate.Vote {
	return predicate.Vote(sql.FieldNEQ(FieldSelection, v))
}

// SelectionIn applies the In predicate on the "selection" field.
func SelectionIn(vs ...string) predicate.Vote {
	return predicate.Vote(sql.FieldIn(FieldSelection, vs...))
}

// SelectionNotIn applies the NotIn predicate on the "selection" field.
func SelectionNotIn(vs ...string) predicate.Vote {
	return predicate.Vote(sql.FieldNotIn(FieldSelection, vs...))
}

// SelectionGT applies the GT predicate on the "selection" field.
func SelectionGT(v string) predicate.Vote {
	return predicate.Vote(sql.FieldGT(FieldSelection, v))
}

// SelectionGTE applies the GTE predicate on the "selection" field.
func SelectionGTE(v string) predicate.Vote {
	return predicate.Vote(sql.FieldGTE(FieldSelection, v))
}

// SelectionLT applies the LT predicate on the "selection" field.
func SelectionLT(v string) predicate.Vote {
	return predicate.Vote(sql.FieldLT(FieldSelection, v))
}

// SelectionLTE applies the LTE predicate on the "selection" field.
func SelectionLTE(v string) predicate.Vote {
	return predicate.Vote(sql.FieldLTE(FieldSelection, v))
}

// SelectionContains applies the Contains predicate on the "selection" field.
func SelectionContains(v string) predicate.Vote {
	return predicate.Vote(sql.FieldContains(FieldSelection, v))
}

// SelectionHasPrefix applies the HasPrefix predicate on the "selection" field.
func SelectionHasPrefix(v string) predicate.Vote {
	return predicate.Vote(sql.FieldHasPrefix(FieldSelection, v))
}

// SelectionHasSuffix applies the HasSuffix predicate on the "selection" field.
func SelectionHasSuffix(v string) predicate.Vote {
	return predicate.Vote(sql.FieldHasSuffix(FieldSelection, v))
}

// SelectionEqualFold applies the EqualFold predicate on the "selection" field.
func SelectionEqualFold(v string) predicate.Vote {
	return predicate.Vote(sql.FieldEqualFold(FieldSelection, v))
}

// SelectionContainsFold applies the ContainsFold predicate on the "selection" field.
func SelectionContainsFold(v string) predicate.Vote {
	return predicate.Vote(sql.FieldContainsFold(FieldSelection, v))
}

// RankEQ applies the EQ predicate on the "rank" field.
func RankEQ(v int) predicate.Vote {
	return predicate.Vote(sql.FieldEQ(FieldRank, v))
}

// RankNEQ applies the NEQ predicate on the "rank" field.
func RankNEQ(v int) predicate.Vote {
	return predicate.Vote(sql.FieldNEQ(FieldRank, v))
}

// RankIn applies the In predicate on the "rank" field.
func RankIn(vs ...int) predicate.Vote {
	return predicate.Vote(sql.FieldIn(FieldRank, vs...))
}

// RankNotIn applies the NotIn predicate on the "rank" field.
func RankNotIn(vs ...int) predicate.Vote {
	return predicate.Vote(sql.FieldNotIn(FieldRank, vs...))
}

// RankGT applies the GT predicate on the "rank" field.
func RankGT(v int) predicate.Vote {
	return predicate.Vote(sql.FieldGT(FieldRank, v))
}

// RankGTE applies the GTE predicate on the "rank" field.
func RankGTE(v int) predicate.Vote {
	return predicate.Vote(sql.FieldGTE(FieldRank, v))
}

// RankLT applies the LT predicate on the "rank" field.
func RankLT(v int) predicate.Vote {
	return predicate.Vote(sql.FieldLT(FieldRank, v))
}

// RankLTE applies the LTE predicate on the "rank" field.
func RankLTE(v int) predicate.Vote {
	return predicate.Vote(sql.FieldLTE(FieldRank, v))
}

// VoteIDEQ applies the EQ predicate on the "vote_id" field.
func VoteIDEQ(v string) predicate.Vote {
	return predicate.Vote(sql.FieldEQ(FieldVoteID, v))
}

// VoteIDNEQ applies the NEQ predicate on the "vote_id" field.
func VoteIDNEQ(v string) predicate.Vote {
	return predicate.Vote(sql.FieldNEQ(FieldVoteID, v))
}

// VoteIDIn applies the In predicate on the "vote_id" field.
func VoteIDIn(vs ...string) predicate.Vote {
	return predicate.Vote(sql.FieldIn(FieldVoteID, vs...))
}

// VoteIDNotIn applies the NotIn predicate on the "vote_id" field.
func VoteIDNotIn(vs ...string) predicate.Vote {
	return predicate.Vote(sql.FieldNotIn(FieldVoteID, vs...))
}

// VoteIDGT applies the GT predicate on the "vote_id" field.
func VoteIDGT(v string) predicate.Vote {
	return predicate.Vote(sql.FieldGT(FieldVoteID, v))
}

// VoteIDGTE applies the GTE predicate on the "vote_id" field.
func VoteIDGTE(v string) predicate.Vote {
	return predicate.Vote(sql.FieldGTE(FieldVoteID, v))
}

// VoteIDLT applies the LT predicate on the "vote_id" field.
func VoteIDLT(v string) predicate.Vote {
	return predicate.Vote(sql.FieldLT(FieldVoteID, v))
}

// VoteIDLTE applies the LTE predicate on the "vote_id" field.
func VoteIDLTE(v string) predicate.Vote {
	return predicate.Vote(sql.FieldLTE(FieldVoteID, v))
}

// VoteIDContains applies the Contains predicate on the "vote_id" field.
func VoteIDContains(v string) predicate.Vote {
	return predicate.Vote(sql.FieldContains(FieldVoteID, v))
}

// VoteIDHasPrefix applies the HasPrefix predicate on the "vote_id" field.
func VoteIDHasPrefix(v string) predicate.Vote {
	return predicate.Vote(sql.FieldHasPrefix(FieldVoteID, v))
}

// VoteIDHasSuffix applies the HasSuffix predicate on the "vote_id" field.
func VoteIDHasSuffix(v string) predicate.Vote {
	return predicate.Vote(sql.FieldHasSuffix(FieldVoteID, v))
}

// VoteIDEqualFold applies the EqualFold predicate on the "vote_id" field.
func VoteIDEqualFold(v string) predicate.Vote {
	return predicate.Vote(sql.FieldEqualFold(FieldVoteID, v))
}

// VoteIDContainsFold applies the ContainsFold predicate on the "vote_id" field.
func VoteIDContainsFold(v string) predicate.Vote {
	return predicate.Vote(sql.FieldContainsFold(FieldVoteID, v))
}

// HasUser applies the HasEdge predicate on the "user" edge.
func HasUser() predicate.Vote {
	return predicate.Vote(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserWith applies the HasEdge predicate on the "user" edge with a given conditions (other predicates).
func HasUserWith(preds ...predicate.User) predicate.Vote {
	return predicate.Vote(func(s *sql.Selector) {
		step := newUserStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Vote) predicate.Vote {
	return predicate.Vote(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Vote) predicate.Vote {
	return predicate.Vote(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Vote) predicate.Vote {
	return predicate.Vote(func(s *sql.Selector) {
		p(s.Not())
	})
}
