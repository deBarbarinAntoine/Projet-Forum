package models

type Models struct {
	UserModel     *UserModel
	CategoryModel *CategoryModel
	ThreadModel   *ThreadModel
	PostModel     *PostModel
	TagModel      *TagModel
}

func NewModels(token string) Models {
	return Models{
		UserModel:     &UserModel{},
		CategoryModel: &CategoryModel{},
		ThreadModel:   &ThreadModel{},
		PostModel:     &PostModel{},
		TagModel:      &TagModel{},
	}
}
