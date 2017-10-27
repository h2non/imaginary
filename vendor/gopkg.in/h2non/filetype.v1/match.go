package filetype

import (
	"io"
	"os"

	"gopkg.in/h2non/filetype.v1/matchers"
	"gopkg.in/h2non/filetype.v1/types"
)

// Matchers is an alias to matchers.Matchers
var Matchers = matchers.Matchers

// NewMatcher is an alias to matchers.NewMatcher
var NewMatcher = matchers.NewMatcher

// Match infers the file type of a given buffer inspecting its magic numbers signature
func Match(buf []byte) (types.Type, error) {
	length := len(buf)
	if length == 0 {
		return types.Unknown, ErrEmptyBuffer
	}

	for _, checker := range Matchers {
		match := checker(buf)
		if match != types.Unknown && match.Extension != "" {
			return match, nil
		}
	}

	return types.Unknown, nil
}

// Get is an alias to Match()
func Get(buf []byte) (types.Type, error) {
	return Match(buf)
}

// MatchFile infers a file type for a file
func MatchFile(filepath string) (types.Type, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return types.Unknown, err
	}
	defer file.Close()

	return MatchReader(file)
}

// MatchReader is convenient wrapper to Match() any Reader
func MatchReader(reader io.Reader) (types.Type, error) {
	buffer := make([]byte, 512)

	_, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return types.Unknown, err
	}

	return Match(buffer)
}

// AddMatcher registers a new matcher type
func AddMatcher(fileType types.Type, matcher matchers.Matcher) matchers.TypeMatcher {
	return matchers.NewMatcher(fileType, matcher)
}

// Matches checks if the given buffer matches with some supported file type
func Matches(buf []byte) bool {
	kind, _ := Match(buf)
	return kind != types.Unknown
}

// MatchMap performs a file matching againts a map of match functions
func MatchMap(buf []byte, matchers matchers.Map) types.Type {
	for kind, matcher := range matchers {
		if matcher(buf) {
			return kind
		}
	}
	return types.Unknown
}

// MatchesMap is an alias to Matches() but using matching againts a map of match functions
func MatchesMap(buf []byte, matchers matchers.Map) bool {
	return MatchMap(buf, matchers) != types.Unknown
}
