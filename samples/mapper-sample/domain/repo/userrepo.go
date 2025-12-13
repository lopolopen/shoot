package repo

import "mappersample/domain/model"

type UserRepo interface {
	RepoBase[uint, *model.User]
}
