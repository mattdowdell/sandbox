/*global module*/

/**
 *
 */
module.exports = async ({ core, exec }) => {
  const short = await exec.getExecOutput(
    "git",
    ["describe", "--match", "v[0-9]*"],
    { ignoreReturnCode: true },
  );

  if (short.exitCode != 0) {
    const fallback = await makeFallback({ exec });
    console.debug(`no tags found, using fallback: ${fallback}`);

    core.setOutput("short", fallback);
    core.setOutput("long", fallback);
    return;
  }

  const long = await exec.getExecOutput("git", [
    "describe",
    "--long",
    "--match",
    "v[0-9]*",
  ]);

  core.setOutput("short", short.stdout.replace(/^v/, ""));
  core.setOutput("long", long.stdout.replace(/^v/, ""));
};

/**
 *
 */
async function makeFallback({ exec }) {
  const commit = await exec.getExecOutput("git", [
    "rev-parse",
    "--short",
    "HEAD",
  ]);

  return `0.0.0-0-g${commit.stdout}`;
}
