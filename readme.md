# SSGO - ⚠️HEAVILY WIP⚠️

Very small static site generator.

It currently support basic folder structures, but the navigation is clearly not optimal

I am focusing on making the thing work right now, it is nowhere near usable in prod

* [SSGO - ⚠️HEAVILY WIP⚠️](#ssgo---️heavily-wip️)
* [Use](#use)
  * [Install](#install)
  * [Structure](#structure)
  * [Write](#write)
  * [Style](#style)
* [Developpement](#developpement)
  * [Templates](#templates)
  * [Dependencies](#dependencies)
* [License](#license)
* [Author](#author)
  
# Use

## Install

`go get github.com/Karitham/ssgo`

This will install as a binary ready to use if the right directory is in your path.

Just reproduce the following folder structure and use `ssgo` to generate HTML

## Structure

Here is an exemple folder structure, you only need the assets directory layed out with the templates in, and make posts in `posts` the rest is flexible.

You can change the folder structure by changing the source code.

```tree
.
├───assets
│   ├───css
│   │       style.css
│   │       style.sass
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
│           Random-RSS.md
│
└───public
    │   about.html
    │   index.html
    │
    └───Projects
            index.html
            Random-RSS.html
```

## Write

Write a markdown post in `posts`

All the files starting with `_` will not be generated, so you can use that to make drafts

The output folder is `public`. It generates index files to navigate in the folder structure easily.

## Style

The theme is gruvbox-like, because I love gruvbox.

You can customize everything in `assets/css` to make your own style

There is technically no restriction on style or how you structure it, just make sure you update the templates accordingly

# Developpement

## Templates

For now it only support 2 templates and the code is not modular

## Dependencies

[https://github.com/yuin/goldmark](https://github.com/yuin/goldmark)

[https://github.com/alecthomas/chroma](github.com/alecthomas/chroma)

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

PL "Karitham" Pery
