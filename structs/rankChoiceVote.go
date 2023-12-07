package structs

import (
	"fmt"
)

// Choice is a choice in a rank choice vote
type Choice struct {
	// Selection is the choice
	Selection string

	// Votes is the number of votes for the choice
	Votes int
}

// Round is a round of voting
type Round struct {
	// Choices is the choices in the round
	Choices []Choice
}

// RankChoiceVote is a rank choice vote
type RankChoiceVote struct {
	// Title is the title of the vote
	Title string

	// Options is the options in the vote
	Options []string

	// Winner is the winner of the vote
	Winner string

	// Rounds is the rounds of the vote
	Rounds []Round

	// Eliminations is the eliminations of the vote
	Eliminations []string
}

// SanKeyRow is a row in a sankey diagram
type SanKeyRow struct {
	// From is the from node
	From string

	// To is the to node
	To string

	// Weight is the weight of the node
	Weight int
}

// SanKeyRows is a list of rows in a sankey diagram
type SanKeyRows struct {
	// Rows is the rows in the sankey diagram
	Rows []SanKeyRow
}

// String returns the string representation of a SanKeyRow
func (s SanKeyRow) String() string {
	return fmt.Sprintf("['%s', '%s', %d]", s.From, s.To, s.Weight)
}

// String returns the string representation of a SanKeyRows
func (s SanKeyRows) String() string {
	output := "["
	for i, row := range s.Rows {
		output += row.String()
		if i != len(s.Rows)-1 {
			output += ", "
		}
	}
	output += "]"

	return output
}

// ConvertToSanKeyRows converts a map of rounds to a SanKeyRows
func ConvertToRound(round map[string]int, options []string) Round {
	var r Round

	for _, option := range options {
		r.Choices = append(r.Choices, Choice{
			Selection: option,
			Votes:     round[option],
		})
	}

	return r
}

// ConvertToSanKeyRows converts a map of rounds to a SanKeyRows
func (r RankChoiceVote) String() string {
	response := r.Title + "\n\n"

	for i, round := range r.Rounds {
		response += fmt.Sprintf("=== Round %d ===\n", i+1)
		for _, choice := range round.Choices {
			response += fmt.Sprintf("%s: %d\n", choice.Selection, choice.Votes)
		}
		response += "\n"
	}

	response += fmt.Sprintf("\nWinner: %s", r.Winner)

	return response
}

// ConvertToSanKeyRows converts a map of rounds to a SanKeyRowss
func (r RankChoiceVote) HTML() string {
	HTML := `
	<html>
	<body>
		<script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>

	<div id="sankey_multiple" style="width: 900px; height: 300px;"></div>

	<script type="text/javascript">
		google.charts.load("current", {packages:["sankey"]});
		google.charts.setOnLoadCallback(drawChart);
		function drawChart() {
		var data = new google.visualization.DataTable();
		data.addColumn('string', 'From');
		data.addColumn('string', 'To');
		data.addColumn('number', 'Weight');
		data.addRows(%v);

		// Set chart options
		var options = {
			width: 600,
		};

		// Instantiate and draw our chart, passing in some options.
		var chart = new google.visualization.Sankey(document.getElementById('sankey_multiple'));
		chart.draw(data, options);
		}
	</script>
	</body>
	</html>`

	var rows SanKeyRows

	for i, round := range r.Rounds {
		if i == 0 {
			continue
		}
		for _, choice := range round.Choices {
			if choice.Selection == r.Eliminations[i-1] {
				fromRound := Round{}
				toRound := Round{}

				fromRound.Choices = make([]Choice, len(r.Rounds[i-1].Choices))
				copy(fromRound.Choices, r.Rounds[i-1].Choices)

				toRound.Choices = make([]Choice, len(r.Rounds[i].Choices))
				copy(toRound.Choices, r.Rounds[i].Choices)

				rows.Rows = append(rows.Rows, createFlows(fromRound, toRound, choice.Selection, i)...)
			} else {
				rows.Rows = append(rows.Rows, SanKeyRow{
					From:   fmt.Sprintf("%v - Round %d", choice.Selection, i-1),
					To:     fmt.Sprintf("%v - Round %d", choice.Selection, i),
					Weight: choice.Votes,
				})
			}
		}
	}

	return fmt.Sprintf(HTML, rows.String())
}

// createFlows creates the flows between rounds
func createFlows(fromRound Round, toRound Round, eliminated string, round int) []SanKeyRow {
	var rows []SanKeyRow

	for _, from := range fromRound.Choices {
		for _, to := range toRound.Choices {
			if from.Selection != to.Selection {
				if from.Votes < to.Votes {
					rows = append(rows, SanKeyRow{
						From:   fmt.Sprintf("%v - Round %d", eliminated, round-1),
						To:     fmt.Sprintf("%v - Round %d", to.Selection, round),
						Weight: to.Votes - from.Votes,
					})
				}
			}
		}
	}

	return rows
}
