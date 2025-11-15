package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"gogen/internal/project"
)

func NewRootCommand() *cobra.Command {
	flags := &Flags{}

	cmd := &cobra.Command{
		Use:   "gogen",
		Short: "Генератор кода для Go проектов",
		Long: `gogen - мощный генератор кода для Go проектов с поддержкой Clean Architecture.

Примеры использования:
  # Простая генерация
  gogen -d User -r User -uc CreateUser
  
  # С тестами и моками
  gogen -d Order -r Order -uc ProcessOrder -t -m
  
  # Интерактивный режим
  gogen -d User --interactive
  
  # Множественная генерация
  gogen -d User -d Product -d Order -r User -r Product -uc CreateOrder`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(flags)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	RegisterFlags(cmd, flags)

	cmd.AddCommand(NewInitCommand())
	cmd.AddCommand(NewVersionCommand())
	cmd.AddCommand(NewInteractiveCommand())

	return cmd
}

func NewInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Инициализировать проект",
		Long:  "Создаёт конфигурационный файл gogen.yaml в корне проекта",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit()
		},
	}
}

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Показать версию",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("gogen version 1.0.0")
		},
	}
}

func NewInteractiveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "interactive",
		Short: "Запустить полностью интерактивный режим",
		Long:  "Запускает интерактивный режим с пошаговым выбором компонентов",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFullInteractive()
		},
	}
}

func runInit() error {
	finder := project.NewFinder("")
	root, err := finder.FindRoot()
	if err != nil {
		return fmt.Errorf("не удалось найти корень проекта: %w", err)
	}

	configPath := filepath.Join(root, "gogen.yaml")

	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("⚠️  Файл %s уже существует\n", configPath)
		return nil
	}

	defaultConfig := `version: "1.0"

# Переопределение путей (опционально)
# paths:
#   domain: "internal/domain"
#   repository: "internal/repository"
#   usecase: "internal/usecase"

# Переопределение стиля именования (опционально)
# naming:
#   style: "snake_case"  # pascal_case | snake_case | camel_case
#   suffixes:
#     repository: "Repo"

# Кастомные шаблоны (опционально)
# templates:
#   entity: "templates/my_entity.tmpl"`

	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("не удалось создать конфиг: %w", err)
	}

	fmt.Printf("✓ Создан файл конфигурации: %s\n", configPath)
	fmt.Println("\n💡 Теперь вы можете:")
	fmt.Println("  1. Отредактировать gogen.yaml под ваши нужды")
	fmt.Println("  2. Запустить генерацию: gogen -d User -r User")

	return nil
}

func runFullInteractive() error {

	return fmt.Errorf("полностью интерактивный режим будет реализован")
}
