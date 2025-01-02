/**
 *
 */
module.exports = async ({core, exec}) => {
	const version = process.env.version;
	if (version != 'latest') {
		core.setOutput('version', version);
	}

	const pkg = process.env.package;

	let mod = pkg;

	while (true) {
		let output = '';

		const options = {
			ignoreReturnCode: true,
			listeners: {
				stdout: (data) => {
					output += data.toString();
				},
			},
		};

		const code = await exec.exec(
			'go',
			['list', '-m', '-versions', '-mod=readonly', '-json',  mod],
			options
		);

		if (code == 0) {
			const data = JSON.parse(output);
			core.setOutput('version', data.Versions[-1]);
			console.debug(data.Versions[-1]);
			return;
		}

		if (mod.lastIndexOf('/') == -1) {
			core.setFailed('failed to identify go module');
			return;
		}

		mod = mod.split('/').slice(0, -1).join('/');
	}
}
