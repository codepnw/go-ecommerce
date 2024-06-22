package server

import (
	"github.com/codepnw/go-ecommerce/internal/monitor"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	r fiber.Router
	s *server
}

func InitModule(r fiber.Router, s *server) IModuleFactory {
	return &moduleFactory{
		r: r,
		s: s,
	}
}

func (m *moduleFactory) MonitorModule() {
	handler := monitor.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}