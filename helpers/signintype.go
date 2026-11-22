package helpers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/ent/signin"
)

// Converts signin type string to signin type
func StringToType(signinString string) signin.Type {
	var entSigninType signin.Type
	switch signinString {
	case "General Meeting":
		entSigninType = signin.TypeGeneralMeeting
	case "Contagion":
		entSigninType = signin.TypeContagion
	case "CCDC":
		entSigninType = signin.TypeCCDC
	case "CCDC Tryouts":
		entSigninType = signin.TypeCCDCTryouts
	case "IR":
		entSigninType = signin.TypeIR
	case "Ops":
		entSigninType = signin.TypeOps
	case "Ops IG":
		entSigninType = signin.TypeOpsIG
	case "Red Team":
		entSigninType = signin.TypeRedTeam
	case "Red Team Recruiting":
		entSigninType = signin.TypeRedTeamRecruiting
	case "RVAPT":
		entSigninType = signin.TypeRVAPT
	case "Reversing":
		entSigninType = signin.TypeReversing
	case "Physical":
		entSigninType = signin.TypePhysical
	case "Wireless":
		entSigninType = signin.TypeWireless
	case "WiCyS":
		entSigninType = signin.TypeWiCyS
	case "Vulnerability Research":
		entSigninType = signin.TypeVulnerabilityResearch
	case "Mentorship":
		entSigninType = signin.TypeMentorship
	case "Zero To Hero":
		entSigninType = signin.TypeZeroToHero
	case "OT Security":
		entSigninType = signin.TypeOTSecurity
	case "Other":
		entSigninType = signin.TypeOther
	case "All":
		entSigninType = "All"
	}
	return entSigninType
}

// Return array of all signin types (update return array length)
func SigninTypeArray() [19]string {
	return [...]string{
		"General Meeting",
		"Contagion",
		"CCDC",
		"CCDC Tryouts",
		"IR",
		"Zero To Hero",
		"OT Security",
		"Ops",
		"Ops IG",
		"Red Team",
		"Red Team Recruiting",
		"RVAPT",
		"Reversing",
		"Physical",
		"Wireless",
		"WiCyS",
		"Vulnerability Research",
		"Mentorship",
		"Other",
	}

}

// SigninTypeChoices returns the dropdown choices for all signin types, derived from SigninTypeArray().
func SigninTypeChoices() []*discordgo.ApplicationCommandOptionChoice {
	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(SigninTypeArray()))
	for _, t := range SigninTypeArray() {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{Name: t, Value: t})
	}
	return choices
}
