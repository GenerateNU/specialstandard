import antfu from "@antfu/eslint-config";

export default antfu(
  {
    stylistic: false,
    formatters: true,
    rules: {
      "no-console": "off",
      "perfectionist/sort-imports": "off",
    },
  },
  {
    ignores: ["src/lib/api/theSpecialStandardAPI.schemas.ts"],
  }
);
