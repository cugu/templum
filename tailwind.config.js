const disabledCss = {
    'code::before': false,
    'code::after': false,
    'blockquote p:first-of-type::before': false,
    'blockquote p:last-of-type::after': false,
    pre: false,
    code: false,
    'pre code': false,
    'code::before': false,
    'code::after': false,
}

const defaultCss = {
    ...disabledCss,
    maxWidth: '96ch',
}

/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./theme/plain/static/in.css",
        "./theme/**/*.js",
        "./theme/**/*.templ",
        "./content/**/*.md"
    ],
    safelist: [
        "anchor",
    ],
    theme: {
        extend: {
            typography: {
                DEFAULT: {css: defaultCss},
                sm: {css: disabledCss},
                lg: {css: disabledCss},
                xl: {css: disabledCss},
                '2xl': {css: disabledCss},
            },
        },
    },
    plugins: [
        require('@tailwindcss/typography'),
    ],
}
