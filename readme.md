# SSGO ⚠️ REWRITE ⚠️

Very small static site generator.

## Changelog

- Added a Filtrer func functional style
- Simplify a few method / func signatures
- Implement a config package
- Fix the making of index files

## TODO

- Server
  - Live reload
  - File watcher
- Configuration
  - Loading from file
  - CLI
- Upgrade Poster interface
- Use metadata maps more efficiently
- Fix indexes
  - Preview of the content
  - Tags / metadata display
  - Ability to move to the next index
- Work on CSS

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
