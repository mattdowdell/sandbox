/**
 *
 */
module.exports = async ({exec}) => {
	const pkg = process.env.package;
	const version = process.env.version;

	await exec.exec('go', ['install', `${pkg}@${version}`])
}
