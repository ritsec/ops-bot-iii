package helpers

import "github.com/bwmarrin/discordgo"

func Username(s *discordgo.Session, userID string) (string, error) {
	user, err := s.User(userID)
	if err != nil {
		return "", err
	}

	return user.Username, nil
}
