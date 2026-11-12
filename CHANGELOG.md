Changelog
========
( English / [Japanese](./CHANGELOG_ja.md) )

v0.5.0
------
Nov 12, 2026

- Added a `Shutdown Service` link to stop the web server. (#10)
- Support graceful shutdown via Ctrl-C as well as the `Shutdown Service` link. (#11, #12)
- Replace GitHub Markdown CSS with a simpler, lightweight stylesheet (#13)
- Log the local URL when the web server starts and a message when it stops (#14)
- Support FreeBSD and macOS experimantally (#15)

v0.4.0
------
Oct 7, 2026

- Show Markdown files only by default; use `?all` to show all files. (#8)
- Added a `(Markdown Only / All)` toggle to switch between the two modes. (#8)

v0.3.0
------
Sep 18, 2026

- Added the `-start` option to open a web browser when starting the server. (#6)

v0.2.0
------
Sep 16, 2026

- Removed the `Preview` button and enabled real-time preview updates. (#2-1)
- The editing screen now displays the Markdown source on the left and the rendered result on the right. (#2-2)
- Enabled strikethrough using `~~`. (#2-3)

v0.1.0
------
Aug 26, 2026

- The first release
