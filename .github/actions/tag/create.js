/*global module*/

/**
 *
 */
module.exports = async ({ exec }) => {
  await configure({ exec });

  // TODO: detect when the version already exists.
  const version = process.env.version;
  await createAndPush({ exec, version });
};

/**
 * Create the given version as an annotated git tag and push the result.
 */
async function createAndPush({ exec, version }) {
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
async function configure({ exec }) {
  await exec.exec("git", ["config", "user.name", "github-actions[bot]"]);
  await exec.exec("git", ["config", "user.email", "41898282+github-actions[bot]@users.noreply.github.com"]);
}
