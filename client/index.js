window.onload = function () {
	const MAX_FILES = 4;

	async function downloadVideo(filename) {
		const response = await fetch(
			`http://10.0.0.73:8000/api/files/${filename}`,
		).catch(console.error);
		try {
			const data = await response.blob();
			var a = document.createElement('a');
			a.href = window.URL.createObjectURL(data);
			a.download = filename;
			a.click();
		} catch (err) {
			console.error(err);
		}
	}

	function renderFileData(thumbnails) {
		for (const curr of thumbnails) {
			const fileDiv = document.createElement('div');
			fileDiv.style.display = 'flex';
			fileDiv.style.flexDirection = 'column';

			const filenameDiv = document.createElement('div');
			filenameDiv.innerText = curr.Name;
			fileDiv.appendChild(filenameDiv);

			const imgDiv = document.createElement('div');
			const img = document.createElement('img');
			img.src = curr.Thumbnail;
			img.height = 200;
			img.width = 200;
			imgDiv.appendChild(img);
			fileDiv.appendChild(imgDiv);

			const downloadDiv = document.createElement('div');
			const downloadButton = document.createElement('button');
			downloadButton.type = 'button';
			downloadButton.innerText = 'Download video';
			downloadButton.onclick = function () {
				downloadVideo(curr.Name);
				console.log('download');
			};

			downloadDiv.appendChild(downloadButton);
			fileDiv.appendChild(downloadDiv);

			document.querySelector('#files').appendChild(fileDiv);
		}
	}

	async function getThumbnails() {
		const response = await fetch('http://10.0.0.73:8000/api/files').catch(
			console.error,
		);
		try {
			const thumbnails = await response.json();
			renderFileData(thumbnails);
		} catch (e) {
			console.error(e);
		}
	}

	(async () => await getThumbnails())();

	const previews = document.querySelector('#previews');

	const uploadButton = document.querySelector('#submit');
	uploadButton.style.display = 'none';
	uploadButton.addEventListener('click', async function () {
		// revoke object urls
		const formData = new FormData();
		for (let i = 0; i < input.files.length; i++) {
			const file = input.files[i];
			console.log({ i, file });
			formData.append('files', file);
		}

		await fetch('http://10.0.0.73:8000/api/save', {
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
	});
};
