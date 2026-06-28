/*global module, process*/

/**
 * ZeroVer - Versioning with a ZeroVer philosophy.
 *
 * Main entry point for the ZeroVer action. Manages version calculation and tag creation
 * following a ZeroVer approach where the major version is never incremented.
 *
 * Behavior:
 * - When create=true: Creates and pushes a new tag with incremented version
 *   * Strategy 'rollover': Increments patch, rolls over to minor at rollover-max
 *   * Strategy 'conventional-commits': Analyzes commits to determine version bump
 *   * Returns clean version (e.g., "0.1.5")
 *
 * - When create=false: Returns current version without creating a tag
 *   * If commit is tagged: Returns clean version (e.g., "0.1.4")
 *   * If commit is not tagged: Returns dev version with timestamp and hash
 *     (e.g., "0.1.4-dev.20260628143052.abc123def456")
 *
 * Outputs:
 * - version: The calculated version (without 'v' prefix)
 * - created: Boolean indicating if a new tag was created
 */
module.exports = async ({ core, exec }) => {
  const shouldCreate = process.env.create === "true";

  let tag;
  let created = false;

  if (shouldCreate) {
    const result = await getAndCreateNextTag({ core, exec });
    tag = result.tag;
    created = result.created;
  } else {
    tag = await getCurrentVersion({ core, exec });
  }

  core.info(`Current version: ${tag}`);

  // Trim leading 'v' from version output
  const version = tag.replace(/^v/, "");
  core.setOutput("version", version);
  core.setOutput("created", created.toString());
};

/**
 * Get the current tag and create the next tag by incrementing the version.
 *
 * Checks if the current commit is already tagged to avoid duplicate tags.
 * Uses the configured strategy (rollover or conventional-commits) to determine
 * the next version, then creates and pushes the tag.
 *
 * @returns {Object} { tag: string, created: boolean }
 */
async function getAndCreateNextTag({ core, exec }) {
  const currentTag = await getLatestTag({ core, exec });

  // Check if current commit already has a tag
  const hasTag = await checkIfCurrentCommitHasTag({ exec, tag: currentTag });
  if (hasTag) {
    core.info(
      `Current commit already tagged as ${currentTag}, skipping create`,
    );
    return { tag: currentTag, created: false };
  }

  const strategy = process.env.strategy || "rollover";
  let nextTag;

  if (strategy === "conventional-commits") {
    nextTag = await incrementVersionConventionalCommits({
      core,
      exec,
      tag: currentTag,
    });
  } else {
    nextTag = incrementVersionRollover({ tag: currentTag });
  }

  core.info(`Creating new tag: ${nextTag}`);

  await configureGit({ exec });
  await createTag({ exec, tag: nextTag });
  await pushTag({ exec, tag: nextTag });

  return { tag: nextTag, created: true };
}

/**
 * Get the latest git tag in short form (without commit count/hash).
 *
 * Uses `git describe --abbrev=0` to get just the tag name.
 * Defaults to v0.0.0 if no tags exist.
 *
 * @returns {string} The latest tag (e.g., "v0.1.4")
 */
async function getLatestTag({ core, exec }) {
  const result = await exec.getExecOutput(
    "git",
    ["describe", "--match", "v[0-9]*", "--abbrev=0"],
    { ignoreReturnCode: true },
  );

  if (result.exitCode !== 0) {
    core.info("No existing tags found, defaulting to v0.0.0");
    return "v0.0.0";
  }

  return result.stdout.trim();
}

/**
 * Get the current version for read-only mode (create=false).
 *
 * Returns different formats depending on the current commit state:
 * - If commit is tagged: Returns the clean tag (e.g., "v0.1.4")
 * - If commit is not tagged but tags exist: Returns dev version with base tag, timestamp (UTC), and
 *   hash (e.g., "v0.1.4-dev.20260628143052.abc123def456")
 * - If no tags exist: Returns v0.0.0 with timestamp (UTC) and hash
 *   (e.g., "v0.0.0-dev.20260628143052.abc123def456")
 *
 * Timestamp format: YYYYMMDDHHmmss in UTC (controlled by TZ=UTC environment variable)
 *
 * @returns {string} The current version string
 */
