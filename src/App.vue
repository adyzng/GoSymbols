<template>
	<div id="app">
		<!-- top navbar -->
		<div class="header">
			<div class="header-container">
				<myHeader></myHeader>
			</div>
		</div>

		<!-- content -->
		<div class="page-container">
			<el-row :gutter="20">
				<el-col>
					<el-breadcrumb separator="/">
						<el-breadcrumb-item :to="{ path: '/' }">HOME</el-breadcrumb-item>
						<el-breadcrumb-item v-for="(route, index) in $route.matched" v-bind:key="index">
							{{route.name | fltUpperCase}}
						</el-breadcrumb-item>
					</el-breadcrumb>
				</el-col>

				<el-col class="bottom-gap">
					<transition name="fade">
						<keep-alive>
							<router-view></router-view>
						</keep-alive>
					</transition>
				</el-col>
			</el-row>
		</div>
	</div>
</template>

<script>
import myHeader from "./components/header.vue"
import myBranch from "./components/branchs.vue"
import mySymbol from "./components/symbols.vue"

export default {
	components: {
		myHeader,
		myBranch,
		mySymbol,
	},
}
</script>

<style lang="scss">

a {
	text-decoration: none;
	color: #4db3ff;
}

body {
	margin: 0;
	font-family: "Helvetica Neue",Helvetica,"PingFang SC","Hiragino Sans GB","Microsoft YaHei","微软雅黑",Arial,sans-serif;

	.header {
		position: fixed;
		top: 0;
		width: 100vw;
		height: 50px;
		padding: 5px 0px;
		background-color: #eef1f6;
		box-shadow: 0 0 1px rgba(0,0,0,0.25);
		transition: background-color 0.3s ease-in-out; 
		z-index: 2;

		.header-container {
			width: 1140px;
			margin: 0 auto;
			padding-right: 15px;
		}
	}

	.page-container {
		padding: 20px;
		margin: 60px auto;
		width: 1140px;
		
		.el-breadcrumb {
			//font-size: 14px;
			margin: 10px 0;
		}
		.with-boader {
			border: 1px solid rgb(223, 230, 236);
			margin: 0 10px;
		}
		.queryBox {
			width: 120px;
		}
		h2, h3, h4, h5 {
			font-weight: 400;
			color: #1f2f3d;
		}
	}
}

.fade-enter-active, .fade-leave-active {
  transition: opacity .5s
}
.fade-enter, .fade-leave-to {
  opacity: 0
}

/*
.test {
	color: #E74856', '#10893E', '#FF8C00', '#0078D7', '#BF0077', '#00B7C3', '#EF6950', '#567C73';
}
*/

</style>