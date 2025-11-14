import antfu from "@antfu/eslint-config";

export default antfu(
  {
    formatters: true,
    stylistic: false,
    rules: {
      "perfectionist/sort-imports": "off",
    },
  },
  {
    ignores: [
      "src/lib/api/theSpecialStandardAPI.schemas.ts",
      "src/api/models.ts",
    ],
  }
);