async function getCurrentVersion({ core, exec }) {
  // Check if current commit has a tag
  const exactMatchResult = await exec.getExecOutput(
    "git",
    ["describe", "--exact-match", "--match", "v[0-9]*"],
    { ignoreReturnCode: true },
  );

  if (exactMatchResult.exitCode === 0) {
    // Current commit is tagged, return the tag
    core.info("Current commit is tagged");
    return exactMatchResult.stdout.trim();
  }

  // Current commit is not tagged, use dev format
  // Try to get version with describe format
  let result = await exec.getExecOutput(
    "git",
    [
      "--no-pager",
      "show",
      "--quiet",
      "--abbrev=12",
      "--date=format-local:%Y%m%d%H%M%S",
      "--format=%(describe:abbrev=0)-dev.%cd.%h",
    ],
    { ignoreReturnCode: true },
  );

  if (result.exitCode === 0) {
    return result.stdout.trim();
  }

  // No tags exist, construct version manually with 0.0.0 prefix
  core.info("No existing tags found, constructing version with v0.0.0 prefix");
  result = await exec.getExecOutput(
    "git",
    [
      "--no-pager",
      "show",
      "--quiet",
      "--abbrev=12",
      "--date=format-local:%Y%m%d%H%M%S",
      "--format=%cd.%h",
    ],
    { ignoreReturnCode: true },
  );

  if (result.exitCode !== 0) {
    return "v0.0.0-dev";
  }

  const timestampAndHash = result.stdout.trim();
  return `v0.0.0-dev.${timestampAndHash}`;
}

/**
 * Check if the current commit already has the given tag.
 *
 * Used to prevent creating duplicate tags on the same commit.
 *
 * @param {string} tag - The tag to check for
 * @returns {boolean} True if the tag exists on the current commit
 */
async function checkIfCurrentCommitHasTag({ exec, tag }) {
  const result = await exec.getExecOutput(
    "git",
    ["describe", "--exact-match", "--match", tag],
    { ignoreReturnCode: true },
  );

  return result.exitCode === 0;
}

/**
 * Parse a version tag into [major, minor, patch] array.
 *
 * Strips the 'v' prefix and any git describe suffixes (e.g., "-4-gabcdef").
 * Ensures the result always has at least 3 parts (major.minor.patch).
 *
 * @param {string} tag - Version tag to parse (e.g., "v0.1.4" or "v0.1.4-5-gabcdef")
 * @returns {number[]} Array of [major, minor, patch]
 */
function parseVersion(tag) {
  // Remove 'v' prefix if present
  const version = tag.replace(/^v/, "");

  // Split by '-' to remove any git describe suffix (e.g., "1.2.3-4-gabcdef")
  const baseVersion = version.split("-")[0];

  // Split into parts
  const parts = baseVersion.split(".").map((p) => parseInt(p));

  // Ensure we have at least major.minor.patch format
  while (parts.length < 3) {
    parts.push(0);
  }

  return parts;
}

/**
 * Format version parts into a tag string with 'v' prefix.
 *
 * @param {number[]} parts - Array of [major, minor, patch]
 * @returns {string} Formatted tag (e.g., "v0.1.5")
 */
function formatVersion(parts) {
  return "v" + parts.map((p) => p.toString()).join(".");
}

/**
 * Increment version using the rollover strategy.
 *
 * Increments the patch version. When the patch version would reach the configured
 * rollover-max value, increments the minor version and resets patch to 0.
 * The major version is never modified (ZeroVer philosophy).
 *
 * Example with rollover-max=100:
 * - v0.5.50 → v0.5.51
 * - v0.5.99 → v0.6.0
 *
 * @param {string} tag - Current version tag
 * @returns {string} Next version tag
 */
function incrementVersionRollover({ tag }) {
  const rolloverMax = parseInt(process.env.rollover_max || "100");
  const parts = parseVersion(tag);

  // Increment patch version
  parts[2] = parts[2] + 1;

  // Check for rollover
  if (parts[2] >= rolloverMax) {
    parts[1] = parts[1] + 1; // Increment minor
    parts[2] = 0; // Reset patch
  }

  return formatVersion(parts);
}

