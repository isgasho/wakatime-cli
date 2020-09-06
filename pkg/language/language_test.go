package language_test

import (
	"fmt"
	// "io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/wakatime/wakatime-cli/pkg/language"

	"github.com/alecthomas/chroma/lexers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetect_ByFileExtension(t *testing.T) {
	tests := map[string]struct {
		Filepath string
		Expected string
	}{
		"assembly not gas": {
			Filepath: "path/to/file.s",
			Expected: "Assembly",
		},
		"c": {
			Filepath: "path/to/file.c",
			Expected: "C",
		},
		"c header": {
			Filepath: "path/to/file.h",
			Expected: "C",
		},
		"c++": {
			Filepath: "path/to/file.cpp",
			Expected: "C++",
		},
		"c++ 2": {
			Filepath: "path/to/file.cxx",
			Expected: "C++",
		},
		"c sharp": {
			Filepath: "path/to/file.cs",
			Expected: "C#",
		},
		"coldfusion": {
			Filepath: "path/to/file.cfm",
			Expected: "ColdFusion",
		},
		"elm": {
			Filepath: "path/to/file.elm",
			Expected: "Elm",
		},
		"f sharp not forth": {
			Filepath: "path/to/file.fs",
			Expected: "F#",
		},
		"golang": {
			Filepath: "path/to/file.go",
			Expected: "Go",
		},
		"golang modfile": {
			Filepath: "path/to/go.mod",
			Expected: "Go",
		},
		"haskell": {
			Filepath: "path/to/file.hs",
			Expected: "Haskell",
		},
		"haxe": {
			Filepath: "path/to/file.hx",
			Expected: "Haxe",
		},
		"html": {
			Filepath: "path/to/file.html",
			Expected: "HTML",
		},
		"java": {
			Filepath: "path/to/file.java",
			Expected: "Java",
		},
		"javascript": {
			Filepath: "path/to/file.js",
			Expected: "JavaScript",
		},
		"json": {
			Filepath: "path/to/file.json",
			Expected: "JSON",
		},
		"kotlin": {
			Filepath: "path/to/file.kt",
			Expected: "Kotlin",
		},
		"matlab": {
			Filepath: "path/to/file.m",
			Expected: "Matlab",
		},
		"objective c": {
			Filepath: "path/to/file.mm",
			Expected: "Objective-C",
		},
		"perl not prolog": {
			Filepath: "path/to/file.pl",
			Expected: "Perl",
		},
		"python": {
			Filepath: "path/to/file.py",
			Expected: "Python",
		},
		"rust": {
			Filepath: "path/to/file.rs",
			Expected: "Rust",
		},
		"scala": {
			Filepath: "path/to/file.scala",
			Expected: "Scala",
		},
		"swift": {
			Filepath: "path/to/file.swift",
			Expected: "Swift",
		},
		"textfile": {
			Filepath: "path/to/file.txt",
			Expected: "plaintext",
		},
		"typescript": {
			Filepath: "path/to/file.ts",
			Expected: "TypeScript",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			lang, err := language.Detect(test.Filepath)
			require.NoError(t, err)

			assert.Equal(t, test.Expected, lang)
		})
	}
}

// lexer := lexers.Match("foo.go")
// lexer := lexers.Analyse("package main\n\nfunc main()\n{\n}\n")

func _TestLanguage(t *testing.T) {
	// files, err := ioutil.ReadDir("testdata/codefiles")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// for _, file := range files {
	// 	t.Log(file.Name())
	// }

	// lexer := lexers.Match("foo.go")

	err := filepath.Walk("testdata/codefiles",
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				filename := path.Base(p)

				lexer := lexers.Match(filename)
				if lexer == nil {
					fmt.Printf("Could not detect language for file %q\n", filename)
					return nil
				}

				cfg := lexer.Config()
				fmt.Printf("  >> filename: %q, lexer: %q\n", filename, cfg.Name)
			}

			return nil
		})
	if err != nil {
		log.Println(err)

	}
	t.Fatal("intended fail")
}
