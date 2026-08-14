package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed all:dist
var distFS embed.FS

const notBuiltPage = `<!doctype html>
<html>
<body style="font-family: system-ui; padding: 3rem; color: #9e9ea7; background: #0d0d0f">
<h1 style="color: #ededef">traceknot UI is not built</h1>
<p>Run <code style="color: #6e56cf">cd ui &amp;&amp; pnpm install &amp;&amp; pnpm build</code> and restart the server.</p>
</body>
</html>`

func Handler() http.Handler {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}

	built := func() bool {
		_, err := fs.Stat(sub, "index.html")
		return err == nil
	}()

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !built {
			writer.Header().Set("content-type", "text/html; charset=utf-8")
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(notBuiltPage))
			return
		}

		path := strings.TrimPrefix(request.URL.Path, "/")
		if path == "" || strings.HasSuffix(path, "/") {
			path = "index.html"
		}
		if path == "select" {
			path = "select.html"
		}
		if _, err := fs.Stat(sub, path); err != nil {
			path = "index.html"
		}
		http.ServeFileFS(writer, request, sub, path)
	})
}
