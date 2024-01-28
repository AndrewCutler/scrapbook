function login() {}

function renderLoginError() {
	const loginError = document.createElement('div');
	loginError.id = 'login-error';

	const message = document.createElement('article');
	message.classList.add('message', 'is-danger');
	const messageHeader = document.createElement('div');
	messageHeader.classList.add('message-header');
	messageHeader.innerHTML = '<p>Error</p>';
	message.appendChild(messageHeader);
	const messageBody = document.createElement('div');
	messageBody.classList.add('message-body');
	messageBody.innerText = 'Invalid username/password.';
	message.appendChild(messageBody);

	loginError.appendChild(message);

	document.body.appendChild(loginError);
}

function renderLoginForm() {
	const formContainer = document.createElement('div');
	formContainer.classList.add('form-container');

	const form = document.createElement('form');
	form.id = 'login-form';

	const usernameField = document.createElement('div');
	usernameField.classList.add('field');
	const usernameLabel = document.createElement('div');
	usernameLabel.classList.add('label');
	usernameLabel.innerText = 'Username';
	usernameField.appendChild(usernameLabel);
	const usernameControl = document.createElement('div');
	usernameControl.classList.add('control');
	const usernameInput = document.createElement('input');
	usernameInput.type = 'text';
	usernameInput.name = 'username';
	usernameInput.id = 'username';
	usernameInput.placeholder = 'Username';
	usernameControl.appendChild(usernameInput);
	usernameField.appendChild(usernameControl);

	const passwordField = document.createElement('div');
	passwordField.classList.add('field');
	const passwordLabel = document.createElement('div');
	passwordLabel.classList.add('label');
	passwordLabel.innerText = 'Password';
	passwordField.appendChild(passwordLabel);
	const passwordControl = document.createElement('div');
	passwordControl.classList.add('control');
	const passwordInput = document.createElement('input');
	passwordInput.type = 'password';
	passwordInput.name = 'password';
	passwordInput.id = 'password';
	passwordInput.placeholder = 'Password';
	passwordControl.appendChild(passwordInput);
	passwordField.appendChild(passwordControl);

	form.appendChild(usernameField);
	form.appendChild(passwordField);

	formContainer.appendChild(form);

	document.body.appendChild(formContainer);
}

function renderLoginButton() {
	const loginButtonContainer = document.createElement('div');
	loginButtonContainer.id = 'login-button-container';
	const loginButton = document.createElement('input');
	loginButton.classList.add('button');
	loginButton.type = 'submit';
	loginButton.innerText = 'Login';
	loginButtonContainer.appendChild(loginButton);

	const form = document.getElementById('login-form');
	form.appendChild(loginButtonContainer);
	form.addEventListener('submit', function (e) {
		e.preventDefault();
		const formData = new FormData(e.target);

		if (
			formData.get('username')?.length < 3 ||
			formData.get('password')?.length < 3
		) {
			renderLoginError();

			return;
		}

		fetch('http://10.0.0.73:8000/api/login', {
			method: 'POST',
			body: JSON.stringify({
				username: 'username',
				password: 'password',
			}),
		});
	});
}

export function renderLogin() {
	renderLoginForm();
	renderLoginButton();
}
