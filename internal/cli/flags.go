package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Flags struct {
	Entities     []string
	Repositories []string
	UseCases     []string
	Handlers     []string

	WithTests   bool
	WithMocks   bool
	Interactive bool
	Force       bool
	DryRun      bool

	Verbose bool
	Quiet   bool
	NoColor bool

	ConfigPath string
	OutputDir  string
}

func RegisterFlags(cmd *cobra.Command, flags *Flags) {

	cmd.Flags().StringSliceVarP(&flags.Entities, "entity", "d", []string{},
		"Создать сущность (можно указать несколько раз)")
	cmd.Flags().StringSliceVarP(&flags.Repositories, "repo", "r", []string{},
		"Создать репозиторий (можно указать несколько раз)")
	cmd.Flags().StringSliceVarP(&flags.UseCases, "usecase", "uc", []string{},
		"Создать use case (можно указать несколько раз)")
	cmd.Flags().StringSliceVar(&flags.Handlers, "handler", []string{},
		"Создать HTTP handler (можно указать несколько раз)")

	cmd.Flags().BoolVarP(&flags.WithTests, "with-tests", "t", false,
		"Генерировать тесты для всех компонентов")
	cmd.Flags().BoolVarP(&flags.WithMocks, "with-mocks", "m", false,
		"Генерировать моки для репозиториев")
	cmd.Flags().BoolVar(&flags.Interactive, "interactive", false,
		"Интерактивный режим с дополнительными вопросами")
	cmd.Flags().BoolVar(&flags.Force, "force", false,
		"Перезаписывать существующие файлы без подтверждения")
	cmd.Flags().BoolVar(&flags.DryRun, "dry-run", false,
		"Показать что будет создано без реального создания файлов")

	cmd.Flags().BoolVarP(&flags.Verbose, "verbose", "v", false,
		"Подробный вывод")
	cmd.Flags().BoolVarP(&flags.Quiet, "quiet", "q", false,
		"Минимальный вывод")
	cmd.Flags().BoolVar(&flags.NoColor, "no-color", false,
		"Отключить цветной вывод")

	cmd.Flags().StringVarP(&flags.ConfigPath, "config", "c", "",
		"Путь к конфигурационному файлу")
	cmd.Flags().StringVarP(&flags.OutputDir, "output", "o", "",
		"Директория для генерации (по умолчанию текущая)")
}

func (f *Flags) Validate() error {

	if len(f.Entities) == 0 &&
		len(f.Repositories) == 0 &&
		len(f.UseCases) == 0 &&
		len(f.Handlers) == 0 {

		return nil
	}

	if f.Quiet && f.Verbose {
		return fmt.Errorf("--quiet and --verbose cannot be used together")
	}

	return nil
}

func (f *Flags) HasComponents() bool {
	return len(f.Entities) > 0 ||
		len(f.Repositories) > 0 ||
		len(f.UseCases) > 0 ||
		len(f.Handlers) > 0
}
