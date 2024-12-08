function show_hide_password(target) {
	var input = document.getElementById('password-input');
	if (input.getAttribute('type') == 'password') {
		target.classList.add('view');
		input.setAttribute('type', 'text');
	} else {
		target.classList.remove('view');
		input.setAttribute('type', 'password');
	}
	return false;
}
function togglePasswordVisibility(element, inputId) {
	const input = document.getElementById(inputId);
	if (input.type === "password") {
		input.type = "text";
		element.classList.add("visible");
	} else {
		input.type = "password";
		element.classList.remove("visible");
	}
	return false;
}