package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/time/rate"

	"github.com/maxwellpark/stanzabonanza/backend/internal/config"
	"github.com/maxwellpark/stanzabonanza/backend/internal/handler"
	"github.com/maxwellpark/stanzabonanza/backend/internal/middleware"
	"github.com/maxwellpark/stanzabonanza/backend/internal/repository"
	"github.com/maxwellpark/stanzabonanza/backend/internal/service"
)

type Server struct {
	cfg    *config.Config
	db     *pgxpool.Pool
	router *chi.Mux
	logger *slog.Logger
}

func New(cfg *config.Config, db *pgxpool.Pool, logger *slog.Logger) *Server {
	s := &Server{
		cfg:    cfg,
		db:     db,
		router: chi.NewRouter(),
		logger: logger,
	}
	s.setupMiddleware()
	s.setupRoutes()
	return s
}

func (s *Server) setupMiddleware() {
	rl := middleware.NewRateLimiter(rate.Limit(100), 200)

	s.router.Use(chimw.Recoverer)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logging(s.logger))
	s.router.Use(rl.Middleware)
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   s.cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))
}

func (s *Server) setupRoutes() {
	userRepo := repository.NewUserRepository(s.db)
	sessionRepo := repository.NewSessionRepository(s.db)
	magicLinkRepo := repository.NewMagicLinkRepository(s.db)
	webAuthnRepo := repository.NewWebAuthnRepository(s.db)
	poemRepo := repository.NewPoemRepository(s.db)
	stanzaRepo := repository.NewStanzaRepository(s.db)
	likeRepo := repository.NewLikeRepository(s.db)
	commentRepo := repository.NewCommentRepository(s.db)
	followRepo := repository.NewFollowRepository(s.db)
	notifRepo := repository.NewNotificationRepository(s.db)

	tutorialRepo := repository.NewTutorialRepository(s.db)

	authSvc := service.NewAuthService(userRepo, sessionRepo, magicLinkRepo, webAuthnRepo, s.cfg)
	poemSvc := service.NewPoemService(poemRepo, stanzaRepo, notifRepo)
	socialSvc := service.NewSocialService(likeRepo, commentRepo, followRepo, notifRepo, poemRepo)
	tutorialSvc := service.NewTutorialService(tutorialRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	poemHandler := handler.NewPoemHandler(poemSvc)
	socialHandler := handler.NewSocialHandler(socialSvc)
	tutorialHandler := handler.NewTutorialHandler(tutorialSvc)

	authMiddleware := middleware.Auth(authSvc)
	optionalAuth := middleware.OptionalAuth(authSvc)

	// 3 req/s burst 5 - tight limit for endpoints that send email or initiate auth ceremonies.
	strictRL := middleware.NewRateLimiter(rate.Limit(3), 5)

	s.router.Get("/health", handler.Health(s.db))

	s.router.Route("/api/v1", func(r chi.Router) {
		// Auth (public)
		r.Route("/auth", func(r chi.Router) {
			r.With(strictRL.Middleware).Post("/magic-link", authHandler.RequestMagicLink)
			r.Get("/magic-link/verify", authHandler.VerifyMagicLink)
			r.With(authMiddleware).Post("/register/begin", authHandler.BeginRegistration)
			r.With(authMiddleware).Post("/register/finish", authHandler.FinishRegistration)
			r.With(strictRL.Middleware).Post("/login/begin", authHandler.BeginLogin)
			r.Post("/login/finish", authHandler.FinishLogin)
			r.With(authMiddleware).Post("/logout", authHandler.Logout)
			r.With(authMiddleware).Get("/me", authHandler.Me)
		})

		// Poems
		r.Route("/poems", func(r chi.Router) {
			r.With(optionalAuth).Get("/", poemHandler.List)
			r.With(authMiddleware).Post("/", poemHandler.Create)
			r.Route("/{poemID}", func(r chi.Router) {
				r.With(optionalAuth).Get("/", poemHandler.Get)
				r.With(authMiddleware).Put("/", poemHandler.Update)
				r.With(authMiddleware).Delete("/", poemHandler.Delete)

				// Stanzas
				r.Get("/stanzas", poemHandler.ListStanzas)
				r.With(authMiddleware).Post("/stanzas", poemHandler.SubmitStanza)
				r.With(authMiddleware).Put("/stanzas/{stanzaID}", poemHandler.ReviewStanza)

				// Social on poems
				r.With(authMiddleware).Post("/like", socialHandler.ToggleLike)
				r.Get("/comments", socialHandler.ListComments)
				r.With(authMiddleware).Post("/comments", socialHandler.AddComment)
			})
		})

		// Users
		r.Route("/users", func(r chi.Router) {
			r.With(authMiddleware).Put("/me", authHandler.UpdateProfile)
			r.Get("/{userID}", authHandler.GetProfile)
			r.Get("/{userID}/poems", poemHandler.ListByUser)
			r.With(authMiddleware).Post("/{userID}/follow", socialHandler.ToggleFollow)
			r.Get("/{userID}/followers", socialHandler.ListFollowers)
			r.Get("/{userID}/following", socialHandler.ListFollowing)
		})

		// Feed
		r.With(authMiddleware).Get("/feed", poemHandler.Feed)
		r.Get("/explore", poemHandler.Explore)
		r.Get("/hall-of-fame", poemHandler.HallOfFame)

		// Notifications
		r.With(authMiddleware).Get("/notifications", socialHandler.ListNotifications)
		r.With(authMiddleware).Post("/notifications/read", socialHandler.MarkNotificationsRead)

		// Comments
		r.With(authMiddleware).Delete("/comments/{commentID}", socialHandler.DeleteComment)

		// Tutorials
		r.Get("/tutorials", tutorialHandler.List)
		r.Get("/tutorials/{slug}", tutorialHandler.Get)
	})
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s.logger.Info("server starting", "port", s.cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	s.logger.Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
