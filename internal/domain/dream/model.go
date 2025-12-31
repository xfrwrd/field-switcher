package dream

import (
	"context"
	"fmt"
	"slices"

	domainErrors "github.com/xeniasokk/field-switcher/pkg/errors"
)

type Type string

const (
	TypeFootballer Type = "footballer"
)

type Role string

const (
	RoleTeamLead  Role = "team_lead"
	RoleDeveloper Role = "developer"
)

const (
	QualityPersistence = "Упорство"
)

type RoleConfig struct {
	title       string
	description string
	comment     string
	stack       []string
}

func (r RoleConfig) Title() string {
	return r.title
}

func (r RoleConfig) Description() string {
	return r.description
}

func (r RoleConfig) Comment() string {
	return r.comment
}

func (r RoleConfig) Stack() []string {
	return slices.Clone(r.stack)
}

type RoleOption func(*RoleConfig)

func WithTitle(title string) RoleOption {
	return func(cfg *RoleConfig) {
		cfg.title = title
	}
}

func WithDescription(description string) RoleOption {
	return func(cfg *RoleConfig) {
		cfg.description = description
	}
}

func WithComment(comment string) RoleOption {
	return func(cfg *RoleConfig) {
		cfg.comment = comment
	}
}

func WithStack(stack ...string) RoleOption {
	return func(cfg *RoleConfig) {
		cfg.stack = slices.Clone(stack)
	}
}

