package data

import "time"

type Models struct {
	TokenModel    *TokenModel
	UserModel     *UserModel
	CategoryModel *CategoryModel
	ThreadModel   *ThreadModel
	PostModel     *PostModel
	TagModel      *TagModel
}

func NewModels(uri, clientToken string, pemKey []byte) Models {
	return Models{
		TokenModel: &TokenModel{
			uri:         uri,
			endpoint:    "/tokens",
			clientToken: clientToken,
			pemKey:      pemKey,
		},
		UserModel: &UserModel{
			uri:         uri,
			endpoint:    "/users",
			clientToken: clientToken,
			pemKey:      pemKey,
		},
		CategoryModel: &CategoryModel{
			uri:         uri,
			endpoint:    "/categories",
			clientToken: clientToken,
			pemKey:      pemKey,
		},
		ThreadModel: &ThreadModel{
			uri:         uri,
			endpoint:    "/threads",
			clientToken: clientToken,
			pemKey:      pemKey,
		},
		PostModel: &PostModel{
			uri:         uri,
			endpoint:    "/posts",
			clientToken: clientToken,
			pemKey:      pemKey,
		},
		TagModel: &TagModel{
			uri:         uri,
			endpoint:    "/tags",
			clientToken: clientToken,
			pemKey:      pemKey,
		},
	}
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

type envelope map[string]any

type Tokens struct {
	Authentication Token `json:"authentication_token"`
	Refresh        Token `json:"refresh_token"`
}

type Token struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}

type User struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
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
	ID           int       `json:"id"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Author       User      `json:"author"`
	IDParentPost int       `json:"id_parent_post,omitempty"`
	Thread       struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"thread"`
	Reactions  map[string]int `json:"reactions,omitempty"`
	Popularity int            `json:"popularity,omitempty"`
	Version    int            `json:"version,omitempty"`
}
