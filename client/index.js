window.onload = function () {
	console.log('init');

	const input = document.querySelector('#upload');
	input.addEventListener('change', function () {
		console.log(input.files);
		if (input.files.length === 1) {
			const [file] = input.files;
			const formData = new FormData();
			formData.append('file',  file);
			fetch('http://localhost:8000/save', {
				method: 'POST',
				body: formData,
			});
		}
	});
};
