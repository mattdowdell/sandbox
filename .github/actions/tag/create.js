/*global module, process*/

/**
 * Create a tag for the given version.
 *
 * A "v" is prepending to support Go tagging conventions.
 */
module.exports = async ({ core, exec }) => {
  await configure({ exec });

  const version = process.env.version;

  await create({ core, exec, version });
  await push({ exec, version });
};

/**
 * Create the given version as an annotated git tag.
 */
async function create({ core, exec, version }) {
  const result = await exec.exec(
    "git",
    ["rev-list", "--max-count", "1", `v${version}`],
    { ignoreReturnCode: true },
  );

  if (result.exitCode == 0) {
    core.info(`tag v${version} exists, skipping create`);
    return;
  }

  await exec.exec("git", [
    "tag",
    "--annotate",
    `v${version}`,
    "--message",
    `Version ${version}`,
  ]);
}

/**
 * Push the created tag.
 */
async function push({ exec, version }) {
  await exec.exec("git", ["push", "origin", `v${version}`]);
}

/**
 * Configure the git client to take the identity of the github-actions bot.
 */
async function configure({ exec }) {
  await exec.exec("git", ["config", "user.name", "github-actions[bot]"]);
  await exec.exec("git", [
    "config",
    "user.email",
    "41898282+github-actions[bot]@users.noreply.github.com",
  ]);
}
