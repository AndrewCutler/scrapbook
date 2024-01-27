import {
	createUploadButton,
	handleDownloadHeaderClick,
	handleUploadHeaderClick,
	uploadFile,
} from './listeners.js';

window.onload = function () {
	const BASEURL = 'http://10.0.0.73:8000/api';

	async function downloadVideo(filename) {
		const response = await fetch(`${BASEURL}/files/${filename}`).catch(
			console.error,
		);
		try {
			const data = await response.blob();
			const a = document.createElement('a');
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

			const filenameDiv = document.createElement('div');
			filenameDiv.innerText = curr.Name;
			fileDiv.appendChild(filenameDiv);

			const imgDiv = document.createElement('div');
			const img = document.createElement('img');
			img.classList.add('thumbnail');
			img.src = curr.Thumbnail;
			imgDiv.appendChild(img);
			fileDiv.appendChild(imgDiv);

			const downloadDiv = document.createElement('div');
			const downloadButton = document.createElement('button');
			downloadButton.type = 'button';
			downloadButton.classList.add('button');
			downloadButton.innerText = 'Download video';
			downloadButton.onclick = function () {
				downloadVideo(curr.Name);
			};

			downloadDiv.appendChild(downloadButton);
			fileDiv.appendChild(downloadDiv);

			document.querySelector('#files-tab').appendChild(fileDiv);
		}
	}

	async function getThumbnails() {
		const response = await fetch(`${BASEURL}/files`).catch(console.error);
		try {
			const thumbnails = await response.json();
			renderFileData(thumbnails);
		} catch (e) {
			console.error(e);
		}
	}

	(async () => await getThumbnails())();

	// const previews = document.querySelector('#previews');

	// function createPreviewElement(file, index) {
	// 	const { name, size } = file;

	// 	const previewListItem = document.createElement('li');
	// 	previewListItem.classList.add('preview-item');
	// 	previewListItem.id = `preview-${index}`;

	// 	const previewDescription = document.createElement('span');
	// 	previewDescription.textContent = `Name: ${name}; size: ${size}`;
	// 	previewListItem.appendChild(previewDescription);

	// 	const videoContainer = document.createElement('div');
	// 	videoContainer.classList.add('video-container');
	// 	const video = document.createElement('video');
	// 	video.src = URL.createObjectURL(file);
	// 	video.controls = true;
	// 	video.muted = true;
	// 	videoContainer.appendChild(video);

	// 	const removeButton = document.createElement('a');
	// 	removeButton.textContent = 'Remove';
	// 	removeButton.style.cursor = 'pointer';
	// 	removeButton.style.marginTop = '2px';
	// 	removeButton.classList.add('button');
	// 	removeButton.addEventListener('click', function () {
	// 		console.log('delete ', index);
	// 	});
	// 	videoContainer.appendChild(removeButton);

	// 	previewListItem.appendChild(videoContainer);
	// 	previews.appendChild(previewListItem);
	// }

	const input = document.querySelector('#upload');
	createUploadButton(input);
	input.addEventListener('change', uploadFile(uploadFile, input));

	handleUploadHeaderClick();
	handleDownloadHeaderClick();

	fetch('http://10.0.0.73:8000/api/test')
		.then(console.log)
		.catch(console.error);
};
