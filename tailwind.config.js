/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./theme/**/*.templ"
    ],
    theme: {
        extend: {},
    },
    plugins: [
        require('@tailwindcss/typography'),
    ],
}
