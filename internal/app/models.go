package app

type reqURL struct {
	URL string `json:"url"`
}

type respURL struct {
	URL string `json:"result"`
}

type loginIn struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type singIn struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}
