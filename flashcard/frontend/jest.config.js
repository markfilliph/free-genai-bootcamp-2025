module.exports = {
  transform: {
    "^.+\\.svelte$": "svelte-jester",
    "^.+\\.js$": "babel-jest"
  },
  moduleFileExtensions: ["js", "svelte"],
  testPathIgnorePatterns: ["node_modules"],
  bail: false,
  verbose: true,
  transformIgnorePatterns: ["node_modules"],
  setupFilesAfterEnv: ["./src/__tests__/setup.js"],
  testEnvironment: "jsdom",
  moduleNameMapper: {
    // Handle import.meta.env in Vite
    "\\.(css|less|scss|sass)$": "identity-obj-proxy"
  }
};
