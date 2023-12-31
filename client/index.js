window.onload = function () {
	console.log('init');

	const input = document.querySelector('#upload');
	input.addEventListener('change', async function () {
		console.log(input.files);
		if (input.files.length === 1) {
			const [file] = input.files;
			const formData = new FormData();
			formData.append('file', file);
			await fetch('http://10.0.0.73:8000/save', {
				method: 'POST',
				body: formData,
			}).catch(alert);
		}
	});
};
