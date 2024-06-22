package server

import (
	"github.com/codepnw/go-ecommerce/internal/middleware"
	"github.com/codepnw/go-ecommerce/internal/monitor"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
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