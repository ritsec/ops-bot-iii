package data

import "testing"

func TestEmailLocalPart(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  string
	}{
		{"standard RIT email", "abc1234@rit.edu", "abc1234"},
		{"g prefixed RIT email", "g.abc1234@g.rit.edu", "abc1234"},
		{"uppercase g prefixed email", "G.ABC1234@G.RIT.EDU", "abc1234"},
		{"gmail email", "jdoe@gmail.com", "jdoe"},
		{"tag subdomain email", "user.name+tag@sub.domain.co.uk", "user.name+tag"},
		{"email without at sign", "noatsign", "noatsign"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := emailLocalPart(tt.email); got != tt.want {
				t.Errorf("emailLocalPart(%q) = %q, want %q", tt.email, got, tt.want)
			}
		})
	}
}

func TestNormalizeRITEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  string
	}{
		{"uppercase g rit email", "ABC1234@G.RIT.EDU", "abc1234@rit.edu"},
		{"standard rit email", "abc1234@rit.edu", "abc1234@rit.edu"},
		{"gmail email", "jdoe@gmail.com", "jdoe@gmail.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeRITEmail(tt.email); got != tt.want {
				t.Errorf("normalizeRITEmail(%q) = %q, want %q", tt.email, got, tt.want)
			}
		})
	}
}
