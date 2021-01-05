# SSGO ⚠️ REWRITE ⚠️

Very small static site generator.

## Changelog

- Rewrote the generation and parsing to be modular
  - support yaml metadata
  - will try to support other metadata arguments
  - the [file walker](cmd/post/fstructure.go) is cool
  - using interface to enable extensibility, I'm not sure of the choice yet.
- Planning to rewrite the server next
- I may include the CSS directly in the HTML so I don't have to serve the css somewhere else, not sure yet
- A lot of logging added, which is cool for developpement
- I may swap to using a CSS framework because I can't style
- The index is much more modular, I will add the descriptions of the file when trying to style it.
- I'm aiming for simplicity and not speed, but it's probably one of the fastest due to go properties.

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
