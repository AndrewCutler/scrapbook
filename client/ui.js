import {
	renderUploadButton,
	handleDownloadHeaderClick,
	handleUploadHeaderClick,
} from './listeners.js';

export function renderError(type) {
	if (['login', 'file-format'].includes(type)) {
		const error = document.createElement('div');
		error.id = `${type}-error`;

		const message = document.createElement('article');
		message.classList.add('message', 'is-danger');
		const messageHeader = document.createElement('div');
		messageHeader.classList.add('message-header');
		messageHeader.innerHTML = '<p>Error</p>';
		message.appendChild(messageHeader);
		const messageBody = document.createElement('div');
		messageBody.classList.add('message-body');
		switch (type) {
			case 'login':
				messageBody.innerText = 'Invalid username/password.';
				break;
			case 'file-format':
				messageBody.innerText = 'Invalid file.';
				break;
			default:
				break;
		}
		message.appendChild(messageBody);

		error.appendChild(message);

		document.body.appendChild(error);

		setTimeout(function () {
			error.remove();
		}, 8000);
	}
}

export function renderTabHeaders() {
	const tabContainer = document.createElement('div');
	tabContainer.classList.add('tabs');

	const tabUl = document.createElement('ul');
	tabUl.classList.add('tab-headers');

	const uploadLi = document.createElement('li');
	uploadLi.id = 'upload-header';
	uploadLi.classList.add('is-active');
	const uploadA = document.createElement('a');
	uploadA.innerText = 'Upload';
	uploadLi.appendChild(uploadA);
	tabUl.appendChild(uploadLi);

	const downloadLi = document.createElement('li');
	downloadLi.id = 'download-header';
	const downloadA = document.createElement('a');
	downloadA.innerText = 'Download';
	downloadLi.appendChild(downloadA);
	tabUl.appendChild(downloadLi);

	tabContainer.appendChild(tabUl);

	document.querySelector('body').appendChild(tabContainer);

	handleUploadHeaderClick();
	handleDownloadHeaderClick();
}

export function renderUploadTab() {
	document.querySelector('#files-tab')?.remove();
	const uploadTab = document.createElement('div');
	uploadTab.id = 'upload-tab';

	const header = document.createElement('h2');
	header.innerText = 'Choose videos to upload';
	uploadTab.appendChild(header);

	const uploadContainer = document.createElement('div');
	uploadContainer.classList.add('upload-container');
	const inputContainer = document.createElement('div');
	inputContainer.id = 'input-container';
	const input = document.createElement('input');
	input.type = 'file';
	input.name = 'upload';
	input.id = 'upload';
	input.multiple = true;
	// input.capture = true;
	// input.accept = 'image/*';
	// input.addEventListener('change', function (e) {
	// 	return validateFile(file, ['.mp4']);
	// });
	inputContainer.appendChild(input);
	uploadContainer.appendChild(inputContainer);

	const previewUl = document.createElement('ul');
	previewUl.id = 'previews';
	previewUl.style = 'display: inline-block';
	uploadContainer.appendChild(previewUl);
	uploadContainer.appendChild(document.createElement('br'));
	const uploadButton = document.createElement('button');
	uploadButton.id = 'submit';
	uploadButton.classList.add('button');
	uploadButton.type = 'button';
	uploadButton.innerText = 'Upload';
	uploadContainer.appendChild(uploadButton);

	uploadTab.appendChild(uploadContainer);

	document.querySelector('body').appendChild(uploadTab);

	renderUploadButton();
}

export function renderDownloadTab() {
	document.querySelector('#upload-tab')?.remove();
	const fileContainer = document.createElement('div');
	fileContainer.id = 'files-tab';

	document.querySelector('body').appendChild(fileContainer);
}
