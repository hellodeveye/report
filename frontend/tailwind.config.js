const colors = require("tailwindcss/colors");

/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
    "./node_modules/vue-tailwind-datepicker/**/*.js",
  ],
  theme: {
    extend: {
      colors: {
        "vtd-primary": colors.indigo, // Light mode Datepicker color - now matches button style
        "vtd-secondary": colors.gray, // Dark mode Datepicker color palette
      },
    },
  },
  plugins: [require("@tailwindcss/forms")],
} 