function renderForm() {

}

export function renderLogin() {
	const loginContainer = document.createElement('div');
	loginContainer.id = 'login-container';
	const loginButton = document.createElement('button');
	loginButton.classList.add('button');
	loginButton.type = 'submit';
	loginButton.innerText = 'Login';
	loginContainer.appendChild(loginButton);

	document.querySelector('body').appendChild(loginContainer);


}
