Melw - Markdown Editor Like WikiEngine
======================================
( English / [Japanese](./README_ja.md) )

Melw is a small, local web application for editing Markdown files
with a WikiEngine-like user interface.

- Runs a web server on `localhost`.
- Directly edits Markdown files on the local filesystem.
- Uses a plain `textarea` for editing.
- ~~Requires no JavaScript.~~ __Uses [HTMX](https://htmx.org/)__
- Uses no database or proprietary file format.


![](image.png)

Install
-------

### Manual Installation

Download the binary package from [Releases](https://github.com/hymkor/melw/releases) and extract the executable.

<!-- go run github.com/hymkor/example-into-readme/cmd/how2install@master | -->

### Use [eget] installer (cross-platform)

```sh
brew install eget        # Unix-like systems
# or
scoop install eget       # Windows

cd (YOUR-BIN-DIRECTORY)
eget hymkor/melw
```

[eget]: https://github.com/zyedidia/eget

### Use [scoop]-installer (Windows only)

```
scoop install https://raw.githubusercontent.com/hymkor/melw/master/melw.json
```

or

```
scoop bucket add hymkor https://github.com/hymkor/scoop-bucket
scoop install melw
```

[scoop]: https://scoop.sh/

### Use "go install" (requires Go toolchain)

```
go install github.com/hymkor/melw@latest
```

Note: `go install` places the executable in `$HOME/go/bin` or `$GOPATH/bin`, so you need to add this directory to your `$PATH` to run `melw`.
<!-- -->


Usage
-----

1. Start the `melw` executable in a directory containing Markdown files.
2. Open [http://127.0.0.1:8000](http://127.0.0.1:8000) in your browser.
3. A list of files in the directory where melw was started will be displayed.
   - Clicking a link to a directory navigates to that directory.
   - Clicking a link to a Markdown file displays its contents rendered as HTML.
   - Clicking a link to a non-Markdown file displays its contents as-is.
4. On a Markdown file page, click the `Edit` button to edit the source of the page.
   - ~~Click `Preview` to preview the edited source rendered as HTML.~~  
     __The preview is updated almost in real time while editing.__
   - Click `Save` to save the edited source to the file and finish editing.
   - Click `Cancel` to cancel editing and discard the changes.