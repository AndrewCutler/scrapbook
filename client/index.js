import { uploadFile } from './listeners.js';
import { renderLogin } from './login.js';

export const Config = {
	baseUrl: '',
};

export async function getConfig() {
	try {
		const data = await fetch('./config.json');
		const json = await data.json();
		Config.baseUrl = json.baseUrl;
	} catch (err) {
		console.error(err);
	}
}

window.onload = function () {
	(async () => await getConfig())();

	renderLogin();
};
