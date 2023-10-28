package static

import _ "embed"

//go:embed style.css
var CSS []byte

//go:embed main.js
var JS []byte

//go:embed prism.js
var PrismJS []byte

//go:embed prism.css
var PrismCSS []byte

//go:embed prism-include-languages.js
var PrismIncludeLanguages []byte
