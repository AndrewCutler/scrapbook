import { Config } from './index.js';
import { renderDownloadTab } from './ui.js';

const MAX_FILES = 4;

export function createUploadButton(input) {
	const uploadButton = document.querySelector('#submit');
	uploadButton.style.display = 'none';
	uploadButton.classList.add('button');
	uploadButton.addEventListener('click', async function () {
		// revoke object urls
		const formData = new FormData();
		for (let i = 0; i < input.files.length; i++) {
			const file = input.files[i];
			console.log({ i, file });
			formData.append('files', file);
		}

		await fetch(`${Config.baseUrl}/save`, {
			method: 'POST',
			body: formData,
		}).catch(console.error);
	});
}

function createPreviewElement(file, index) {
	const { name, size } = file;

	const previewListItem = document.createElement('li');
	previewListItem.classList.add('preview-item');
	previewListItem.id = `preview-${index}`;

	const previewDescription = document.createElement('span');
	previewDescription.textContent = `Name: ${name}; size: ${size}`;
	previewListItem.appendChild(previewDescription);

	const videoContainer = document.createElement('div');
	videoContainer.classList.add('video-container');
	const video = document.createElement('video');
	video.src = URL.createObjectURL(file);
	video.controls = true;
	video.muted = true;
	videoContainer.appendChild(video);

	const removeButton = document.createElement('a');
	removeButton.textContent = 'Remove';
	removeButton.style.cursor = 'pointer';
	removeButton.style.marginTop = '2px';
	removeButton.classList.add('button');
	removeButton.addEventListener('click', function () {
		console.log('delete ', index);
	});
	videoContainer.appendChild(removeButton);

	previewListItem.appendChild(videoContainer);

	const previews = document.querySelector('#previews');
	previews.appendChild(previewListItem);
}

export function uploadFile(input) {
	return function () {
		const uploadButton = document.querySelector('#submit');
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
			const file = input.files[i];
			createPreviewElement(file, i);
		}
	};
}

function clearActiveHeaders() {
	const headers = document.querySelectorAll('.tab-headers li');
	headers.forEach((header) => header.classList.remove('is-active'));
}

function setActiveHeader(header) {
	clearActiveHeaders();
	header.classList.add('is-active');
}

export function handleUploadHeaderClick() {
	const uploadHeader = document.querySelector('#upload-header');
	uploadHeader.addEventListener('click', function () {
		setActiveHeader(uploadHeader);
		document.querySelector('#files-tab').style.display = 'none';
		document.querySelector('#upload-tab').style.display = 'block';
	});
}

export function handleDownloadHeaderClick() {
	const downloadHeader = document.querySelector('#download-header');
	downloadHeader.addEventListener('click', function () {
		setActiveHeader(downloadHeader);
		renderDownloadTab();
		document.querySelector('#files-tab').style.display = 'grid';
		document.querySelector('#upload-tab').style.display = 'none';
		(async () => await getFiles())();
	});
}

async function getFiles() {
	const response = await fetch(`${Config.baseUrl}/files`).catch(
		console.error,
	);
	try {
		const thumbnails = await response.json();
		renderFileData(thumbnails);
	} catch (e) {
		console.error(e);
	}
}

async function downloadVideo(filename) {
	const response = await fetch(`${Config.baseUrl}/files/${filename}`).catch(
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
	console.log({ thumbnails });
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
