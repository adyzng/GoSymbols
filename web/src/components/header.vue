<template>
	<div class="nav-container">
		<div class="nav-logo">
			<a href="#"><img src="../assets/logo.png" /></a>
		</div>
		<ul class="nav-menu">
			<router-link class="nav-link" tag="li" :to="{name: 'branchs'}">BRANCHS</router-link>
			<router-link class="nav-link" tag="li" :to="{name: 'symbols'}">SYMBOLS</router-link>
			<li class="nav-link">
				<el-dropdown :hide-on-click="false" @visible-change="showMessage" trigger="click"> 
					<span class="el-dropdown-link">REPORT</span>
					<el-dropdown-menu slot="dropdown" 
						v-loading="loading"
						:show-timeout="100"
						class="dropdown-box" 
						style="min-width: 250px;">
						<el-dropdown-item v-for="(item, index) in messages" v-bind:key="index" class="msg-info">
							<my-message :item="item"></my-message>
						</el-dropdown-item>
					</el-dropdown-menu>
				</el-dropdown>
			</li>
			<li class="nav-link">
				<a v-if="!userLogin" href="/api/auth/login">LOGIN</a> 
				<!-- <a v-if="!userLogin" @click.stop="loginFn">LOGIN</a> -->
				<el-dropdown v-else :hide-on-click="false" @visible-change="showProfile" trigger="click">
					<span class="el-dropdown-link">{{userInfo.shortName | fltUpperCase}}</span>
					<el-dropdown-menu slot="dropdown" class="dropdown-box" v-loading="loading">
						<el-dropdown-item class="user-info">
							<img class="user-avator" src="../assets/user.jpg"/>
							<div class="user-desc">
								<span class="title">{{userInfo.fullName}}</span>
								<span>{{userInfo.mail}}</span>
							</div>
						</el-dropdown-item>
						<el-dropdown-item v-if="userInfo.jobTitle">{{userInfo.jobTitle}}</el-dropdown-item>
						<el-dropdown-item v-if="userInfo.cellPhone">{{userInfo.cellPhone}}</el-dropdown-item>
						<el-dropdown-item v-if="userInfo.bussPhone">{{userInfo.bussPhone}}</el-dropdown-item>
						<el-dropdown-item divided><a class="nav-link" href="/api/auth/logout">Log Out</a></el-dropdown-item>
					</el-dropdown-menu>
				</el-dropdown>
			</li>
		</ul>
	</div>
</template>

<script>
import pdb from '../api/pdb'
import Message from './message.vue'
import openWindow from '../utils/openWin'
import {USER_PROFILE} from "../utils/types"

export default {
	name: "header",
	components: {
		"myMessage": Message,
	},
	data() {
		return {
			loading: false,
			messages : [],
			userInfo : {},
			userLogin : false,
			userInfoTest: {
				givenName: "HUAN",
				displayName: "Zhang, Huan",
				avatar: "../assets/user.png",
				mail: "Huan.Zhang@Arcserve.com",
				jobTitle: "Sr Software Engineer",
				cellPhone: "+8618810514785",
				bussPhone: "+86 10 5089 0543 ext. 66-0543",
			},
		}
	},
	methods : {
		showMessage(val) {
			if (val) {
				let vm = this
				vm.loading = true
				pdb.getTodayMessages(messages => {
					vm.messages = messages.map(function(msg) {
						if (msg.status === 1) {
							return {
								title: msg.branch,
								content: `${msg.date} updated to build ${msg.build}.`,
							}
						} else {
							return {
								type: "warning",
								title: msg.branch,
								content: `Failed to update today.`,
							}
						}
					})
					if (vm.messages.length === 0) {
						vm.messages.push({
							title: '',
							content: `No update message.`,
						})
					}
					vm.loading = false
				});
			}
		},
		showProfile(val) {
			if (val) {
				let vm = this;
				pdb.fetchProfile(resp => {
					if (resp.code === 0) {
						vm.userInfo = {
							'mail': resp.data.mail,
							'jobTitle': resp.data.jobTitle,
							'fullName' : resp.data.displayName,
							'shortName': resp.data.givenName.toUpperCase(),
							'cellPhone': resp.data.mobilePhone,
							'bussPhone': resp.data.businessPhones.length? resp.data.businessPhones[0] : "",
						}
						vm.userLogin = true
					} else {
						vm.userLogin = false
					}
				})
			} 
		},
		loginFn() {
			let params = [
				`https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=15b003d2-fe18-4587-a8c7-45d0f064376a`,
				`redirect_uri=http://localhost:8010/#/login/authorize`,
				`response_mode=query`, //form_post
				`response_type=code`,
				`scope=User.Read Mail.Send`,
			]

			let url = encodeURI(params.join('&'))
			console.log(`redirect : ${url}`)

			let win = openWindow(url, 'Login', 400, 500)
			console.log(win)
		},
		logoutFn() {
			let vm = this
			pdb.userLogout(()=> {
				// push not working if already on the given route
				vm.$router.push("/")
			})
		}
	},
	mounted() {
		this.showProfile(true)
	}
}
</script>


<style scoped lang="scss">

$link-color: #4db3ff;
$text-color: #8391a5;
$title-color: #1f2d3d;

.nav-container {
	display: flex;
	height: 50px;
	align-items: center;
	justify-content: space-between;

	.nav-logo {
		height: 40px;
		display: inline-block;
		img {
			width: 40px;
			height: 40px;
		}
	}
	.nav-menu {
		margin: 0;
		line-height: 30px;
		font-size: 14px;
		list-style-type: none;
		li {
			margin-left: 20px;
			display: inline-block;
		}
	}
	.nav-link {
		cursor: pointer;
		color: $link-color;
		border-bottom: 3px solid transparent;
		a {
			text-decoration: none;
			color: $link-color;
		}
		.el-dropdown-link {
			color: $link-color;
		}
		&:hover, .current {
			border-bottom: 3px solid #42b983;
		}
	}
	.current {
		border-bottom: 3px solid #42b983;
	}
}

.dropdown-box {
	font-size: 12px;
	color: $text-color;
	span {
		line-height: 28px;
		display: block;
	}
}

.user-info {
	padding: 10px 10px;
	display: flex;
	align-items: center;
	justify-content: flex-start;

	.user-avator {
		border: 1px solid rgb(209, 219, 229);
		border-radius: 24px;
		width: 48px;
		height: 48px;
		margin-right: 10px;
	}
	.user-desc {
		align-self: center;
		.title {
			font-size: 14px;
			color: $title-color;
		}
	}
}
</style>

