package users

type Links struct {
	Old       string `json:"old"`
	Short     string `json:"short"`
	CreatedBy string `json:"created_by"`
}

type User struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Links    []Links `json:"urls"`
	User_id  string  `json:"user_id"`
}

type Mix struct {
	Username string `json:"username"`
	Old      string `json:"old"`
}
