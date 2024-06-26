import store from "../state/store";
import {set, unset} from "../state/userSlice";

import axios from "axios";
import qs from "qs";

function first_login(email: string, password: string): Promise<boolean> {
	const login_data = {
		"email": email,
		"password": password
	};

	return axios({
			method: "post",
			headers: {'content-type': 'application/x-www-form-urlencoded'},
			withCredentials: true,
			data: qs.stringify(login_data),
			url: "/api/login",
	}).then(res => {
		if ("data" in res === false)
			throw "unexpected response";

		const ans = res.data;
		if (ans.success) {
			store.dispatch(set({
				type: "login",
				id: ans.userId,
				alias: ans.alias,
				kind: ans.kind
			}));
			return true;
		}
		return false;
	});
}

function continue_login(): Promise<boolean> {
	return axios({
		method: "post",
		url: "/api/continue-login",
		withCredentials: true
	}).then(res => {
		if ("data" in res === false)
			throw "unexpected response";

		const ans = res.data;
		if (ans.success) {
			store.dispatch(set({
				type: "login",
				id: ans.userId,
				alias: ans.alias,
				kind: ans.kind
			}));
			return true;
		}
		return false;
	});
}

async function register(registerData: {
	alias: string,
	email: string,
	phone: string,
	birthDate: string,
	password: string,
	kind: string
}): Promise<boolean> {

	const res = await axios({
		method: "post",
		headers: {'content-type': 'application/x-www-form-urlencoded'},
		withCredentials: true,
		data: qs.stringify(registerData),
		url: "/api/register",
	});

	const data = res?.data;
	if (data.success) {
		store.dispatch(set({
			type: "login",
			id: data.userId,
			alias: data.alias,
			kind: registerData.kind
		}));
		return true;
	}
	return false;
}

async function logout(): Promise<null> {
	store.dispatch(unset());
	return axios({
		method: "post",
		url: "/api/logout",
		withCredentials: true
	});
}

export { continue_login, first_login, logout, register};

