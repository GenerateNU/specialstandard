import antfu from '@antfu/eslint-config'

// I had to set this because my build kept failing??
export default antfu({
  formatters: true,
  // stylistic: {
  //   semi: false,
  //   quotes: "single",
  // },
  // rules: {
  //   "style/semi": "off",
  //   "style/quotes": "off",
  //   "perfectionist/sort-imports": "off",
  //   "style/jsx-one-expression-per-line": "off",
  // },
})
