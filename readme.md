# ⚠️ SSGO REWRITE ⚠️

Very small static site generator.

Rewrite happening, this means a few things:

- Things are breaking
- Things may not work
- Features are missing
- Theming is not done and completely optional
- Performance may not be there
- Code is cleaner

## Plans

Here are the goals planned for this rewrite

If anybody sees this and has any suggestion, please open an issue.

### TODO

- [ ]  Server
  - [ ]  Live reload
    - [ ]  Use of Watcher to re generate the files and indexes
  - [ ]  File watcher
    - [ ]  Interface with the fs
- [ ]  Configuration
  - [ ]  Loading from file
  - [ ]  CLI
- [ ]  Use metadata maps more efficiently
- [ ]  Fix indexes
  - [ ]  Preview of the content
  - [ ]  Tags / metadata display
  - [ ]  Date use ?
  - [ ]  Ability to move to the next index
- [ ]  Work on CSS (devide Index & Article)

### DONE

- [x]  Filter files while flattening
- [x]  Implement a better version of the Runner
- [x]  Fix data race on index building
- [x]  Rewrite the generation and parsing to be modular
  - [x]  YAML metadata
    - [x]  Use a Map for metadata
    - [x]  Made the metadata parsing generic to any `Poster`
  - [x]  build a good [file walker](https://github.com/Karitham/ssgo/blob/rewrite/cmd/post/fs.go)
- [x]  Logging
- [x]  Make index modular

### Maybe

- Use a CSS Framework
- Improve the fs package and release as seperate

## Structure

Here is the default folder structure.

```tree
.
├───assets
│   ├───css
│   │       style.css
│   │
│   └───templates
│           index.tmpl
│           post.tmpl
│
├───posts
│   │   about.md
│   │   _draft.md
│   │
│   └───Projects
│           WaifuBot.md
│
└───public
    │   about.html
    │   index.html
    │
    └───Projects
            index.html
            WaifuBot.html
```

# License

```license
This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <http://unlicense.org/>
```

# Author

Pierre-Louis "Karitham" Pery
