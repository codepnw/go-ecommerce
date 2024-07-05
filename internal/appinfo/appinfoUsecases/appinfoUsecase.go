package appinfoUsecases

import "github.com/codepnw/go-ecommerce/internal/appinfo/appinfoRepositories"

type IAppinfoUsecase interface {

} 

type appinfoUsecase struct {
	repo appinfoRepositories.IAppinfoRepository
}

func AppinfoUsecase(repo appinfoRepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{repo: repo}
} 