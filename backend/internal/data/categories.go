package data

type CategoryModel struct {
	endpoint    string
	clientToken string
	pemKey      []byte
}

func (c *CategoryModel) Create(name, author, parentCategory string) int {
	// TODO implement category creation method
	category := &Category{
		Name:           name,
		Author:         author,
		ParentCategory: parentCategory,
	}
	return category.Create()
}
