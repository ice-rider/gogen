package cli

import (
	"context"
	"fmt"
	"path/filepath"

	"gogen/internal/config"
	"gogen/internal/dependency"
	"gogen/internal/file"
	"gogen/internal/format"
	"gogen/internal/generator"
	"gogen/internal/interactive"
	"gogen/internal/logger"
	"gogen/internal/project"
	"gogen/internal/template"
	"gogen/internal/util"
	"gogen/pkg/models"
)

func runGenerate(flags *Flags) error {

	if err := flags.Validate(); err != nil {
		return err
	}

	if !flags.HasComponents() && !flags.Interactive {

		return runFullInteractive()
	}

	finder := project.NewFinder(flags.OutputDir)
	root, err := finder.FindRoot()
	if err != nil {
		return fmt.Errorf("не удалось найти корень проекта: %w\nПопробуйте запустить 'gogen init'", err)
	}

	modulePath, err := finder.GetModulePath()
	if err != nil {
		return fmt.Errorf("не удалось получить module path: %w", err)
	}

	configLoader := config.NewLoader(root)
	cfg, err := configLoader.Load()
	if err != nil {
		return fmt.Errorf("не удалось загрузить конфигурацию: %w", err)
	}

	logLevel := logger.LevelInfo
	if flags.Verbose {
		logLevel = logger.LevelDebug
	}
	if flags.Quiet {
		logLevel = logger.LevelError
	}

	log := logger.NewLogger(logLevel, true)
	defer log.Close()

	reporter := logger.NewReporter(log)

	parser := NewParser()
	plan, err := parser.BuildPlan(flags)
	if err != nil {
		return fmt.Errorf("ошибка парсинга аргументов: %w", err)
	}

	plan.ModulePath = modulePath
	plan.ProjectRoot = root

	if flags.Interactive {
		interactor := interactive.NewInteractor(log)
		if err := interactor.EnhancePlan(plan, cfg); err != nil {
			return fmt.Errorf("ошибка интерактивного режима: %w", err)
		}
	}

	detector := dependency.NewDetector()
	resolver := dependency.NewResolver(detector)

	if err := resolver.Resolve(plan); err != nil {
		return fmt.Errorf("не удалось разрешить зависимости: %w", err)
	}

	if err := finder.EnsureStructure(cfg); err != nil {
		return fmt.Errorf("не удалось создать структуру папок: %w", err)
	}

	if flags.DryRun {
		return runDryRun(plan, reporter)
	}

	writer := file.NewWriter(root)
	conflictResolver := file.NewConflictResolver(flags.Interactive, flags.Force)

	expectedFiles := collectExpectedFiles(plan, cfg)
	conflicts, err := conflictResolver.CheckConflicts(expectedFiles)
	if err != nil {
		return err
	}

	if len(conflicts) > 0 && !flags.Force {
		reporter.ReportConflicts(conflicts)

		if !flags.Interactive {
			return fmt.Errorf("обнаружены конфликты, используйте --force или --interactive")
		}

		for _, conflict := range conflicts {
			overwrite, err := conflictResolver.ResolveConflict(conflict)
			if err != nil {
				return err
			}
			if !overwrite {
				return fmt.Errorf("генерация отменена пользователем")
			}
		}
	}

	templateLoader := template.NewLoader(root, cfg)
	renderer := template.NewRenderer(templateLoader)
	formatter := format.NewFormatter()
	importsManager := format.NewImportsManager()

	gen := generator.NewGenerator(renderer, writer, formatter, importsManager, cfg)

	reporter.ReportStart(plan)

	ctx := context.Background()

	if err := gen.Generate(ctx, plan); err != nil {

		log.Error("Ошибка генерации: %v", err)
		log.Info("Выполняется откат...")

		if rollbackErr := writer.Rollback(); rollbackErr != nil {
			log.Error("Ошибка отката: %v", rollbackErr)
		}

		reporter.ReportError(err, writer.GetWrittenFiles())
		return err
	}

	reporter.ReportComplete(plan, writer.GetWrittenFiles())

	return nil
}

func runDryRun(plan *models.GenerationPlan, reporter *logger.Reporter) error {
	fmt.Println("🔍 Dry-run режим - показываем что будет создано:\n")

	reporter.ReportStart(plan)

	fmt.Println("\n📋 Будут созданы следующие файлы:\n")

	for _, entity := range plan.Entities {
		fmt.Printf("  📄 internal/domain/%s.go\n", util.ToSnakeCase(entity.Name))
		if plan.WithTests {
			fmt.Printf("  📄 internal/domain/%s_test.go\n", util.ToSnakeCase(entity.Name))
		}
	}

	for _, repo := range plan.Repositories {
		fmt.Printf("  📄 internal/domain/%s_repository.go (интерфейс)\n", util.ToSnakeCase(repo.Name))
		fmt.Printf("  📄 internal/repository/%s_repository.go (реализация)\n", util.ToSnakeCase(repo.Name))
		if plan.WithTests {
			fmt.Printf("  📄 internal/repository/%s_repository_test.go\n", util.ToSnakeCase(repo.Name))
		}
		if plan.WithMocks {
			fmt.Printf("  📄 internal/mocks/%s_repository_mock.go\n", util.ToSnakeCase(repo.Name))
		}
	}

	for _, uc := range plan.UseCases {
		fmt.Printf("  📄 internal/usecase/%s_usecase.go\n", util.ToSnakeCase(uc.Name))
		if plan.WithTests {
			fmt.Printf("  📄 internal/usecase/%s_usecase_test.go\n", util.ToSnakeCase(uc.Name))
		}
	}

	fmt.Println("\n💡 Для реальной генерации уберите флаг --dry-run")

	return nil
}

func collectExpectedFiles(plan *models.GenerationPlan, cfg *models.Config) []string {
	var files []string

	root := plan.ProjectRoot

	for _, entity := range plan.Entities {
		fileName := util.ToSnakeCase(entity.Name) + ".go"
		files = append(files, filepath.Join(root, cfg.Paths.Domain, fileName))
	}

	for _, repo := range plan.Repositories {
		fileName := util.ToSnakeCase(repo.Name) + "_repository.go"
		files = append(files, filepath.Join(root, cfg.Paths.Domain, fileName))
		files = append(files, filepath.Join(root, cfg.Paths.Repository, fileName))
	}

	for _, uc := range plan.UseCases {
		fileName := util.ToSnakeCase(uc.Name) + "_usecase.go"
		files = append(files, filepath.Join(root, cfg.Paths.UseCase, fileName))
	}

	return files
}
