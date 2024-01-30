import { createUploadButton, uploadFile } from './listeners.js';
import { renderLogin } from './login.js';
import { renderTabHeaders, renderUploadTab } from './ui.js';

export const Config = {
	baseUrl: '',
};

export async function getConfig() {
	try {
		const data = await fetch('./config.json');
		const json = await data.json();
		Config.baseUrl = json.baseUrl;
		console.log({ Config });
	} catch (err) {
		console.error(err);
	}
}

window.onload = function () {
	(async () => await getConfig())();

	renderLogin();

	// async function downloadVideo(filename) {
	// 	const response = await fetch(`${BASEURL}/files/${filename}`).catch(
	// 		console.error,
	// 	);
	// 	try {
	// 		const data = await response.blob();
	// 		const a = document.createElement('a');
	// 		a.href = window.URL.createObjectURL(data);
	// 		a.download = filename;
	// 		a.click();
	// 	} catch (err) {
	// 		console.error(err);
	// 	}
	// }

	// function renderFileData(thumbnails) {
	// 	for (const curr of thumbnails) {
	// 		const fileDiv = document.createElement('div');

	// 		const filenameDiv = document.createElement('div');
	// 		filenameDiv.innerText = curr.Name;
	// 		fileDiv.appendChild(filenameDiv);

	// 		const imgDiv = document.createElement('div');
	// 		const img = document.createElement('img');
	// 		img.classList.add('thumbnail');
	// 		img.src = curr.Thumbnail;
	// 		imgDiv.appendChild(img);
	// 		fileDiv.appendChild(imgDiv);

	// 		const downloadDiv = document.createElement('div');
	// 		const downloadButton = document.createElement('button');
	// 		downloadButton.type = 'button';
	// 		downloadButton.classList.add('button');
	// 		downloadButton.innerText = 'Download video';
	// 		downloadButton.onclick = function () {
	// 			downloadVideo(curr.Name);
	// 		};

	// 		downloadDiv.appendChild(downloadButton);
	// 		fileDiv.appendChild(downloadDiv);

	// 		document.querySelector('#files-tab').appendChild(fileDiv);
	// 	}
	// }

	// async function getThumbnails() {
	// 	const response = await fetch(`${BASEURL}/files`).catch(console.error);
	// 	try {
	// 		const thumbnails = await response.json();
	// 		renderFileData(thumbnails);
	// 	} catch (e) {
	// 		console.error(e);
	// 	}
	// }

	// (async () => await getThumbnails())();

	// const input = document.querySelector('#upload');
	// createUploadButton(input);
	// input.addEventListener('change', uploadFile(input));

	// fetch('http://10.0.0.73:8000/api/test')
	// 	.then(console.log)
	// 	.catch(console.error);
};
