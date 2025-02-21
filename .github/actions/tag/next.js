/*global module*/

/**
 *
 */
module.exports = async ({ exec }) => {
  const fallback = makeFallback();

  const current = exec.getExecOutput(
    "git",
    ["describe", "--match", "'v[0-9]*"],
    { ignoreReturnCode: true },
  );

  if (current.exitCode != 0) {
    await create(exec, fallback);
    return;
  }

  const describe = current.stdout.replace(/^v/, "");
  const cmp = compare(fallback, describe);

  if (cmp == -1) {
    await create(exec, fallback);
    return;
  }

  if (cmp == 0) {
    return;
  }

  const next = increment(describe);
  await create(exec, next);

  return;
};

/**
 * Generate a fallback version in the format "YYYY.MM.0".
 */
function makeFallback() {
  const now = new Date(Date.now());

  const year = now.getUTCFullYear();
  const month = now.getUTCMonth() + 1;

  return `${year}.${month}.0`;
}

/**
 * Compare 2 versions together, returning one of the following integers:
 *
 * - 1: If the describe parameter is higher.
 * - 0: If the versions are the same.
 * - -1: If the fallback parameter is higher.
 */
function compare(fallback, describe) {
  const desc = describe.replace(/^v/, "");

  if (fallback == desc) {
    return 0;
  }

  const a = fallback.split(".");
  const b = desc.split("-", 1).split(".");

  for (let i = 0; i < a.length; i++) {
    // guard against describe being shorter than fallback
    if (i > b.length) {
      return -1;
    }

    const c = parseInt(a[i]);
    const d = parseInt(b[i]);

    if (c == d) {
      continue;
    }

    if (c > d) {
      return -1;
    }

    return 1;
  }

  return 1;
}

/**
 * Increment the last segment of the given version.
 */
function increment(version) {
  const parts = version.split("-", 1)[0].split(".");
  let last = parseInt(parts[parts.length - 1]);

  last += 1;
  parts[parts.length - 1] = last.toString();

  return parts.join(".");
}

/**
 * Create the given version as an annotated git tag.
 */
async function create(exec, version) {
  await configure(exec);

  await exec.exec("git", [
    "tag",
    "--annotate",
    `v${version}`,
    "--message",
    `Version ${version}`,
  ]);

  // TODO: push tag
  // await exec.exec("git", ["push", `v${version}`]);
}

/**
 *
 */
async function configure(exec) {
  await exec.exec("git", ["config", "user.name", "github-actions[bot]"]);
  await exec.exec("git", ["config", "user.email", "41898282+github-actions[bot]@users.noreply.github.com"]);
}
