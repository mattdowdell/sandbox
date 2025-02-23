/*global module*/

/**
 * Get the output of "git describe" and "git describe --long" without the "v" prefix.
 *
 * For pull requests, the output is identical. For default branch builds, short will omit the
 * "-<count>-g<commit>" suffix.
 */
module.exports = async ({ core, exec }) => {
  const short = await exec.getExecOutput(
    "git",
    ["describe", "--match", "v[0-9]*"],
    { ignoreReturnCode: true },
  );

  if (short.exitCode != 0) {
    const fallback = await makeFallback({ exec });
    core.info(`no tags found, using fallback: ${fallback}`);

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

  core.setOutput("short", short.stdout.replace(/^v/, "").trim());
  core.setOutput("long", long.stdout.replace(/^v/, "").trim());
};

/**
 * Create a fallback describe output for when no tags have been created.
 *
 * It is assumed that this will only be used for when no tags have been created, such as for the
 * very first pull request.
 */
async function makeFallback({ exec }) {
  const commit = await exec.getExecOutput("git", [
    "rev-parse",
    "--short",
    "HEAD",
  ]);

  return `0.0.0-0-g${commit.stdout.trim()}`;
}