func NewRoleConfig(opts ...RoleOption) RoleConfig {
	cfg := RoleConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func (r Role) Config() (RoleConfig, error) {
	switch r {
	case RoleTeamLead:
		return NewRoleConfig(
			WithTitle("Тимлид"),
			WithDescription("Капитан команды на новом поле: вместо капитанской повязки — "+
				"ответственность за команду, вместо тактики на поле — архитектура и процессы. "+
				"Твоё упорство превратилось в настойчивость в решении сложных задач и поддержку команды."),
			WithComment("Ты не отказался от мечты — ты просто сменил поле и стал капитаном команды. "+
				"Твоё упорство привело тебя сюда."),
			WithStack(
				"System Design",
				"Team Leadership",
				"Agile/Scrum",
				"Code Review",
				"CI/CD",
				"Monitoring & Observability",
				"Technical Documentation",
			),
		), nil
	case RoleDeveloper:
		return NewRoleConfig(
			WithTitle("Разработчик"),
			WithDescription("Игрок на новом поле: вместо бутс — клавиатура, вместо газона — код. "+
				"Твоё упорство помогает преодолевать баги и дедлайны, как когда-то ты преодолевал защиту соперника."),
			WithComment("Ты не отказался от мечты — ты просто сменил поле. Ты всё ещё в игре. "+
				"Твоё упорство осталось с тобой."),
			WithStack("Go", "Git", "Microservices"),
		), nil
	default:
		return RoleConfig{}, domainErrors.NewDomainError(
			fmt.Sprintf("unsupported role: %s", r),
		)
	}
}

func (r Role) String() string {
	return string(r)
}

type Quality struct {
	name        string
	description string
}

func NewQuality(name, description string) (Quality, error) {
	if name == "" {
		return Quality{}, domainErrors.NewValidationError("quality name cannot be empty")
	}
	return Quality{name: name, description: description}, nil
}

func (q Quality) Name() string        { return q.name }
func (q Quality) Description() string { return q.description }

type Field struct {
	name        string
	environment string
}

func NewField(name, environment string) (Field, error) {
	if name == "" {
		return Field{}, domainErrors.NewValidationError("field name cannot be empty")
	}
	return Field{name: name, environment: environment}, nil
}

func (f Field) Name() string        { return f.name }
func (f Field) Environment() string { return f.environment }

type ChildhoodDream struct {
	dreamType     Type
	displayName   string
	desiredRole   string
	field         Field
	coreQualities []Quality
}

func NewChildhoodDream(
	dreamType Type,
	displayName string,
	desiredRole string,
	field Field,
	qualities []Quality,
) (ChildhoodDream, error) {
	if displayName == "" {
		return ChildhoodDream{}, domainErrors.NewValidationError("display name cannot be empty")
	}
	if desiredRole == "" {
		return ChildhoodDream{}, domainErrors.NewValidationError("desired role cannot be empty")
	}
	if len(qualities) == 0 {
		return ChildhoodDream{}, domainErrors.NewValidationError("qualities cannot be empty")
	}

	copied := slices.Clone(qualities)

	return ChildhoodDream{
		dreamType:     dreamType,
		displayName:   displayName,
		desiredRole:   desiredRole,
		field:         field,
		coreQualities: copied,
	}, nil
}

func (d ChildhoodDream) Type() Type           { return d.dreamType }
func (d ChildhoodDream) DisplayName() string  { return d.displayName }
func (d ChildhoodDream) DesiredRole() string  { return d.desiredRole }
func (d ChildhoodDream) Field() Field         { return d.field }
func (d ChildhoodDream) Qualities() []Quality { return slices.Clone(d.coreQualities) }

type Adult struct {
	roleTitle       string
	roleDescription string
	field           Field
	stack           []string
	traits          []Quality
	comment         string
}

func NewAdult(
	roleTitle string,
	roleDescription string,
	field Field,
	stack []string,
	traits []Quality,
	comment string,
) (Adult, error) {
	if roleTitle == "" {
		return Adult{}, domainErrors.NewValidationError("role title cannot be empty")
	}

	copiedStack := slices.Clone(stack)
	copiedTraits := slices.Clone(traits)

	return Adult{
		roleTitle:       roleTitle,
		roleDescription: roleDescription,
		field:           field,
		stack:           copiedStack,
		traits:          copiedTraits,
		comment:         comment,
	}, nil
}

func (a Adult) RoleTitle() string       { return a.roleTitle }
func (a Adult) RoleDescription() string { return a.roleDescription }
func (a Adult) Field() Field            { return a.field }
func (a Adult) Stack() []string         { return slices.Clone(a.stack) }
func (a Adult) Traits() []Quality       { return slices.Clone(a.traits) }
func (a Adult) Comment() string         { return a.comment }

const (
	DevFieldName        = "Поле разработки"
	DevFieldEnvironment = "Команда разработчиков, репозитории, прод-среда"
)

const (
	FootballFieldName        = "Футбольное поле"
	FootballFieldEnvironment = "Стадион, команда, трибуны"
	FootballerDisplayName    = "Футболист"
	FootballerDesiredRole    = "Полевой игрок"
)

const (
	QualityTeamSpirit        = "Командный дух"
	QualityTeamSpiritDesc    = "Играть ради общего результата"
	QualityPlayToWhistle     = "Игра до финального свистка"
	QualityPlayToWhistleDesc = "Не сдаваться до конца"
	QualityResilience        = "Умение держать удар"
	QualityResilienceDesc    = "Переживать промахи и критику"
	QualityGoalOriented      = "Стремление забивать"
	QualityGoalOrientedDesc  = "Ориентированность на результат"
	QualityPersistenceDesc   = "Не останавливаться перед препятствиями, продолжать путь несмотря ни на что"
)

func NewDevelopmentField() (Field, error) {
	return NewField(DevFieldName, DevFieldEnvironment)
}

func NewDefaultFootballerDream() (ChildhoodDream, error) {
	f, err := NewField(FootballFieldName, FootballFieldEnvironment)
	if err != nil {
		return ChildhoodDream{}, domainErrors.Wrap(err, domainErrors.CodeDomainFailure, "failed to create default field")
	}

	q1, err := NewQuality(QualityTeamSpirit, QualityTeamSpiritDesc)
	if err != nil {
		return ChildhoodDream{}, domainErrors.Wrap(
			err, domainErrors.CodeDomainFailure, fmt.Sprintf("failed to create quality: %s", QualityTeamSpirit),
		)
	}
	q2, err := NewQuality(QualityPlayToWhistle, QualityPlayToWhistleDesc)
	if err != nil {
		return ChildhoodDream{}, domainErrors.Wrap(
			err, domainErrors.CodeDomainFailure, fmt.Sprintf("failed to create quality: %s", QualityPlayToWhistle),
		)
	}
	q3, err := NewQuality(QualityResilience, QualityResilienceDesc)
	if err != nil {
		return ChildhoodDream{}, domainErrors.Wrap(
			err, domainErrors.CodeDomainFailure, fmt.Sprintf("failed to create quality: %s", QualityResilience),
		)
	}
	q4, err := NewQuality(QualityGoalOriented, QualityGoalOrientedDesc)
	if err != nil {
		return ChildhoodDream{}, domainErrors.Wrap(
			err, domainErrors.CodeDomainFailure, fmt.Sprintf("failed to create quality: %s", QualityGoalOriented),
		)
	}
	q5, err := NewQuality(QualityPersistence, QualityPersistenceDesc)
	if err != nil {
		return ChildhoodDream{}, domainErrors.Wrap(
			err, domainErrors.CodeDomainFailure, fmt.Sprintf("failed to create quality: %s", QualityPersistence),
		)
	}

	d, err := NewChildhoodDream(
		TypeFootballer,
		FootballerDisplayName,
		FootballerDesiredRole,
		f,
		[]Quality{q1, q2, q3, q4, q5},
	)
	if err != nil {
		return ChildhoodDream{}, domainErrors.Wrap(
			err, domainErrors.CodeDomainFailure, "failed to create default footballer dream",
		)
	}

	return d, nil
}

func (a Adult) Summary(ctx context.Context) string {
	_ = ctx
	return fmt.Sprintf(
		"%s — %s. Поле: %s (%s). Качества: %d",
		a.roleTitle,
		a.roleDescription,
		a.field.Name(),
		a.field.Environment(),
		len(a.traits),
	)
}
