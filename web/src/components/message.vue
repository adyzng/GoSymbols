<template>
	<div class="msg-box">
		<div class="msg-icon">
			<i :class="iconClass"/>
		</div>
		<div class="msg-text">
			<p class="msg-title">{{item.title}}</p>
			<span class="msg-content">{{item.content}}</span>
		</div>
	</div>
</template>

<script>
export default {
	name: 'Message',
	props : {
		item : {
			type: Object,
			required: true,
		},
		width : {
			type: Number,
			default: 300,
		}
	},
	watch : {
		"item" : function(val) {
			if (val.type === 'success') {
				this.iconClass = `${val.icon? val.icon:'el-icon-circle-check'} icon-success`;
			} else if(val.type === 'warning') {
				this.iconClass = `${val.icon? val.icon:'el-icon-warning'} icon-warn`;
			} else if(val.type === 'error') {
				this.iconClass = `${val.icon? val.icon:'el-icon-circle-cross'} icon-error`;
			} else {
				this.iconClass = `${val.icon? val.icon:'el-icon-information'} icon-info`;
			}
		}
	},
	data () {
		return {
			iconClass : "el-icon-information icon-info",
		};
	},
}
</script>


<style lang="scss" scoped>
$text-color: #8391a5;
$title-color: #1f2d3d;
$border-gray: 1px solid rgb(209, 219, 229);

$succ-color: #13ce66; // el-icon-circle-check
$info-color: #50bfff; // el-icon-information
$warn-color: #f7ba2a; // el-icon-warning
$erro-color: #ff4949; // el-icon-circle-cross

.icon-success {
	color: $succ-color,
}
.icon-info {
	color: $info-color,
}
.icon-warn {
	color: $warn-color,
}
.icon-error {
	color: $erro-color,
}


.msg-box {
	display: flex;
	align-items: flex-start;
	padding: 10px 0;
	max-width: 260px;

	.msg-icon {
		width: 36px;
		height: 36px;
		font-size: 36px;
		border-radius: 18px;
	}
	.msg-text {
		margin-left: 10px;
		line-height: 14px; 
		.msg-title {
			margin: 0 0 10px;
			font-size: 14px;
			color: $title-color;
		}
		.msg-content {
			font-size: 12px;
			color: $text-color;
			text-overflow: initial;
		}
	}
}
</style>

