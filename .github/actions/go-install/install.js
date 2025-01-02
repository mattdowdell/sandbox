/**
 *
 */
module.exports = async ({exec}) {
	const package = process.env.package;
	const version = process.env.version;

	await exec.exec(['go', ['install', `${package}@${version}`]])
}
