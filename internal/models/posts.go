package models

type PostModel struct {
}

func (p *PostModel) Create(content, author string, threadId, parentPostId int) int {

	// TODO thread creation method
	post := Post{
		Content:      content,
		Author:       author,
		ThreadId:     threadId,
		ParentPostId: parentPostId,
	}
	return post.Create()
}
