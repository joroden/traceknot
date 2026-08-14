# Third-party notices

This project embeds or adapts some third-party materials directly (listed
below with full license text), and links or bundles a number of open-source
libraries into the built binary and UI (listed by license family further
down, with pointers to their full license text).

## tiktoken (OpenAI)

- **Component:** o200k_base mergeable ranks vocabulary
  (`internal/tokenize/vocab/o200k_base.tiktoken`) and the o200k_base
  pretokenization regular expression (`internal/tokenize/bpe.go`).
- **Source:** https://github.com/openai/tiktoken
- **License:** MIT (see below). Used for estimating token usage of tool
  calls; the Go byte-pair merge implementation in `internal/tokenize` is an
  original implementation of the published BPE algorithm (Sennrich, Haddow,
  Birch, 2016) that matches tiktoken's merge ordering.

```
MIT License

Copyright (c) 2022 OpenAI, Shantanu Jain

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## Runtime dependencies

Every dependency actually compiled into the `traceknot` binary or bundled
into the UI build is permissively licensed (MIT, BSD-3-Clause, Apache-2.0,
ISC, or OFL-1.1 for fonts) — none are copyleft (GPL, LGPL, AGPL, MPL). No
dependency source is copied into this repository; full license text for each
is available where noted below.

### Go modules

Scope: `go list -deps ./cmd/traceknot/...`, i.e. what actually ends up in the
binary — not every module listed in `go.mod` (a few, like `fsnotify`, are
declared but not currently imported by any code path). Full text is in the
module cache: `$(go env GOMODCACHE)/<module>@<version>/LICENSE`.

**MIT** — `charmbracelet/{huh, bubbles, bubbletea, colorprofile, lipgloss,
x/ansi, x/cellbuf, x/exp/strings, x/term}`, `aymanbagabas/go-osc52/v2`,
`catppuccin/go`, `dustin/go-humanize`, `lucasb-eyer/go-colorful`,
`mattn/{go-isatty, go-runewidth}`, `mitchellh/hashstructure/v2`,
`muesli/{ansi, cancelreader, termenv}`, `rivo/uniseg`, `xo/terminfo`

**BSD-3-Clause** — `atotto/clipboard`, `google/uuid`, `remyoudompheng/bigfft`,
`golang.org/x/sync`, `golang.org/x/sys`, `google.golang.org/protobuf`,
`modernc.org/{libc, mathutil, memory, sqlite}`

**Apache-2.0** — `go.opentelemetry.io/proto/otlp`

### UI (npm) packages

Scope: the `dependencies` in `ui/package.json` (what Vite bundles into the
shipped UI), not `devDependencies` (build-time only: TypeScript, Vite, type
stubs — never shipped). Full text is in `ui/node_modules/<package>/LICENSE`
where present, or the package's own repository.

**MIT** — `react`, `react-dom`, `react-router`, `recharts`,
`@radix-ui/react-dialog`, `@radix-ui/react-dropdown-menu`,
`@tanstack/react-table`, `@tanstack/react-virtual`, `date-fns`,
`tailwindcss`, `@tailwindcss/vite`

**ISC** — `lucide-react`

**OFL-1.1** (SIL Open Font License) — `@fontsource/inter`,
`@fontsource/jetbrains-mono`, packaging the Inter and JetBrains Mono
typefaces as bundled font files. OFL is written for exactly this kind of
redistribution (bundling fonts with an application); it requires preserving
the license and font names, not renaming a modified font as the original.
Canonical text: https://openfontlicense.org
