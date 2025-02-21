/*global module*/

/**
 *
 */
module.exports = async ({ core, exec }) => {
  const minimum = process.env.minimum;
  const rollover = process.env.rollover == "true";
  const limit = parseInt(process.env.limit);

  const { describe, ok } = await current({ exec });

  if (!ok || describe == minimum) {
    core.setOutput("next", minimum);
    return;
  }

  const next = increment({ describe, rollover, limit });

  if (compare({minimum, next}) <= 0) {
    core.setOutput("next", minimum);
    return;
  }

  core.setOutput("next", next);
};

/**
 * Get the current version.
 */
async function current({ exec }) {
  const result = exec.getExecOutput(
    "git",
    ["describe", "--match", "'v[0-9]*"],
    { ignoreReturnCode: true },
  );

  if (current.exitCode != 0) {
    return { ok: false };
  }

  const describe = current.stdout.replace(/^v/, "");
  return { describe, ok: true };
}

/**
 * Compare 2 versions together, returning one of the following integers:
 *
 * - -1: If the minimum parameter is higher.
 * - 0: If the values are the same.
 * - 1: If the next parameter is higher.
 */
function compare({ minimum, next }) {
  if (minimum == next) {
    return 0;
  }

  const m = minimum.split(".").map(p => parseInt(p));
  const n = next.split(".").map(p => parseInt(p));

  // TODO: decide what to do if the lengths are different
  for (let i = 0; i < m.length; i++) {
    if (m[i] > n[i]) {
      return -1;
    }

    if (m[i] < n[i]) {
      return 1
    }
  }

  return 0;
}

/**
 * Increment a version.
 */
function increment({ describe, rollover, limit }) {
  const parts = version.split("-", 1)[0].split(".").map(p => parseInt(p));

  for (let i = parts.length - 1; i >= 0; i--) {
    parts[i]++;

    if (!rollover || parts[i] < limit) {
      break;
    }

    parts[i] = 0;
  }

  return parts.map(p => p.toString()).join(".");
}
