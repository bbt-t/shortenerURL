package app

type reqURL struct {
	URL string `json:"url"`
}

type respURL struct {
	URL string `json:"result"`
}

type singIn struct {
	UserName string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}
