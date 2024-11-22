/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./views/**/*.{html,js,templ,go}",  // adjust paths based on your project structure
    "./templates/**/*.{html,js,templ,go}",
    "./**/*.{html,js,templ,go}"
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
