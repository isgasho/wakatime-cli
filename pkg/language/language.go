package language

import (
	"fmt"
	"strings"

	"github.com/wakatime/wakatime-cli/pkg/heartbeat"

	"github.com/alecthomas/chroma/lexers"
	jww "github.com/spf13/jwalterweatherman"
)

// Config contains configurations for language detection.
type Config struct {
	Alternative string
	Overwrite   string
	LocalFile   string
}

// WithDetection initializes and returns a heartbeat handle option, which
// can be used in a heartbeat processing pipeline to detect and add programming
// language info to heartbeats of entity type 'file'.
func WithDetection(config Config) heartbeat.HandleOption {
	return func(next heartbeat.Handle) heartbeat.Handle {
		return func(hh []heartbeat.Heartbeat) ([]heartbeat.Result, error) {
			for n, h := range hh {
				if h.EntityType == heartbeat.FileType {
					// 1. standardize language
					// 2. Detect language
					language, err := Detect(h.Entity)
					if err != nil {
						jww.ERROR.Printf("failed to detect language on file entity %q: %s", h.Entity, err)
						continue
					}

					if language != "" {
						hh[n].Language = heartbeat.String(language)
					}
				}
			}

			return next(hh)
		}
	}
}

// Detect detects the language of a specific file.
func Detect(filepath string) (string, error) {
	switch {
	case strings.HasSuffix(filepath, ".mm"):
		return "Objective-C", nil
	case strings.HasSuffix(filepath, ".s"):
		return "Assembly", nil
	case strings.HasSuffix(filepath, ".fs"):
		return "F#", nil
	case strings.HasSuffix(filepath, ".cfm"):
		return "ColdFusion", nil
	case strings.HasSuffix(filepath, "go.mod"):
		return "Go", nil
	}

	lexer := lexers.Match(filepath)
	if lexer != nil {
		return lexer.Config().Name, nil
	}

	return "", fmt.Errorf("Could not detect language for file %q\n", filepath)
}
