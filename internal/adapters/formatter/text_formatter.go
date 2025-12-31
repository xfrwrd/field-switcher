package formatter

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/xeniasokk/field-switcher/internal/domain/dream"
	"github.com/xeniasokk/field-switcher/internal/ports"
)

var _ ports.FormatterPort = (*TextFormatter)(nil)

type TextFormatter struct{}

func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}

func (f *TextFormatter) Format(ctx context.Context, vm ports.ViewModel) (string, error) {
	_ = ctx

	colors := f.InitColors()
	var b strings.Builder

	f.WriteTitle(&b, vm, colors)
	f.WriteChildhoodSection(&b, vm, colors)
	b.WriteString("\n")
	f.WriteAdultSection(&b, vm, colors)
	f.WriteNote(&b, vm, colors)
	f.WriteComment(&b, vm, colors)

	return b.String(), nil
}

type colorScheme struct {
	title            *color.Color
	childhoodSection *color.Color
	adultSection     *color.Color
	label            *color.Color
	value            *color.Color
	quality          *color.Color
	persistence      *color.Color
	stack            *color.Color
	note             *color.Color
	comment          *color.Color
	bullet           *color.Color
	secondary        *color.Color
}

func (f *TextFormatter) InitColors() colorScheme {
	return colorScheme{
		title:            color.New(color.FgCyan, color.Bold),
		childhoodSection: color.New(color.FgYellow, color.Bold),
		adultSection:     color.New(color.FgGreen, color.Bold),
		label:            color.New(color.FgWhite, color.Bold),
		value:            color.New(color.FgHiWhite),
		quality:          color.New(color.FgHiCyan),
		persistence:      color.New(color.FgHiRed, color.Bold),
		stack:            color.New(color.FgHiMagenta),
		note:             color.New(color.FgHiBlue),
		comment:          color.New(color.FgHiYellow, color.Italic),
		bullet:           color.New(color.FgHiGreen),
		secondary:        color.New(color.FgHiBlack),
	}
}

func (f *TextFormatter) WriteTitle(b *strings.Builder, cvm ports.ViewModel, colors colorScheme) {
	if cvm.Title() != "" {
		b.WriteString(colors.title.Sprint(cvm.Title()))
		b.WriteString("\n\n")
	}
}

func (f *TextFormatter) WriteChildhoodSection(b *strings.Builder, cvm ports.ViewModel, colors colorScheme) {
	b.WriteString(colors.childhoodSection.Sprint("ДЕТСКАЯ МЕЧТА\n"))
	b.WriteString(colors.childhoodSection.Sprint("-----------------\n"))
	childhood := cvm.Childhood()
	fmt.Fprintf(b, "%s %s\n",
		colors.label.Sprint("Название:"),
		colors.value.Sprint(childhood.DisplayName()),
	)
	fmt.Fprintf(b, "%s %s\n",
		colors.label.Sprint("Роль:"),
		colors.value.Sprint(childhood.DesiredRole()),
	)
	fmt.Fprintf(b, "%s %s %s\n",
		colors.label.Sprint("Поле:"),
		colors.value.Sprint(childhood.Field().Name()),
		colors.secondary.Sprint("("+childhood.Field().Environment()+")"),
	)

	qualities := childhood.Qualities()
	if len(qualities) > 0 {
		b.WriteString(colors.label.Sprint("Качества:\n"))
		for _, q := range qualities {
			qualityNameColor := colors.quality
			if q.Name() == dream.QualityPersistence {
				qualityNameColor = colors.persistence
			}
			fmt.Fprintf(b, "  %s %s %s %s\n",
				colors.bullet.Sprint("•"),
				qualityNameColor.Sprint(q.Name()),
				colors.secondary.Sprint("—"),
				colors.value.Sprint(q.Description()),
			)
		}
	}
}

func (f *TextFormatter) WriteAdultSection(b *strings.Builder, cvm ports.ViewModel, colors colorScheme) {
	b.WriteString(colors.adultSection.Sprint("ВЗРОСЛАЯ РОЛЬ\n"))
	b.WriteString(colors.adultSection.Sprint("--------------\n"))
	adult := cvm.Adult()
	fmt.Fprintf(b, "%s %s\n",
		colors.label.Sprint("Роль:"),
		colors.value.Sprint(adult.RoleTitle()),
	)
	if desc := adult.RoleDescription(); desc != "" {
		fmt.Fprintf(b, "%s %s\n",
			colors.label.Sprint("Описание:"),
			colors.value.Sprint(desc),
		)
	}
	fmt.Fprintf(b, "%s %s %s\n",
		colors.label.Sprint("Поле:"),
		colors.value.Sprint(adult.Field().Name()),
		colors.secondary.Sprint("("+adult.Field().Environment()+")"),
	)

	stack := adult.Stack()
	if len(stack) > 0 {
		b.WriteString(colors.label.Sprint("Стек:\n"))
		for _, s := range stack {
			fmt.Fprintf(b, "  %s %s\n",
				colors.bullet.Sprint("•"),
				colors.stack.Sprint(s),
			)
		}
	}

	traits := adult.Traits()
	if len(traits) > 0 {
		b.WriteString(colors.label.Sprint("Сохранённые качества:\n"))
		for _, q := range traits {
			traitNameColor := colors.quality
			if q.Name() == dream.QualityPersistence {
				traitNameColor = colors.persistence
			}
			fmt.Fprintf(b, "  %s %s %s %s\n",
				colors.bullet.Sprint("•"),
				traitNameColor.Sprint(q.Name()),
				colors.secondary.Sprint("—"),
				colors.value.Sprint(q.Description()),
			)
		}
	}
}

func (f *TextFormatter) WriteNote(b *strings.Builder, cvm ports.ViewModel, colors colorScheme) {
	if note := cvm.Note(); note != "" {
		b.WriteString("\n")
		if strings.Contains(note, dream.QualityPersistence) {
			parts := strings.Split(note, "|")
			if len(parts) == 2 {
				b.WriteString(colors.note.Sprint(parts[0]))
				b.WriteString(" | ")
				b.WriteString(colors.persistence.Sprint(parts[1]))
			} else {
				b.WriteString(colors.note.Sprint(note))
			}
		} else {
			b.WriteString(colors.note.Sprint(note))
		}
		b.WriteString("\n")
	}
}

func (f *TextFormatter) WriteComment(b *strings.Builder, cvm ports.ViewModel, colors colorScheme) {
	if comment := cvm.Adult().Comment(); comment != "" {
		b.WriteString("\n")
		b.WriteString(colors.label.Sprint("Комментарий:\n"))
		b.WriteString(colors.comment.Sprint(comment))
		b.WriteString("\n")
	}
}
