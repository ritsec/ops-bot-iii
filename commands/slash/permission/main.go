package permission

import "github.com/bwmarrin/discordgo"

var (
	// Member Permissions
	Member int64 = discordgo.PermissionChangeNickname

	// Admin Permissions
	Admin int64 = discordgo.PermissionManageChannels

	// IG Lead Permissions
	IGLead int64 = discordgo.PermissionVoiceMoveMembers
)
