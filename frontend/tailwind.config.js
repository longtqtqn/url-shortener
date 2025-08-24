/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  safelist: [
    'h-3', 'h-4', 'h-5', 'h-6', 'h-8',
    'w-3', 'w-4', 'w-5', 'w-6', 'w-8',
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
