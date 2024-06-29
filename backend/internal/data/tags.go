package data

type TagModel struct {
	uri         string
	endpoint    string
	clientToken string
	pemKey      []byte
}

func (t *TagModel) Create(name, author string) int {

	// TODO tag creation method
	tag := Tag{
		Name:   name,
		Author: author,
	}
	return tag.Create()
}
