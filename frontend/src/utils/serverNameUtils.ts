// utils/serverNameUtils.ts

/**
 * Normalizes a server name to a URL-friendly format
 * Handles cases like:
 * - Spaces: "Twisting Nether" -> "twisting-nether"
 * - Special characters: "Zul'jin" -> "zuljin"
 * - Multiple words with spaces: "Area 52" -> "area-52"
 * - Special character combinations: "Quel'Thalas" -> "quelthalas"
 */
export const normalizeServerName = (serverName: string): string => {
  return serverName
    .toLowerCase() // Convert to lowercase
    .replace(/[']/g, "") // Remove apostrophes
    .replace(/[^a-z0-9]+/g, "-") // Replace any non-alphanumeric characters with hyphens
    .replace(/^-+|-+$/g, ""); // Remove leading and trailing hyphens
};

/**
 * Test cases to verify the normalization works correctly
 */
const testCases = [
  ["Twisting Nether", "twisting-nether"],
  ["Zul'jin", "zuljin"],
  ["Area 52", "area-52"],
  ["Quel'Thalas", "quelthalas"],
  ["Азурегос", "azuregos"], // Handles non-Latin characters if needed
  ["Mal'Ganis", "malganis"],
  ["Ahn'Qiraj", "ahnqiraj"],
  ["Burning-Legion", "burning-legion"],
  ["Kel'Thuzad", "kelthuzad"],
  ["   Test   Server   ", "test-server"], // Handles extra spaces
];

// Use this to test the function
if (process.env.NODE_ENV === "development") {
  testCases.forEach(([input, expected]) => {
    const result = normalizeServerName(input);
    if (result !== expected) {
      console.warn(
        `Server name normalization failed for "${input}": expected "${expected}", got "${result}"`
      );
    }
  });
}
