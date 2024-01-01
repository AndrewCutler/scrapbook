window.onload = function () {
	const MAX_FILES = 4;

	const previews = document.querySelector('#previews');

	const uploadButton = document.querySelector('#submit');
	uploadButton.style.display = 'none';
	uploadButton.addEventListener('click', async function () {
		// revoke object urls
		const formData = new FormData();
		// for (const file of input.files) {
		// 	formData.append('file[]', file);
		// }
		for (let i = 0; i < input.files.length; i++) {
			const file = input.files[i];
			console.log({ i, file });
			formData.append('files', file);
		}

		await fetch('http://10.0.0.73:8000/save', {
			method: 'POST',
			body: formData,
		}).catch(console.error);
	});

	function createPreviewElement(file, index) {
		const { name, size } = file;

		const previewListItem = document.createElement('li');
		previewListItem.style.display = 'flex';
		previewListItem.style.flexDirection = 'column';
		previewListItem.style.padding = '6px';
		previewListItem.style.marginBottom = '8px';
		previewListItem.style.border = '1px solid #b5b5b5';
		previewListItem.style.borderRadius = '5px';
		previewListItem.id = `preview-${index}`;

		const previewDescription = document.createElement('span');
		previewDescription.textContent = `Name: ${name}; size: ${size}`;
		previewListItem.appendChild(previewDescription);

		const videoContainer = document.createElement('div');
		videoContainer.style.display = 'flex';
		videoContainer.style.flexDirection = 'column';
		const video = document.createElement('video');
		video.src = URL.createObjectURL(file);
		video.height = 200;
		video.width = 240;
		video.controls = true;
		video.muted = true;
		videoContainer.appendChild(video);

		const removeButton = document.createElement('a');
		removeButton.textContent = 'Remove';
		removeButton.style.cursor = 'pointer';
		removeButton.style.color = '#dd4';
		removeButton.addEventListener('click', function () {
			console.log('delete ', index);
		});
		videoContainer.appendChild(removeButton);

		previewListItem.appendChild(videoContainer);
		previews.appendChild(previewListItem);
	}

	const input = document.querySelector('#upload');
	input.addEventListener('change', function () {
		uploadButton.style.display = input.files.length > 0 ? 'block' : 'none';

		if (input.files.length > MAX_FILES) {
			uploadButton.disabled = true;

			const warning = document.createElement('small');
			warning.textContent =
				'A maximum of 5 files may be uploaded at once';
			warning.style.color = '#e66';
			document.querySelector('#input-container').appendChild(warning);

			return;
		}

		for (let i = 0; i < input.files.length; i++) {
			file = input.files[i];
			createPreviewElement(file, i);
		}

		// console.log(input.files);
		// if (input.files.length === 1) {
		// 	const [file] = input.files;
		// 	const formData = new FormData();
		// 	formData.append('file', file);

		// 	await fetch('http://10.0.0.73:8000/save', {
		// 		method: 'POST',
		// 		body: formData,
		// 	}).catch(console.error);
		// }
	});
};
