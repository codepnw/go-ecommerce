package appinfoUsecases

import (
	"github.com/codepnw/go-ecommerce/internal/appinfo"
	"github.com/codepnw/go-ecommerce/internal/appinfo/appinfoRepositories"
)

type IAppinfoUsecase interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
	DeleteCategory(id int) error
}

type appinfoUsecase struct {
	repo appinfoRepositories.IAppinfoRepository
}

func AppinfoUsecase(repo appinfoRepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{repo: repo}
}

func (u *appinfoUsecase) FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	category, err := u.repo.FindCategory(req)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (u *appinfoUsecase) InsertCategory(req []*appinfo.Category) error {
	if err := u.repo.InsertCategory(req); err != nil {
		return err
	}
	return nil
}

func (u *appinfoUsecase) DeleteCategory(id int) error {
	if err := u.repo.DeleteCategory(id); err != nil {
		return err
	}
	return nil
}
