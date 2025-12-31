package app

import "github.com/xeniasokk/field-switcher/internal/domain/dream"

type Config struct {
	TargetRole dream.Role
}

func DefaultConfig() Config {
	return Config{
		TargetRole: dream.RoleTeamLead,
	}
}
