/*global module, process*/

/**
 * Calculate the next version.
 *
 * If no tags are present, the minimum input is used. If the minimum input is greater than the
 * calculated next tag, said minimum input is used again.
 *
 * The next tag is calculated by incrementing the last segment of the version. For example, "1.2.3"
 * becomes "1.2.4". If the rollover input is enabled, then the incrementing causes the next last
 * segment to be increased if incrementing that last segment would cause the limit to be reached.
 * For example, with a limit of "100", "1.2.99" becomes "1.3.0". Similarly, "1.99.99" becomes
 * "2.0.0".
 */
module.exports = async ({ core, exec }) => {
  const minimum = process.env.minimum;
  const rollover = process.env.rollover == "true";
  const limit = parseInt(process.env.limit);

  const { describe, ok } = await current({ exec });

  if (!ok || describe == minimum) {
    core.info(`no existing tags, defaulting to minimum: ${minimum}`);
    core.setOutput("next", minimum);
    return;
  }

  if (!hasSuffix({ describe })) {
    core.info(`rebuild detected, using current tag: ${describe}`);
    core.setOutput("next", describe);
    return;
  }

  const next = increment({ describe, rollover, limit });
  core.info(`calculated next: ${next}`);

  if (compare({ minimum, next }) <= 0) {
    core.info(`minimum (${minimum}) >= next (${next}), using minimum`);
    core.setOutput("next", minimum);
    return;
  }

  core.info(`next (${next}) > minimum (${minimum}), using next`);
  core.setOutput("next", next);
};

/**
 * Get the current version.
 */
async function current({ exec }) {
  const result = await exec.getExecOutput(
    "git",
    ["describe", "--match", "v[0-9]*"],
    { ignoreReturnCode: true },
  );

  if (result.exitCode != 0) {
    return { ok: false };
  }

  const describe = result.stdout.replace(/^v/, "").trim();
  return { describe, ok: true };
}

/**
 * Test whether the describe output has a "-<count>-g<commit>" suffix.
 */
function hasSuffix({ describe }) {
  return describe.split("-", 1)[0] != describe;
}

/**
 * Increment a version.
 */
function increment({ describe, rollover, limit }) {
  const parts = describe
    .split("-", 1)[0]
    .split(".")
    .map((p) => parseInt(p));

  for (let i = parts.length - 1; i >= 0; i--) {
    parts[i]++;

    if (!rollover || parts[i] < limit) {
      break;
    }

    parts[i] = 0;
  }

  return parts.map((p) => p.toString()).join(".");
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

  const m = minimum.split(".").map((p) => parseInt(p));
  const n = next.split(".").map((p) => parseInt(p));

  // TODO: decide what to do if the lengths are different
  for (let i = 0; i < m.length; i++) {
    if (m[i] > n[i]) {
      return -1;
    }

    if (m[i] < n[i]) {
      return 1;
    }
  }

  return 0;
}
