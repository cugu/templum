# tempel

A static site generator for [templ](https://github.com/a-h/templ).

## Install

```bash
go install github.com/cugu/templum/cmd/tempel@latest
```

## Usage

Generate a site from the `content` folder to the `public` folder:

```bash
tempel --content content --output public
``` 

## Content

Content is written in [Markdown](https://www.markdownguide.org/cheat-sheet/). 

The content directory structure is:

```
content
├── config.yaml
├── docs
│   ├── index.md
│   ├── 1 Something.md
│   ├── 2 Something else.md
│   │   ├── 1 My Topic .md
│   │   └── 2 My Other Topic.md
│   └── 3 Another thing.md
└── static
    └── logo.svg
```

### config.yaml

The config file contains the site base url:

```yaml
base_url: "https://example.com/site/"
```

### docs

The `docs` folder contains the content of the site. 
The folder structure is used to create the navigation.
Number prefixes are used to order the pages and are removed from the navigation.

### static

The `static` folder contains static files that are copied to the output folder.

