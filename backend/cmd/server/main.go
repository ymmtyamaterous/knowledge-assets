package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"asenare/backend/internal/config"
	"asenare/backend/internal/handler"
	"asenare/backend/internal/repository"
	"asenare/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	cfg := config.Load()

	// DB接続
	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("database connected")

	// Repositories
	userRepo := repository.NewPostgresUserRepository(db)
	courseRepo := repository.NewPostgresCourseRepository(db)
	sectionRepo := repository.NewPostgresSectionRepository(db)
	lessonRepo := repository.NewPostgresLessonRepository(db)
	progressRepo := repository.NewPostgresProgressRepository(db)
	glossaryRepo := repository.NewPostgresGlossaryRepository(db)
	quizRepo := repository.NewPostgresQuizRepository(db)
	noteRepo := repository.NewPostgresNoteRepository(db)
	badgeRepo := repository.NewPostgresBadgeRepository(db)
	searchRepo := repository.NewPostgresSearchRepository(db)

	// Use cases
	authUC := usecase.NewAuthUseCase(userRepo, cfg.JWTSecret)
	courseUC := usecase.NewCourseUseCase(courseRepo)
	sectionUC := usecase.NewSectionUseCase(sectionRepo, courseRepo)
	lessonUC := usecase.NewLessonUseCase(lessonRepo, sectionRepo)
	progressUC := usecase.NewProgressUseCase(progressRepo, lessonRepo, courseRepo, sectionRepo, quizRepo, noteRepo, badgeRepo)
	glossaryUC := usecase.NewGlossaryUseCase(glossaryRepo)
	quizUC := usecase.NewQuizUseCase(quizRepo, lessonRepo, sectionRepo)
	noteUC := usecase.NewNoteUseCase(noteRepo, lessonRepo)
	searchUC := usecase.NewSearchUseCase(searchRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authUC)
	userHandler := handler.NewUserHandler(userRepo, authUC)
	courseHandler := handler.NewCourseHandler(courseUC)
	sectionHandler := handler.NewSectionHandler(sectionUC)
	lessonHandler := handler.NewLessonHandler(lessonUC)
	progressHandler := handler.NewProgressHandler(progressUC)
	glossaryHandler := handler.NewGlossaryHandler(glossaryUC)
	quizHandler := handler.NewQuizHandler(quizUC)
	noteHandler := handler.NewNoteHandler(noteUC)
	searchHandler := handler.NewSearchHandler(searchUC)

	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		handler.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(api chi.Router) {
		// 認証
		api.Route("/auth", func(auth chi.Router) {
			auth.Post("/register", authHandler.Register)
			auth.Post("/login", authHandler.Login)
		})

		// コース（公開）
		api.Get("/courses", courseHandler.List)
		api.Get("/courses/{id}", courseHandler.Get)
		api.Get("/courses/{courseID}/sections", sectionHandler.ListByCourse)

		// セクション内レッスン（公開）
		api.Get("/sections/{sectionID}/lessons", lessonHandler.ListBySection)

		// レッスン詳細（公開）
		api.Get("/lessons/{id}", lessonHandler.Get)

		// 用語辞典（公開）
		api.Get("/glossary", glossaryHandler.ListWithFilter)
		api.Get("/glossary/tags", glossaryHandler.ListTags)
		api.Get("/glossary/daily", glossaryHandler.GetDaily)
		api.Get("/glossary/{id}", glossaryHandler.Get)

		// クイズ（公開）
		api.Get("/lessons/{lessonId}/quiz", quizHandler.GetByLesson)
		api.Get("/quizzes/{id}", quizHandler.Get)

		// 検索（公開）
		api.Get("/search", searchHandler.Search)

		// 要認証エンドポイント
		api.Group(func(private chi.Router) {
			private.Use(handler.JWTAuthMiddleware(cfg.JWTSecret))

			private.Get("/users/me", userHandler.Me)
			private.Put("/users/me", userHandler.UpdateMe)
			private.Put("/users/me/password", userHandler.ChangePassword)

			private.Post("/lessons/{id}/complete", progressHandler.CompleteLesson)
			private.Delete("/lessons/{id}/complete", progressHandler.UncompleteLesson)
			private.Get("/users/me/progress", progressHandler.GetMyProgress)
			private.Get("/users/me/course-progress", progressHandler.GetMyCourseProgress)
			private.Get("/users/me/streak", progressHandler.GetMyStreak)
			private.Get("/users/me/stats", progressHandler.GetMyStats)
			private.Get("/users/me/calendar", progressHandler.GetMyCalendar)
			private.Get("/users/me/badges", progressHandler.GetMyBadges)
			private.Post("/quizzes/{id}/submit", quizHandler.Submit)
			private.Get("/users/me/quiz-results", quizHandler.ListMyResults)

			// メモ
			private.Get("/lessons/{lessonId}/note", noteHandler.GetByLesson)
			private.Put("/lessons/{lessonId}/note", noteHandler.Save)
			private.Get("/users/me/notes", noteHandler.ListAll)
		})
	})

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.APIPort)
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
		}
	}()

	log.Printf("API server started on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
