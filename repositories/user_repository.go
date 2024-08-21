package repositories

import (
	"dating-app-api/entities/models"
	"dating-app-api/entities/requests"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	CreateUser(model *models.UserModel, tx *gorm.DB) (*models.UserModel, error)
	GetDetailUser(whereClause interface{}, whereNotClause interface{}, orClause interface{}, relations []string) (*models.UserModel, error)
	GetListUser(meta *requests.MetaPaginationRequest, whereClause interface{}, whereNotClause interface{}, orClause interface{}, relations []string) ([]*models.UserModel, int64, error)
	UpdateUser(model *models.UserModel, tx *gorm.DB) (*models.UserModel, error)
	DeleteUser(model *models.UserModel, tx *gorm.DB) error
}

type userReposiotry struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userReposiotry{
		db: db,
	}
}

func (repo *userReposiotry) CreateUser(model *models.UserModel, tx *gorm.DB) (*models.UserModel, error) {
	err := tx.Create(&model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (repo *userReposiotry) GetDetailUser(whereClause interface{}, whereNotClause interface{}, orClause interface{}, relations []string) (*models.UserModel, error) {
	var user *models.UserModel

	queryBuilder := repo.db.Where(whereClause)
	if whereNotClause != nil {
		queryBuilder.Not(whereNotClause)
	}

	if orClause != nil {
		queryBuilder.Or(orClause)
	}

	err := queryBuilder.First(&user).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return nil, nil
	case nil:
		return user, nil
	default:
		return nil, err
	}
}

func (repo *userReposiotry) GetListUser(meta *requests.MetaPaginationRequest, whereClause interface{}, whereNotClause interface{}, orClause interface{}, relations []string) ([]*models.UserModel, int64, error) {

	var users []*models.UserModel

	queryBuilder := repo.db.Debug().Table("users")

	if whereClause != nil {
		queryBuilder.Where(whereClause)
	}

	if whereNotClause != nil {
		queryBuilder.Not(whereNotClause)
	}

	if orClause != nil {
		queryBuilder.Or(orClause)
	}

	if len(relations) > 0 {
		for _, relation := range relations {
			queryBuilder.Preload(relation)
		}
	}

	var totalRows int64
	if err := queryBuilder.Where("deleted_at is null").Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	queryBuilder.Limit(meta.Limit).Offset(meta.Offset).Order(meta.Order)

	if err := queryBuilder.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalRows, nil
}

func (repo *userReposiotry) UpdateUser(model *models.UserModel, tx *gorm.DB) (*models.UserModel, error) {
	err := tx.Where("id = ?", model.Id).Updates(&model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (repo *userReposiotry) DeleteUser(model *models.UserModel, tx *gorm.DB) error {
	err := tx.Where("id = ?", model.Id).Delete(&model).Error
	if err != nil {
		return err
	}

	return nil

}
