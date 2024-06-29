package data

type ThreadModel struct {
	endpoint    string
	clientToken string
	pemKey      []byte
}

func (t *ThreadModel) Create(title, description string, isPublic bool, author, category string) int {

	// TODO thread creation method
	thread := Thread{
		Title:       title,
		Description: description,
		IsPublic:    isPublic,
		Author:      author,
		Category:    category,
	}
	return thread.Create()
}
