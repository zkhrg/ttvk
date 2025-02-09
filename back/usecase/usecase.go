package usecase

import (
	"back/domain/entities"
	"back/domain/interfaces"

	"github.com/sirupsen/logrus"
)

type Usecase struct {
	Repo interfaces.EntityRepository
}

func (u *Usecase) GetAddressByID(id string) (*entities.Entity, error) {
	return u.Repo.GetByID(id)
}

func (u *Usecase) EditAddressByIP(newAddress entities.EntityRequest) (*entities.Entity, error) {
	res, err := u.Repo.EditByID(newAddress)
	if err == nil {
		logrus.Error(err)
	}
	return res, err
}

func (u *Usecase) CreateAddress(AddressData entities.EntityRequest) (*entities.Entity, error) {
	return u.Repo.CreateAddress(AddressData)
}

func (u *Usecase) GetFullInfo() ([]*entities.Entity, error) {
	return u.Repo.GetFullInfo()
}

func NewUsecase(repo interfaces.EntityRepository) *Usecase {
	return &Usecase{Repo: repo}
}
