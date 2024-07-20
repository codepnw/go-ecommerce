package server

import (
	"github.com/codepnw/go-ecommerce/internal/appinfo/appinfoHandlers"
	"github.com/codepnw/go-ecommerce/internal/appinfo/appinfoRepositories"
	"github.com/codepnw/go-ecommerce/internal/appinfo/appinfoUsecases"
	"github.com/codepnw/go-ecommerce/internal/files/filesHandlers"
	"github.com/codepnw/go-ecommerce/internal/files/filesUsecases"
	"github.com/codepnw/go-ecommerce/internal/middleware"
	"github.com/codepnw/go-ecommerce/internal/monitor"
	"github.com/codepnw/go-ecommerce/internal/orders/orderHandlers"
	"github.com/codepnw/go-ecommerce/internal/orders/orderRepositories"
	"github.com/codepnw/go-ecommerce/internal/orders/orderUsecases"
	"github.com/codepnw/go-ecommerce/internal/products/productHandlers"
	"github.com/codepnw/go-ecommerce/internal/products/productRepositories"
	"github.com/codepnw/go-ecommerce/internal/products/productUsecases"
	"github.com/codepnw/go-ecommerce/internal/users/usersHandlers"
	"github.com/codepnw/go-ecommerce/internal/users/usersRepositories"
	"github.com/codepnw/go-ecommerce/internal/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FileModule()
	ProductModule()
	OrderModule()
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
	router.Get("/categories", m.m.ApiKeyAuth(), handler.FindCategory)
	router.Post("/categories", m.m.JwtAuth(), m.m.Authotize(2), handler.InsertCategory)
	router.Delete("/categories/:id", m.m.JwtAuth(), m.m.Authotize(2), handler.DeleteCategory)
}

func (m *moduleFactory) FileModule() {
	usecase := filesUsecases.FilesUsecase(m.s.cfg)
	handler := filesHandlers.FilesHandler(m.s.cfg, usecase)

	router := m.r.Group("/files")

	router.Post("/upload", m.m.JwtAuth(), m.m.Authotize(2), handler.UploadFiles)
	router.Delete("/delete", m.m.JwtAuth(), m.m.Authotize(2), handler.DeleteFile)
}

func (m *moduleFactory) ProductModule() {
	fileUsecase := filesUsecases.FilesUsecase(m.s.cfg)
	repo := productRepositories.ProductRepository(m.s.db.Get(), m.s.cfg, fileUsecase)
	usecase := productUsecases.ProductUsecase(repo)
	handler := productHandlers.ProductHandler(m.s.cfg, usecase, fileUsecase)

	router := m.r.Group("/products")

	router.Get("/", m.m.ApiKeyAuth(), handler.FindAllProducts)
	router.Post("/", m.m.JwtAuth(), m.m.Authotize(2), handler.InsertProduct)
	router.Get("/:product_id", m.m.ApiKeyAuth(), handler.FindOneProduct)
	router.Patch("/:product_id", m.m.JwtAuth(), m.m.Authotize(2), handler.UpdateProduct)
	router.Delete("/:product_id", m.m.JwtAuth(), m.m.Authotize(2), handler.DeleteProduct)
}

func (m *moduleFactory) OrderModule() {
	fileUsecase := filesUsecases.FilesUsecase(m.s.cfg)
	productRepo := productRepositories.ProductRepository(m.s.db.Get(), m.s.cfg, fileUsecase)

	orderRepo := orderRepositories.OrderRepository(m.s.db.Get())
	usecase := orderUsecases.OrderUsecase(orderRepo, productRepo)
	handler := orderHandlers.OrderHandler(m.s.cfg, usecase)

	router := m.r.Group("/orders")

	_ = handler
	_ = router
}