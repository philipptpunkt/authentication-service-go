/** @type {import("eslint").Linter.Config} */
module.exports = {
  root: true,
  extends: ["@authentication-service-go/eslint-config/next.js"],
  parser: "@typescript-eslint/parser",
  parserOptions: {
    project: true,
  },
  overrides: [
    {
      files: ["*.js"], // Target JavaScript files
      env: {
        node: true, // Enable Node.js global variables
      },
      parserOptions: {
        project: null, // Disable TypeScript parser for JS files
      },
    },
  ],
}
