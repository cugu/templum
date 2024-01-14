<img src="./content/templum.png" width="400" height="400" align="right">

# templum

A static site generator for [templ](https://github.com/a-h/templ).

## Install

```bash
go install github.com/cugu/templum/cmd/templum@latest
```

## Usage

Generate a site from the `content` folder to the `public` folder:

```bash
templum --content content --output public --user "https://cugu.github.io/templum/"
```

## Content

Content is written in [Markdown](https://www.markdownguide.org/cheat-sheet/).

The content directory structure is:

```
content
├── config.yaml
├── index.md
├── 1 Something.md
├── 2 Something else.md
│   ├── 1 My Topic .md
│   └── 2 My Other Topic.md
├── 3 Another thing.md
└── logo.svg
```

### config.yaml

The config file contains the site base url and the GitHub url:

```yaml
github_url: "https://github.com/cugu/templum"
logo: "templum.svg"
title: "templum"
```

### Folders and Markdown files

The folder structure is used to create the navigation.
Number prefixes are used to order the pages and are removed from the navigation.
