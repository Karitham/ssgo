# SSGO

Very small static site generator.

It currently support folder structures, but the navigation is not optimal

I am focused on making the thing work right now, don't use this in prod

* [SSGO](#ssgo)
  * [Installing](#installing)
  * [Usage](#usage)
  * [Structure](#structure)
  * [Writing](#writing)
  * [Styling](#styling)
* [Developpement](#developpement)
  * [Templates](#templates)
  * [Dependencies](#dependencies)
* [License](#license)
* [Author](#author)
  
## Installing

`go get github.com/Karitham/ssgo`

This will install as a binary ready to use if the right directory is in your path.

Just reproduce the following folder structure and use `ssgo` to generate HTML

## Usage

There are a few commands you can use to help you during the writing of your site.

````help
❯ ssgo --help
NAME:
   SSGO - Generate HTML based on the markdown you provide, easily customizable with your own theme

USAGE:
   ssgo [global options] command [command options]

COMMANDS:
   server, s  serve your files with a live reloading server
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --post value  Change the post directory (default: "posts")
   --publ value  Change the publication directory (default: "public")
   --tmpl value  Change the template directory (default: "templates")
   --help, -h    show help (default: false
````

The server is useful if you want to have a preview of what your site looks like. At each save of the file you are working on, the according part of the site is reloaded (be it css or text).

There are available command options for each command, which can be used to change the directories you read / write to etc

## Structure

Here is the default folder structure. You can change where the templates / posts / publication directory are by using global options

```tree
.
├───assets
│   └───css
│           style.css
│           style.sass
│
├───posts
│   │   about.md
│   │   _empty.md
│   │
│   └───Projects
│           Random-RSS.md
│
├───public
│   │   about.html
│   │   index.html
│   │
│   └───Projects
│           index.html
│           Random-RSS.html
│
└───templates
        index.tmpl
        post.tmpl
```

## Writing

Write a markdown post in `posts`, or your own folder

All the files starting with `_` will not be generated, so you can use that to make drafts

Index files are generated to make the navigation easy.

## Styling

The theme is gruvbox-like, because I love gruvbox.

You can customize everything in `assets/css` to make your own style

There is technically no restriction on style or how you structure it, just make sure you update the templates accordingly

# Developpement

## Templates

For now it only support 2 templates and the code is not exactly modular

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

Pierre-Louis "Karitham" Pery
