package models

type StringInt struct {
	Key    string
	Value  float32
	IsLast bool
}

type Project struct {
	Name    string `json:"name"`
	Id      string
	About   string            `json:"desc"`
	Links   map[string]string `json:"links"`
	Files   []FileObject
	Langs   []StringInt
	LOC     int
	LastMod string
	ReadMe  string
	Tags    []string `json:"tags"`
	Wiki    Wiki
	Saved   bool
	Path    string
	Status  string `json:"status"`
}
