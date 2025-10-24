import antfu from "@antfu/eslint-config";

export default antfu({
  formatters: true,
  stylistic: {
    semi: false,
    quotes: "single",
  },
  rules: {
    "style/semi": "off",
    "style/quotes": "off",
    "perfectionist/sort-imports": "off",
    "style/jsx-one-expression-per-line": "off",
    "style/member-delimiter-style": "off",
    "style/operator-linebreak": "off",
    "style/brace-style": "off",
    "style/arrow-parens": "off",
    "node/prefer-global/process": "off",
    "no-console": ["warn", { allow: ["warn", "error", "log"] }], // Allow console.log
    "jsdoc/check-alignment": "off",
  },
});
