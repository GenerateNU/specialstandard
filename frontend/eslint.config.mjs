import antfu from "@antfu/eslint-config";

export default antfu({
  formatters: true,
  stylistic: {
    semi: false, // Disable semicolon enforcement
    quotes: "single",
  },
  rules: {
    // Disable ALL stylistic rules if you want maximum flexibility
    "style/semi": "off",
    "style/quotes": "off",
    "style/member-delimiter-style": "off", // This is causing most of your errors
    "style/brace-style": "off", // For the curly brace errors
    "style/operator-linebreak": "off", // For operator placement errors
    "style/spaced-comment": "off", // For comment spacing
    "style/no-trailing-spaces": "off", // For trailing spaces

    // Other rules you already had
    "perfectionist/sort-imports": "off",
    "style/jsx-one-expression-per-line": "off",

    // Optional: Disable more strict rules if needed
    "ts/consistent-type-definitions": "off", // Allow both interface and type
    "ts/no-unused-vars": "warn", // Make it a warning instead of error
  },

  // Optional: Disable stylistic rules entirely
  stylistic: false, // This will disable ALL stylistic rules
});
