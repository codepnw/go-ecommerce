package monitor

import (
	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/gofiber/fiber/v2"
)

type IMonitorHandler interface {
	HealthCheck(*fiber.Ctx) error
}

type monitorHandler struct {
	cfg config.Config
}

type Monitor struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func MonitorHandler(cfg config.Config) IMonitorHandler {
	return &monitorHandler{cfg: cfg}
}

func (m *monitorHandler) HealthCheck(c *fiber.Ctx) error {
	res := &Monitor{
		Name: m.cfg.App().Name(),
		Version: m.cfg.App().Version(),
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, res).Res()
}
