export const serverURL = 'http://localhost:3000';

export const getCookie = (name: string) => {
	const value = `; ${document.cookie}`;
	const parts = value.split(`; ${name}=`);
	return parts.length === 2 ? parts.pop().split(';').shift() : '';
};

export const getCurrentSemester = () => {
	const d = new Date();
	const month = d.getMonth();
	const year = d.getFullYear();

	let term: string;
	if (month <= 4) {
		term = 'Sp';
	} else if (month <= 7) {
		term = 'Su';
	} else {
		term = 'Fa';
	}

	return term.concat(`${year % 1000}`);
};

export const isUserLoggedIn = (): [boolean, string] => {
	const sessionToken = getCookie('session_token');
	const id = getCookie('net_id');

	if (sessionToken == '' || id == '') {
		return [false, ''];
	}

	return [true, id];
};

export const deleteCookie = (name: string) => {
	document.cookie = name + '=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
};

