import {
	handleDownloadHeaderClick,
	handleUploadHeaderClick,
} from './listeners.js';

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
	input.capture = true;
	input.accept = 'video/*';
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

	const fileContainer = document.createElement('div');
	fileContainer.id = 'files-tab';
	fileContainer.classList.add('files-container');

	uploadTab.appendChild(fileContainer);

	document.querySelector('body').appendChild(uploadTab);
}
