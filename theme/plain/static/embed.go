package static

import (
	_ "embed"
)

//go:embed style.css
var CSS []byte

//go:embed main.js
var JS []byte

//go:embed accordion.js
var AccordionJS []byte

//go:embed search.js
var SearchJS []byte
