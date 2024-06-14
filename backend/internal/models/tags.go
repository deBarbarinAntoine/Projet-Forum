package models

type TagModel struct {
}

func (t *TagModel) Create(name, author string) int {

	// TODO tag creation method
	tag := Tag{
		Name:   name,
		Author: author,
	}
	return tag.Create()
}
