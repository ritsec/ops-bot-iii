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
	<!DOCTYPE html>
	<html>
	<head>
		<title>%s</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				margin: 0;
				padding: 0;
				background-image: url("https://ritsec.club/assets/images/bg-dark.svg");
				text-align: center;
			}
			.header {
				background-color: #0d0d0f;
				color: white;
				padding: 10px 0;
				margin-bottom: 20px;
			}
			.container {
				width: 1000px;
				margin: 20px auto;
				margin-top: 150px;
				padding: 30px;
				color: white;
				background-color: #222222;
				box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
			}
		</style>
	</head>
	<body>

		<div class="header">
			<h1>%s</h1>
		</div>

		<div class="container">
			<h1>Winner: %s</h1>
			<div id="sankey_multiple"></div>
		</div>

		<script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
		<script type="text/javascript">
			google.charts.load("current", {packages:["sankey"]});
			google.charts.setOnLoadCallback(drawChart);
			function drawChart() {
				var data = new google.visualization.DataTable();
				data.addColumn('string', 'From');
				data.addColumn('string', 'To');
				data.addColumn('number', 'Weight');
				data.addRows([
					%s
				]);

				var options = {
					width: 1000,
					sankey: {
						node: {
							colors: ['#ff7613', '#ff7613', '#ff7613', '#ff7613', '#ff7613', '#ff7613', '#ff7613', '#ff7613', '#ff7613', '#ff7613', '#ff7613', '#ff7613'],

							label: {
								color: '#ffffff',
								fontSize: 14,
							},
						},
						link: {
							colors: ['#ff7613'],

						},
					},
				};

				var chart = new google.visualization.Sankey(document.getElementById('sankey_multiple'));
				chart.draw(data, options);
			}
		</script>

	</body>
	</html>`

	var rows SanKeyRows

	for i, round := range r.Rounds {
		if i == len(r.Rounds)-1 {
			continue
		}
		for _, choice := range round.Choices {
			if choice.Selection == r.Eliminations[i] {
				continue
			}

			rows.Rows = append(rows.Rows, SanKeyRow{
				From:   fmt.Sprintf("%v - Round %d", choice.Selection, i),
				To:     fmt.Sprintf("%v - Round %d", choice.Selection, i+1),
				Weight: choice.Votes,
			})
		}

		rows.Rows = append(rows.Rows, createFlows(r.Rounds[i], r.Rounds[i+1], r.Eliminations[i], i)...)
	}

	return fmt.Sprintf(HTML, r.Title, r.Title, r.Winner, rows.String())
}

// createFlows creates the flows between rounds
func createFlows(fromRound Round, toRound Round, eliminated string, round int) []SanKeyRow {
	var rows []SanKeyRow

	for _, from := range fromRound.Choices {
		for _, to := range toRound.Choices {
			if from.Selection == to.Selection {
				if from.Votes < to.Votes {
					rows = append(rows, SanKeyRow{
						From:   fmt.Sprintf("%v - Round %d", eliminated, round),
						To:     fmt.Sprintf("%v - Round %d", to.Selection, round+1),
						Weight: to.Votes - from.Votes,
					})
				}
			}
		}
	}

	return rows
}
