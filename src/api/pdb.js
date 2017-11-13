/**
 *  query data from backend
 *  
 */

import axios from 'axios'
import store from '../utils/store';

const http = axios.create({
	baseURL: '/api',	// base api
	timeout: 5000,		// request timeout
	// `xsrfCookieName` is the name of the cookie to use as a value for xsrf token
	//xsrfCookieName: 'XSRF-TOKEN', // default
	// `xsrfHeaderName` is the name of the http header that carries the xsrf token value
	//xsrfHeaderName: 'X-XSRF-TOKEN', // default
});

// request interceptor
http.interceptors.request.use(config => {
	let token = store.getters.token
	if (token) {
		// set token in request header
		config.headers['X-Token'] = token; 
	}
	return config;
  }, 
  error => {
	// Do something with request error
	console.log(`axois request error : ${error}`);
	Promise.reject(error);
});

// fake build data set
const fakeBuilds = [{
		id: "0000000002",
		date: "2017-10-22 05:16:24",
		branch: "UDPv6.5U2",
		version: "4175.2-646",
		comment: "2017-10-22_05:16:23"
	},
	{
		id: "0000000001",
		date: "2017-10-21 17:17:54",
		branch: "UDPv6.5U2",
		version: "4175.2-645",
		comment: "2017-10-21_17:17:51"
	}
];

// fake branch data set
const fakeBranchs = [{
		buildName: "UDP_6_5_U2",
		storeName: "UDPv6.5U2",
		buildPath: "Z:\\UDP_6_5_U2\\Release",
		storePath: "c:\\Go\\huan\\src\\github.com\\adyzng\\GoSymbols\\testdata\\UDPv6.5U2",
		updateDate: "2017-10-22 15:31:00",
		latestBuild: "4175.2-646",
		buildsCount: 9,
	},
	{
		buildName: "UDP_6_5_U1",
		storeName: "UDPv6.5U1",
		buildPath: "Z:\\UDP_6_5_U1\\Release",
		storePath: "c:\\Go\\huan\\src\\github.com\\adyzng\\GoSymbols\\testdata\\UDPv6.5U1",
		updateDate: "2017-10-22 15:31:28",
		latestBuild: "4175.1-385",
		buildsCount: 65,
	}
];


export default {
	getFakeBranchs() {
		return [...fakeBranchs, ...fakeBranchs]
	},

	getFakeBuilds() {
		let d = []
		for (let i in 20) {
			d.concat(...fakeBuilds)
		}
		return d
	},

	getFakeMessages(cb) {
		const msgs = [ {
			succeed: true,
			branch: "UDPv6.5U2",
			updateDate: "2017-10-22 15:31:00",
		}, {
			succeed: true,
			branch: "UDPv6.5U1",
			updateDate: "2017-10-23 08:21:12",
		}, {
			succeed: false,
			branch: "UDPv6",
			updateDate: "",
		}];
		if (cb) {
			cb(msgs);
		}
	},
	
	getTodayMessages(cb) {
		return http.get('/messages').then(resp => {
			let res = resp.data
			if (Array.isArray(res.data)) {
				cb(res.data)
			} else {
				cb([])
			}
		})
		.catch(err => {
			cb([])
			console.log("getTodayMessages failed:", err);
		})
	},

	fetchBranchs(cb) {
		if (!cb) {
			return
		}
		return http.get("/branches").then(resp => {
				let res = resp.data;
				let data = [];
				if (res.data) {
					data = [...res.data.branchs];
					data.sort((a,b) => a.updateDate < b.updateDate);
				} else {
					console.log("fetchBranchs empty:", res);
				}
				cb(data);
			})
			.catch(error => {
				console.log("fetchBranchs failed:", error);
				cb([]);
			});
	},

	fetchBuilds(branch, cb) {
		if (!branch || !cb) {
			return
		}
		return http.get(`/branches/${branch}`).then(resp => {
				let res = resp.data;
				let	data = [];
				if (res.data) {
					data = [...res.data.builds];
					data.sort((a, b) => b.id - a.id);
				} else {
					console.log("fetchBuilds empty:", res)
				}
				cb(data);
			})
			.catch(err => {
				console.log("fetchBuilds failed:", err);
				cb([]);
			});
	},

	fetchSymbols(branch, build, cb) {
		if (!branch || !build || !cb) {
			return
		}
		return http.get(`/branches/${branch}/${build}`).then(resp => {
				let res = resp.data;
				let data = []
				if (res.data) {
					data = [...res.data.symbols];
					data.sort((a, b) => b.id - a.id);
				} else {
					console.log("fetchSymbols empty:", res)
				}
				cb(data);
			})
			.catch(err => {
				console.log("fetchSymbols failed:", err);
				cb([]);
			});
	},

	deleteBranch(branch, cb) {
		if (!branch) {
			return
		}
		return http.delete(`/branches/${branch.storeName}`).then(resp => {
			if (resp.data) {
				cb && cb(resp.data);
			}
		})
		.catch(err => {
			console.log(`deleteBranch error ${err}.`);
			cb && cb(err);
		})
	},

	validateBranch(branch, cb) {
		if (!branch) {
			return;
		}
		return http.post(`/branches/check`, JSON.stringify(branch)).then(resp => {
			if (resp.data) {
				cb && cb(resp.data);
			}
		})
		.catch(err => {
			console.log(`validateBranch error ${err}.`);
			cb && cb(err);
		})
	},

	modifyBranch(branch, cb) {
		if (!branch) {
			return;
		}
		return http.post(`/branches/modify`, JSON.stringify(branch)).then(resp => {
			if (resp.data) {
				cb && cb(resp.data);
			}
		})
		.catch(err => {
			console.log(`modifyBranch error ${err}.`);
			cb && cb(err);
		})
	},

	fetchProfile(cb) {
		return http.get(`/user/profile`).then(resp => {
			cb && cb(resp.data)
			console.log(`user ${JSON.stringify(resp.data.data)}`)
		})
		.catch(err => {
			console.log(`get user profile error ${err}`)
			cb && cb({})
		})
	},

	userLogout(cb) {
		return http.get('/auth/logout').then(resp => {
			cb && cb(resp.data)
			console.log(`user logout ${JSON.stringify(resp.data)}`)
		})
		.catch(err => {
			console.log(`user logout error ${err}`)
		})
	}
}