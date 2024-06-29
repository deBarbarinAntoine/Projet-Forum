package data

import "time"

type Models struct {
	UserModel     *UserModel
	CategoryModel *CategoryModel
	ThreadModel   *ThreadModel
	PostModel     *PostModel
	TagModel      *TagModel
}

func NewModels(clientToken string, pemKey []byte) Models {
	return Models{
		UserModel: &UserModel{
			endpoint:    "/users",
			clientToken: clientToken,
			pemKey:      pemKey,
		},
		CategoryModel: &CategoryModel{
			endpoint:    "/categories",
			clientToken: clientToken,
			pemKey:      pemKey,
		},
		ThreadModel: &ThreadModel{
			endpoint:    "/threads",
			clientToken: clientToken,
			pemKey:      pemKey,
		},
		PostModel: &PostModel{
			endpoint:    "/posts",
			clientToken: clientToken,
			pemKey:      pemKey,
		},
		TagModel: &TagModel{
			endpoint:    "/tags",
			clientToken: clientToken,
			pemKey:      pemKey,
		},
	}
}

type Token struct {
	Plaintext string    `json:"token"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"scope,omitempty"`
}

type password struct {
	plaintext *string
	hash      []byte
}

type User struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password      password  `json:"-"`
	Role          string    `json:"role"`
	BirthDate     time.Time `json:"birth_date"`
	Bio           string    `json:"bio,omitempty"`
	Signature     string    `json:"signature,omitempty"`
	Avatar        string    `json:"avatar,omitempty"`
	Status        string    `json:"status"`
	Version       int       `json:"-"`
	FollowingTags []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"following_tags,omitempty"`
	FavoriteThreads []struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"favorite_threads,omitempty"`
	CategoriesOwned []Category `json:"categories_owned,omitempty"`
	TagsOwned       []Tag      `json:"tags_owned,omitempty"`
	ThreadsOwned    []Thread   `json:"threads_owned,omitempty"`
	Posts           []Post     `json:"posts,omitempty"`
	Friends         []Friend   `json:"friends,omitempty"`
	Invitations     struct {
		Received []Friend `json:"received,omitempty"`
		Sent     []Friend `json:"sent,omitempty"`
	} `json:"invitations,omitempty"`
}

type Friend struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Author    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"author"`
	ParentCategory struct {
		ID   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"parent_category,omitempty"`
	Version    int        `json:"version,omitempty"`
	Categories []Category `json:"categories,omitempty"`
	Threads    []Thread   `json:"threads,omitempty"`
}

type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Author    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"author"`
	Version    int      `json:"version,omitempty"`
	Popularity int      `json:"popularity,omitempty"`
	Threads    []Thread `json:"threads,omitempty"`
}

type Thread struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      string    `json:"status"`
	Author      struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"author"`
	Category struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Version    int    `json:"version,omitempty"`
	Popularity int    `json:"popularity"`
	Posts      []Post `json:"posts,omitempty"`
	Tags       []Tag  `json:"tags,omitempty"`
}

type Post struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Author    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"author"`
	IDParentPost int `json:"id_parent_post,omitempty"`
	Thread       struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"thread"`
	Reactions  map[string]int `json:"reactions,omitempty"`
	Popularity int            `json:"popularity,omitempty"`
	Version    int            `json:"version,omitempty"`
}