package middleware

type IMiddlewareUsecase interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*Role, error)

}

type middlewareUsecase struct {
	repo IMiddlewareRepository
}

func MiddlewareUsecase(repo IMiddlewareRepository) IMiddlewareUsecase {
	return &middlewareUsecase{repo: repo}
}

func (u *middlewareUsecase) FindAccessToken(userId, accessToken string) bool {
	return u.repo.FindAccessToken(userId, accessToken)
}

func (u *middlewareUsecase) FindRole() ([]*Role, error) {
	return u.repo.FindRole()
}