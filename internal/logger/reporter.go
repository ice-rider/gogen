package logger

import (
	"fmt"
	"strings"
	"time"

	"gogen/pkg/models"
)

type Reporter struct {
	logger    *Logger
	startTime time.Time
}

func NewReporter(logger *Logger) *Reporter {
	return &Reporter{
		logger:    logger,
		startTime: time.Now(),
	}
}

func (r *Reporter) ReportStart(plan *models.GenerationPlan) {
	r.logger.Section("🚀 Генератор кода Go")

	fmt.Println("Планируется создать:")

	if len(plan.Entities) > 0 {
		fmt.Printf("  ✓ %d сущностей: %s\n",
			len(plan.Entities),
			r.joinNames(plan.Entities))
	}

	if len(plan.Repositories) > 0 {
		fmt.Printf("  ✓ %d репозиториев: %s\n",
			len(plan.Repositories),
			r.joinRepositoryNames(plan.Repositories))
	}

	if len(plan.UseCases) > 0 {
		fmt.Printf("  ✓ %d use cases: %s\n",
			len(plan.UseCases),
			r.joinUseCaseNames(plan.UseCases))
	}

	if plan.WithTests {
		fmt.Println("  ✓ Юнит-тесты для всех компонентов")
	}

	if plan.WithMocks {
		fmt.Println("  ✓ Моки для всех репозиториев")
	}

	fmt.Println()
}

func (r *Reporter) ReportProgress(componentType, name string) {
	r.logger.Success("Создано: %s %s", componentType, name)
}

func (r *Reporter) ReportComplete(plan *models.GenerationPlan, files []string) {
	duration := time.Since(r.startTime)

	r.logger.Section("🎉 Генерация завершена успешно!")

	fmt.Printf("Создано файлов: %d\n", len(files))
	fmt.Printf("Время выполнения: %s\n\n", formatDuration(duration))

	fmt.Println("Созданные файлы:")
	for _, file := range files {
		fmt.Printf("  ✓ %s\n", file)
	}

	fmt.Println("\n💡 Следующие шаги:")
	fmt.Println("  1. Проверьте сгенерированный код")
	fmt.Println("  2. Запустите: go mod tidy")
	fmt.Println("  3. Запустите тесты: go test ./...")

	if plan.WithMocks {
		fmt.Println("  4. Моки готовы для использования в тестах")
	}
}

func (r *Reporter) ReportError(err error, filesCreated []string) {
	r.logger.Section("❌ Генерация завершена с ошибкой")

	fmt.Printf("Ошибка: %v\n\n", err)

	if len(filesCreated) > 0 {
		fmt.Println("Файлы, созданные до ошибки:")
		for _, file := range filesCreated {
			fmt.Printf("  • %s\n", file)
		}
		fmt.Println("\n💡 Используйте --force для перезаписи или удалите файлы вручную")
	}
}

func (r *Reporter) ReportConflicts(conflicts []string) {
	r.logger.Section("⚠️  Обнаружены конфликты")

	fmt.Println("Следующие файлы уже существуют:")
	for _, file := range conflicts {
		fmt.Printf("  • %s\n", file)
	}

	fmt.Println("\n💡 Используйте --force для перезаписи")
}

func (r *Reporter) joinNames(entities []models.EntityConfig) string {
	names := make([]string, len(entities))
	for i, e := range entities {
		names[i] = e.Name
	}
	return strings.Join(names, ", ")
}

func (r *Reporter) joinRepositoryNames(repos []models.RepositoryConfig) string {
	names := make([]string, len(repos))
	for i, r := range repos {
		names[i] = r.Name
	}
	return strings.Join(names, ", ")
}

func (r *Reporter) joinUseCaseNames(usecases []models.UseCaseConfig) string {
	names := make([]string, len(usecases))
	for i, uc := range usecases {
		names[i] = uc.Name
	}
	return strings.Join(names, ", ")
}
