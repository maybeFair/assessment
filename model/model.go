package model

type Resbody struct {
	Id     int      `json:"id"`
	Title  string   `json:"title"`
	Note   string   `json:"note"`
	Amount int      `json:"amount"`
	Tags   []string `json:"tags"`
}

type Reqbody struct {
	Title  string   `json:"title"`
	Note   string   `json:"note"`
	Amount int      `json:"amount"`
	Tags   []string `json:"tags"`
}

type Errmsg struct {
	Message string `json:"message"`
}
