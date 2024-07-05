package server

import (
	"github.com/codepnw/go-ecommerce/internal/appinfo/appinfoHandlers"
	"github.com/codepnw/go-ecommerce/internal/appinfo/appinfoRepositories"
	"github.com/codepnw/go-ecommerce/internal/appinfo/appinfoUsecases"
	"github.com/codepnw/go-ecommerce/internal/middleware"
	"github.com/codepnw/go-ecommerce/internal/monitor"
	"github.com/codepnw/go-ecommerce/internal/users/usersHandlers"
	"github.com/codepnw/go-ecommerce/internal/users/usersRepositories"
	"github.com/codepnw/go-ecommerce/internal/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
}

type moduleFactory struct {
	r fiber.Router
	s *server
	m middleware.IMiddlewareHandler
}

func InitModule(r fiber.Router, s *server, m middleware.IMiddlewareHandler) IModuleFactory {
	return &moduleFactory{
		r: r,
		s: s,
		m: m,
	}
}

func InitMiddleware(s *server) middleware.IMiddlewareHandler {
	repo := middleware.MiddlewareRepository(s.db.Get())
	usecase := middleware.MiddlewareUsecase(repo)
	return middleware.MiddlewareHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitor.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repo := usersRepositories.UsersRepository(m.s.db.Get())
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repo)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/signup", m.m.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.m.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.m.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.m.ApiKeyAuth(), handler.SignOut)

	// Initial 1 admin in DB (insert sql)
	// Generate admin key
	router.Get("/admin/secret", m.m.JwtAuth(), m.m.Authotize(2), handler.GenerateAdminToken)
	router.Post("/signup-admin", m.m.JwtAuth(), m.m.Authotize(2), handler.SignUpAdmin)

	router.Get("/:user_id", m.m.JwtAuth(), m.m.ParamsCheck(), handler.GetUserProfile)
}

func (m *moduleFactory) AppinfoModule() {
	repo := appinfoRepositories.AppinfoRepository(m.s.db.Get())
	usecase := appinfoUsecases.AppinfoUsecase(repo)
	handler := appinfoHandlers.AppinfoHandler(m.s.cfg, usecase)

	router := m.r.Group("/appinfo")

	router.Get("/apikey", m.m.JwtAuth(), m.m.Authotize(2), handler.GenerateApiKey)
}
