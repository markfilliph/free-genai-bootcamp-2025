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
  setupFilesAfterEnv: ["./src/__tests__/setup/index.js"],
  testEnvironment: "jsdom",
  moduleNameMapper: {
    // Handle import.meta.env in Vite
    "\\.(css|less|scss|sass)$": "identity-obj-proxy"
  },
  collectCoverageFrom: [
    "src/**/*.{js,svelte}",
    "!src/main.js",
    "!src/__tests__/**",
    "!**/node_modules/**"
  ],
  coverageReporters: ["text", "lcov", "json-summary"],
  coverageThreshold: {
    global: {
      statements: 50,
      branches: 40,
      functions: 50,
      lines: 50
    },
    // Critical components with low coverage that need improvement
    "src/components/StudySession.svelte": {
      statements: 30,
      branches: 20,
      functions: 30,
      lines: 30
    },
    "src/components/FlashcardReview.svelte": {
      statements: 30,
      branches: 20,
      functions: 30,
      lines: 30
    },
    "src/routes/Login.svelte": {
      statements: 30,
      branches: 20,
      functions: 30,
      lines: 30
    },
    // Secondary components that need coverage improvement
    "src/components/Deck.svelte": {
      statements: 20,
      branches: 15,
      functions: 20,
      lines: 20
    },
    "src/components/Navbar.svelte": {
      statements: 20,
      branches: 15,
      functions: 20,
      lines: 20
    },
    "src/routes/DeckManagement.svelte": {
      statements: 20,
      branches: 15,
      functions: 20,
      lines: 20
    }
  }
};
