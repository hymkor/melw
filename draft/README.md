Melw - Markdown Editor Like WikiEngine
======================================
English / [Japanese](./README_ja.md)

Melw is a tool for browsing local files with a web user interface like a WikiEngine and editing Markdown files.

Install
-------

    go install github.com/hymkor/melw@latest


Usage
-----

1. Start the `melw` executable in a directory containing Markdown files.
2. Open [http://127.0.0.1:8000](http://127.0.0.1:8000) in your browser.
3. A list of files in the directory where melw was started will be displayed.
   - Clicking a link to a directory navigates to that directory.
   - Clicking a link to a Markdown file displays its contents rendered as HTML.
   - Clicking a link to a non-Markdown file displays its contents as-is.
4. On a Markdown file page, click the `Edit` button to edit the source of the page.
   - Click `Preview` to preview the edited source rendered as HTML.
   - Click `Save` to save the edited source to the file and finish editing.
   - Click `Cancel` to cancel editing and discard the changes.

Screenshot
----------

![](image.png)