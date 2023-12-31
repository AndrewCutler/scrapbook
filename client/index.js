window.onload = function () {
	console.log('init');

	const input = document.querySelector('#upload');
	input.addEventListener('change', function () {
		console.log(input.files);
	});
};
