import { uploadFile } from './listeners.js';
import { renderLogin } from './login.js';
import { renderApp } from './ui.js';

export const Config = {
	baseUrl: '',
};

export async function initialize() {
	try {
		const data = await fetch('./config.json');
		const json = await data.json();
		Config.baseUrl = json.baseUrl;
		fetch(`${Config.baseUrl}/test`)
			.then(function (response) {
				if (response.status == 200) {
					renderApp();
				} else {
					renderLogin();
				}
			})
			.catch(console.error);
	} catch (err) {
		console.error(err);
	}
}

window.onload = function () {
	(async () => await initialize())();

	// renderLogin();
};
