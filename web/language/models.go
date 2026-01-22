package language

type Commentaire struct {
	Start string
	End   string
}

type Language struct {
	IsCode bool
	Name   string
	Exts   []string
	Coms   []Commentaire
}

type LanguageStat struct {
	Name string
	Loc  float32
	File int
}

type IndexLangs struct {
	Name  string
	Loc   int
	Files int
	Projs int
}