import Vue from 'vue'
import Vuex from 'vuex'
import * as types from './types'

Vue.use(Vuex)

export default new Vuex.Store({
	state : {
		branch : '', 		// current selected branch
		branchList : [],	// all branchs list
		build : '',			// current selected build
		buildList : [],		// current selected build list
		token: "",
		userProfile : {},
	},
	getters : {
		curBranch: store => {
			return store.branch;
		},
		branchesList : store => {
			return store.branchList
		},
		curBuild : store => {
			return store.build;
		},
		buildsList : store => {
			return store.buildList
		},
		token : store => {
			return store.token
		},
		userProfile : store => {
			return store.userProfile
		},
	},
	mutations : {
		[types.CHANGE_BRANCH](store, val) {
			store.branch = val;
		},
		[types.BRANCH_LIST](store, val) {
			store.branchList = val;
		},
		[types.CHANGE_BUILD](store, val) {
			store.build = val;
		},
		[types.BUILD_LIST](store, val) {
			store.buildList = val;
		},
		[types.USER_PROFILE](store, val) {
			store.userProfile = val;
		},
	},
	actions : {
		[types.CHANGE_BRANCH] ({ commit, state }, val) {
			commit(types.CHANGE_BRANCH, val)
		},
		[types.BRANCH_LIST] ({ commit, state }, val) {
			commit(types.BRANCH_LIST, val)
		},
		[types.CHANGE_BUILD] ({ commit, state }, val) {
			commit(types.CHANGE_BUILD, val)
		},
		[types.BUILD_LIST] ({ commit, state }, val) {
			commit(types.BUILD_LIST, val)
		},
		[types.USER_PROFILE] ({ commit, state }, val) {
			commit(types.USER_PROFILE, val)
		},
	},
})