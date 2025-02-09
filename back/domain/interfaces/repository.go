package interfaces

import "back/domain/entities"

type EntityRepository interface {
	GetByID(id string) (*entities.Entity, error)
	EditByID(newAddress entities.EntityRequest) (*entities.Entity, error)
	CreateAddress(addressData entities.EntityRequest) (*entities.Entity, error) // у криейта не должно быть того же типа потому что у него нет и должно быть передаваемого ууид
	GetFullInfo() ([]*entities.Entity, error)
}
