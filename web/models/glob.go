package models

var Projects []Project
var ExcludeList []string = []string{
	// Git
	".git",
	".gitignore",
	".gitattributes",

	// JavaScript - TypeScript
	"node_modules",
	"package-lock.json",
	"tsconfig.json",

	// Rust
	"target",
	"Cargo.lock",

	// Golang
	"go.mod",
	"go.sum",

	// Zig
	"zig-build",

	// Python
	"__pycache__",
	"requirements.txt",

	// Compiled programs
	"*.out",
	"*.exe",

	// Cpp
	".cache",
}
