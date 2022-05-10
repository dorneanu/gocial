module.exports = {
  content: ["./templates/**/*.{html,js}"],
  theme: {
    extend: {},
    container: {
      center: true,
    },
  },
  plugins: [
    require("tailwindcss"),
    require("autoprefixer"),
    require("@tailwindcss/forms"),
    require("@tailwindcss/aspect-ratio"),
  ],
};
