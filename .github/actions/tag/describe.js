/*global module*/

/**
 *
 */
module.exports = async ({ core, exec }) => {
  // TODO: on error, fallback to 0.0.0-0-g<short commit>

  const short = await exec.getExecOutput("git", [
    "describe",
    "--always",
    "--match",
    "'v[0-9]*'",
  ]);
  core.setOutput("short", short.stdout);

  const long = await exec.getExecOutput("git", [
    "describe",
    "--always",
    "--long",
    "--match",
    "'v[0-9]*'",
  ]);
  core.setOutput("long", long.stdout);
};