/**
 * Increment version using the conventional commits strategy.
 *
 * Analyzes commit messages since the last tag following the Conventional Commits specification
 * to determine the appropriate version increment:
 *
 * - Breaking changes (BREAKING CHANGE footer or '!' after type/scope) → increment minor, reset
 *   patch to 0
 * - New features (feat: prefix) → increment minor, reset patch to 0
 * - Other commits (fix:, chore:, docs:, etc.) → increment patch
 *
 * The major version is never modified (ZeroVer philosophy).
 *
 * Examples:
 * - "feat: add new API" → v0.5.10 → v0.6.0
 * - "fix: resolve bug" → v0.5.10 → v0.5.11
 * - "feat!: breaking change" → v0.5.10 → v0.6.0
 *
 * @param {string} tag - Current version tag
 * @returns {string} Next version tag
 */
async function incrementVersionConventionalCommits({ core, exec, tag }) {
  const parts = parseVersion(tag);

  // Get commits since the last tag
  const commits = await getCommitsSinceTag({ exec, tag });

  if (commits.length === 0) {
    core.info("No new commits since last tag, incrementing patch version");
    parts[2] = parts[2] + 1;
    return formatVersion(parts);
  }

  // Analyze commits for conventional commit patterns
  let hasBreakingChange = false;
  let hasFeature = false;

  for (const commit of commits) {
    // Check for breaking change indicator (BREAKING CHANGE: or ! after type/scope)
    if (
      commit.includes("BREAKING CHANGE:") ||
      commit.includes("BREAKING-CHANGE:")
    ) {
      hasBreakingChange = true;
      break;
    }

    // Check for breaking change with ! syntax: type(scope)!: or type!:
    const breakingPattern = /^[a-z]+(\([^)]*\))?!:/;
    if (breakingPattern.test(commit)) {
      hasBreakingChange = true;
      break;
    }

    // Check for feature commits
    const featPattern = /^feat(\([^)]*\))?:/;
    if (featPattern.test(commit)) {
      hasFeature = true;
    }
  }

  // Determine version increment
  if (hasBreakingChange || hasFeature) {
    core.info(
      `Conventional commits analysis: ${hasBreakingChange ? "breaking change" : "feature"} detected, incrementing minor version`,
    );
    parts[1] = parts[1] + 1; // Increment minor
    parts[2] = 0; // Reset patch
  } else {
    core.info(
      "Conventional commits analysis: no breaking changes or features, incrementing patch version",
    );
    parts[2] = parts[2] + 1; // Increment patch
  }

  return formatVersion(parts);
}

/**
 * Get commit messages since the last tag.
 *
 * Retrieves the subject line (first line) of each commit between the specified tag and HEAD.
 *
 * @param {string} tag - The tag to compare against
 * @returns {string[]} Array of commit message subjects
 */
async function getCommitsSinceTag({ exec, tag }) {
  const result = await exec.getExecOutput(
    "git",
    ["log", `${tag}..HEAD`, "--pretty=format:%s"],
    { ignoreReturnCode: true },
  );

  if (result.exitCode !== 0 || !result.stdout.trim()) {
    return [];
  }

  return result.stdout
    .trim()
    .split("\n")
    .filter((line) => line.length > 0);
}

/**
 * Configure the git client with the specified user name and email.
 *
 * Uses values from environment variables (user_name and user_email) if provided,
 * otherwise defaults to the github-actions bot identity.
 */
async function configureGit({ exec }) {
  const userName = process.env.user_name || "github-actions[bot]";
  const userEmail =
    process.env.user_email ||
    "41898282+github-actions[bot]@users.noreply.github.com";

  await exec.exec("git", ["config", "user.name", userName]);
  await exec.exec("git", ["config", "user.email", userEmail]);
}

/**
 * Create an annotated git tag.
 *
 * @param {string} tag - The tag name to create (e.g., "v0.1.5")
 */
async function createTag({ exec, tag }) {
  await exec.exec("git", [
    "tag",
    "--annotate",
    tag,
    "--message",
    `Version ${tag}`,
  ]);
}

/**
 * Push the created tag to the remote repository.
 *
 * @param {string} tag - The tag name to push (e.g., "v0.1.5")
 */
async function pushTag({ exec, tag }) {
  await exec.exec("git", ["push", "origin", tag]);
}
