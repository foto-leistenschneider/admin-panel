package view

import (
	"context"
	"crypto/md5"
	"embed"
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"time"

	unixpath "path"

	"github.com/a-h/templ"
	"github.com/workos/workos-go/v4/pkg/usermanagement"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

//go:embed *.min.css *.svg
var embedFs embed.FS

var (
	filenameToServepath = map[string]string{}
	servepathToFilename = map[string]string{}
)

func init() {
	if err := fs.WalkDir(embedFs, ".", func(filename string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		f, err := embedFs.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		content, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		hash := md5.Sum(content)
		dirname := unixpath.Dir(filename)

		servepath := unixpath.Join("/", dirname, base64.RawURLEncoding.EncodeToString(hash[:]))

		filenameToServepath[filename] = servepath
		servepathToFilename[servepath] = filename

		return nil
	}); err != nil {
		fmt.Fprintf(os.Stderr, "error walking embed FS: %v", err)
		os.Exit(1)
		return
	}
}

func EmbedFSHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if filename, ok := servepathToFilename[path]; ok {
		switch unixpath.Ext(filename) {
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		}
		w.Header().Set("Cache-Control", "public, max-age=2592000") // 30 days
		f, err := embedFs.Open(filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		if _, err := io.Copy(w, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.NotFound(w, r)
	}
}

func GetServePath(filename string) templ.SafeURL {
	return templ.SafeURL(filenameToServepath[filename])
}

func timeSince(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return "N/A"
	}
	span := time.Since(ts.AsTime())
	if span > time.Hour*24 {
		return fmt.Sprintf("%d days ago", int(span.Hours()/24))
	} else if span > time.Hour {
		return fmt.Sprintf("%d hours ago", int(span.Hours()))
	} else if span > time.Minute {
		return fmt.Sprintf("%d minutes ago", int(span.Minutes()))
	} else {
		return fmt.Sprintf("%d seconds ago", int(span.Seconds()))
	}
}

func timeBetween(start, end *timestamppb.Timestamp) string {
	if start == nil || end == nil {
		return "N/A"
	}
	span := end.AsTime().Sub(start.AsTime())
	if span > time.Hour*24 {
		return fmt.Sprintf("%d days", int(span.Hours()/24))
	} else if span > time.Hour {
		return fmt.Sprintf("%d hours", int(span.Hours()))
	} else if span > time.Minute {
		return fmt.Sprintf("%d minutes", int(span.Minutes()))
	} else {
		return fmt.Sprintf("%d seconds", int(span.Seconds()))
	}
}

func getUser(ctx context.Context) (*usermanagement.User, bool) {
	v := ctx.Value("user")
	if v == nil {
		return nil, false
	}
	if user, ok := v.(*usermanagement.User); ok {
		return user, true
	}
	return nil, false
}
