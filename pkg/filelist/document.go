package filelist

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
)

var errMissingFile = errors.New("missing / partial file")

// Document is a parsed file the viewer can validate, describe, and show
// metadata for. ACH files and each FedACH report kind implement this.
type Document interface {
	Kind() string
	Title() string
	Validate() error
	RecordCount() int64
	Metadata() map[string]string
}

// HumanWriter is implemented by documents that can render themselves as
// human-readable text. ACH files are handled separately so they can apply
// display masking.
type HumanWriter interface {
	WriteHuman(w io.Writer) error
}

// Format parses one family of files (ACH, FAHK/ACK, future FedACH reports).
type Format struct {
	// Kind is a short identifier such as "ach" or "ack".
	Kind string
	// Title is a human-readable label for the format.
	Title string
	// Extensions are filename suffixes this format owns, including the dot.
	Extensions []string
	// Parse turns file bytes into a Document. Pending formats return a
	// Document whose Validate error is ErrUnsupported.
	Parse func(name string, r io.Reader) (Document, error)
}

var (
	formatsMu sync.RWMutex
	formats   []Format
	byExt     = map[string]Format{}
)

// RegisterFormat adds a parser for the given extensions. Later registrations
// for the same extension replace earlier ones. Call this from init() in the
// file that implements a new FedACH report kind.
func RegisterFormat(f Format) {
	if f.Parse == nil {
		panic("filelist: RegisterFormat requires Parse") //nolint:forbidigo
	}

	formatsMu.Lock()
	defer formatsMu.Unlock()

	replaced := false
	for i, existing := range formats {
		if existing.Kind == f.Kind {
			formats[i] = f
			replaced = true
			break
		}
	}
	if !replaced {
		formats = append(formats, f)
	}
	for _, ext := range f.Extensions {
		byExt[normalizeExt(ext)] = f
	}
}

// Formats returns the registered parsers in registration order.
func Formats() []Format {
	formatsMu.RLock()
	defer formatsMu.RUnlock()
	out := make([]Format, len(formats))
	copy(out, formats)
	return out
}

// Read parses an ACH file or FedACH report from r using name's extension.
func Read(name string, r io.Reader) (*File, error) {
	return readFile(name, r)
}

func readFile(name string, r io.Reader) (*File, error) {
	ext := normalizeExt(filepath.Ext(name))

	formatsMu.RLock()
	format, ok := byExt[ext]
	formatsMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown file extension for %s", name)
	}

	doc, err := format.Parse(name, r)
	out := &File{Name: name, Document: doc}
	return out, err
}

func mergeParsed(dst, parsed *File) *File {
	if dst == nil {
		dst = &File{}
	}
	if parsed == nil {
		return dst
	}
	if dst.Name == "" {
		dst.Name = parsed.Name
	}
	dst.Document = parsed.Document
	return dst
}

func normalizeExt(ext string) string {
	ext = strings.ToLower(strings.TrimSpace(ext))
	if ext != "" && !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return ext
}

// ErrUnsupported is returned when a filename belongs to a known FedACH
// (or other) report type that is not parsed yet.
type ErrUnsupported struct {
	Kind  string
	Title string
	Ext   string
}

func (e ErrUnsupported) Error() string {
	title := e.Title
	if title == "" {
		title = e.Kind
	}
	if e.Ext != "" {
		return fmt.Sprintf("%s (%s) files are not yet supported", title, e.Ext)
	}
	return fmt.Sprintf("%s files are not yet supported", title)
}

func (e ErrUnsupported) Is(target error) bool {
	other, ok := target.(ErrUnsupported)
	if !ok {
		return false
	}
	if other.Kind != "" && other.Kind != e.Kind {
		return false
	}
	if other.Ext != "" && other.Ext != e.Ext {
		return false
	}
	return true
}
