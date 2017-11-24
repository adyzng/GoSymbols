<template>
	<el-row>
		<h3>{{ curBrBuild }}</h3>	
		<el-card class="box-card">
			<div slot="header" class="clearfix">
				<el-row :gutter="20">
					<el-col :span="12">
						<el-input placeholder="Query Symbol ..."
							icon="close"
							v-model.lazy="queryWord"
							:on-icon-click="handleIconClean" >
							<el-select v-model="queryType" slot="prepend" class="queryBox">
								<el-option label="Name" value="name"></el-option>
								<el-option label="Hash" value="hash"></el-option>
							</el-select>
							<!-- <el-button slot="append" icon="search"></el-button>-->
						</el-input>
					</el-col>
					<el-col :span="12">
						<el-pagination
							layout="total, sizes, prev, pager, next"
							@current-change="handlePageChange"
							@size-change="handleSizeChange"
							:current-page="pageCur"
							:page-sizes="pageSizeLst"
							:page-size="pageSize"
							:total="pageTotal">
						</el-pagination>
					</el-col>
				</el-row>
			</div>
			<div class="text item">
				<el-table :data="tableData" 
					v-loading.body="loading"
					border 
					stripe >
					<el-table-column type="index" width="60"/>
					</el-table-column>
					<el-table-column label="NAME" prop="name" width="200" show-overflow-tooltip>
						<template slot-scope="props">
							<i class="el-icon-document"></i>
							<a :href="symURL(props.row.name, props.row.hash)">
								{{ props.row.name }}
							</a>
						</template>
					</el-table-column>
					<el-table-column label="ARCH" prop="arch" width="100"
						:filters="[{ text: 'x86', value: 'x86' }, { text: 'x64', value: 'x64' }]"
						:filter-method="filterArch"
						filter-placement="bottom-end">
						<template slot-scope="scope">
							<el-tag :type="scope.row.arch === 'x86' ? 'primary' : 'success'"
								close-transition > 
								{{scope.row.arch}}
							</el-tag>
						</template>
					</el-table-column>
					<el-table-column label="VERSION" prop="version" width="120">
					</el-table-column>
					<el-table-column label="HASH" prop="hash" width="350">
					</el-table-column>
					<el-table-column label="PATH" prop="path" show-overflow-tooltip>
					</el-table-column>
				</el-table>
			</div>
		</el-card>
	</el-row>
</template>

<script>
import pdb from "../api/pdb"
import {mapGetters} from "vuex"
import {CHANGE_BUILD} from "../utils/types"

export default {
	data() {
		return {
			loading: false,
			downloadURL: '/api/symbol',
			queryType: 'name',
			queryWord: '',
			refreshData: false,
			updateData: 0,
			tableData:[],
			cacheData:[],
			pageTotal: 0,
			pageSizeLst: [20, 50, 100, 150, 200, 300],
            pageSize: 20,
			pageCur: 1,
		};
	},
	computed : {
		...mapGetters([
			'curBuild',
			'curBranch',
		]),
		curBrBuild : function() {
			if (this.curBranch && this.curBuild) {
				this.refreshData = true
			}
			return this.curBranch + ' / ' + this.curBuild
		}
	},
	watch : {
		updateData: function(newVal, oldVal) {
			let idx = this.pageCur == 0? 0 : this.pageCur - 1;
			let start = idx * this.pageSize;
			this.tableData = this.cacheData.slice(start, start + this.pageSize);
		},
		queryWord: function(newVal, oldVal) {
			if (!this.queryType) {
				return
			}
			if (!this.queryWord) {
				this.updateData++;
				return
			}
			let size = 0;
			let word = this.queryWord.toLowerCase();
			this.tableData = this.cacheData.filter(item=>{
				let ret = false;
				if (size >= this.pageSize) {
					return false;
				}
				if (this.queryType === 'name') {
					ret = item.name.toLowerCase().indexOf(word) !== -1;
				}
				if (this.queryType === 'hash') {
					ret = item.hash.toLowerCase().indexOf(word) !== -1;
				}
				if (ret) {
					size++;
				}
				return ret;
			})
		},
	},
	methods: {
		symURL(name, hash) {
			// /api/v1/symbol/:branch/:hash?name=name
			// console.log(`symbol url ${this.curBranch}`);
			return `${this.downloadURL}/${this.curBranch}/${hash}/${name}`;
		},
		filterArch(value, row) {
        	return row.arch === value;
		},
		handleIconClean() {
			this.queryWord = '';
		},
        handleSizeChange(val) {
			this.pageSize = val;
			this.updateData++;
        },
        handlePageChange(val) {
			this.pageCur = val;
			this.updateData++;
		},
	},
	activated() {
		let vm = this;
		if (vm.refreshData) {
			vm.loading = true
			pdb.fetchSymbols(vm.curBranch, vm.curBuild, data => {
				if (Array.isArray(data)) {
					vm.pageCur = 1
					vm.cacheData = data
					vm.pageTotal = data.length
				}
				vm.updateData++
				vm.loading = false
				vm.refreshData = false
			})
			console.log(`symbols activated for ${vm.curBrBuild}.`)
		}
	}
}
</script>


<style scoped>

</style>

