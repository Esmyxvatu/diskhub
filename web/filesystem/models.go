package filesystem

type GeneralConfig struct {
	Title  string
	Auhtor string
	About  string
	Tags   []string
	Status string
}

type ExceptConfig struct {
	Files []string
	Dirs  []string
}

type ProjectFileConfig struct {
	General GeneralConfig
	Except  ExceptConfig
	Links   map[string]string
	Doc     map[string]string
}
