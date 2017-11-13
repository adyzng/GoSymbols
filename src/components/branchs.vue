<template>
	<el-row>
		<el-col>	
			<h3> BRANCHES </h3>
			<el-card class="box-card">
				<div class="item">
					<el-row :gutter="20">
						<el-col :span="15" >
							<el-table 
								ref="branchTable"
								:data="branchesList"
								:max-height="maxHeight"
								highlight-current-row
								border
								@row-click="chooseBranch" >
								<el-table-column type="index" width="60"/>
								<el-table-column label="BRANCH" prop="storeName">
									<template slot-scope="props">
										<el-button type="text">{{ props.row.storeName }}</el-button>
									</template>
								</el-table-column>
								<el-table-column label="LATEST BUILD" prop="latestBuild">
								</el-table-column>
								<el-table-column label="UPDATE DATE" prop="updateDate" width="180">
								</el-table-column>
								<el-table-column label="OPERATION" width="140">
									<template slot-scope="props">
										<el-button
											@click.native.prevent="evtEditBranch(props.row)"
											type="text"
											size="small">
											EDIT
										</el-button>
										<el-button 
											@click.native.prevent="evtDeleteBranch(props.row)"
											type="text" 
											size="small">
											DELETE
										</el-button>
									</template>
								</el-table-column>
							</el-table>

							<el-dialog size="tiny" 
								ref="frmEditBranch"
								v-loading.body="loading"
								:title="editBranchTitle"
								:visible.sync="showEditDlg" >
								<el-form label-position="top" :label-width="frmLabelWidth" :model="editBranch">
									<el-form-item label="StoreName">
										<el-input v-model="editBranch.storeName" :disabled="true"></el-input>
									</el-form-item>
									<el-form-item label="StorePath">
										<el-input v-model="editBranch.storePath" auto-complete="off"></el-input>
									</el-form-item>
									<el-form-item label="BuildName">
										<el-input v-model="editBranch.buildName" auto-complete="off"></el-input>
									</el-form-item>
									<el-form-item label="BuildPath" style="margin-bottom: 0;">
										<el-input v-model="editBranch.buildPath" auto-complete="off"></el-input>
									</el-form-item>
								</el-form>
								<div slot="footer" class="dialog-footer">
									<el-button @click="showEditDlg = false">CANCEL</el-button>
									<el-button type="primary" @click="evtDlgEditBranch">CONFIRM</el-button>
								</div>
							</el-dialog>

						</el-col>
						<el-col :span="9">
							<pie-chart :chart-label="chartLabels" :chart-data="chartDataset" :height="maxHeight"></pie-chart>
						</el-col>
					</el-row>
				</div>
			</el-card>
		</el-col>

		<el-col class="bottom-gap">
			<my-build></my-build>
		</el-col>
	</el-row>
</template>

<script>
import pdb from "../api/pdb";
import {mapGetters} from "vuex"
import {PieChart} from "../utils/chart"
import * as types from '../utils/types'
import myBuild from "./builds.vue"

export default {
	components: {
		PieChart,
		myBuild,
	},
	data() {
		return {
			frmLabelWidth: "100px",
			loading: false,
			maxHeight: 300,
			editBranch: {},
			showCurrRow: {},
			showEditDlg: false,
			chartLabels:[],
			chartDataset: [{
				data: [],
				backgroundColor: [],
			}],
			colorList : ['#E46651', '#00B7C3', '#FF8C00', '#10893E', '#BF0077', '#20A0FF', '#EF6950', '#567C73'],
		};
	},
	computed : {
		...mapGetters([
			'curBuild',
			'curBranch',
			'branchesList',
		]),
		editBranchTitle : function() {
			return `EDIT BRANCH : ${this.editBranch.storeName}`;
		},
	},
	methods: {
		chooseBranch(row) {
			this.showCurrRow = row
			this.$store.commit(types.CHANGE_BRANCH, row.storeName)
			console.log(`change branch to ${row.storeName}.`)
		},
		evtEditBranch(row) {
			this.editBranch = {
				storeName: row.storeName,
				buildName: row.buildName,
				storePath: row.storePath,
				buildPath: row.buildPath,
			};
			this.showEditDlg = !this.showEditDlg;
			console.log(`start edit branch ${row.storeName}`);
		},
		evtDlgEditBranch() {
			let vm = this;
			let ok = false;
			vm.loading = true;
			pdb.modifyBranch(vm.editBranch, (data) => {
				vm.loading = false;
				let msgBox = {
					type : "error",
					title : "",
					message: "",
					confirmButtonText: "OK",
				}
				if (data instanceof Error) {
					msgBox.title = "ERROR";
					msgBox.message = `Modify branch ${vm.editBranch.storeName} error ${data.response}`;
				} else {
					if (data.code === 0) {
						ok = true;
						msgBox.type = 'info';
						msgBox.title = "INFO";
						msgBox.message = `Modify branch ${vm.editBranch.storeName} succeed.`;
					} else {
						msgBox.type = 'warning';
						msgBox.title = "WARNING";
						msgBox.message = `Modify branch ${vm.editBranch.storeName} failed: ${data.message}.`;
					}
					console.log(`modify branch, response: ${JSON.stringify(data)}.`)
				}
				this.$msgbox(msgBox).then(action => {
					vm.branchesList.forEach((val) => {
						if (val.storeName === vm.editBranch.storeName) {
							val.buildName = vm.editBranch.buildName;
							val.storePath = vm.editBranch.storePath;
							val.buildPath = vm.editBranch.buildPath;
						}
					});
					vm.showEditDlg = !ok;
					console.log(`modify branch action: ${action}`);
				});
			});
		},
		evtDeleteBranch(row) {
			pdb.deleteBranch(row, data => {
				if (data instanceof Error) {
					this.$msgbox({
						type: "warning",
						title: "WARNING",
						message: `NOT authorized to delete branch ${row.storeName}.`,
						confirmButtonText: "OK",
					});
				} else {
					console.log(`delete branch ${row.storeName}, return ${JSON.stringify(data)}.`);
				}
			});
		},
		validateBranch(id) {
			let vm = this;
			console.log(`validator: ${id}, branch: ${this.editBranch}.`)
			pdb.validateBranch(this.editBranch, (data)=>{
				console.log(`validate returns: ${JSON.stringify(data)}`);
				if (data instanceof Error) {
					console.log("validate branch error:", data);
					return;
				}
				vm.errValidates[id] = data.code == 0?
					null
					:
					new Error(data.message)
				console.log(`validate errors: ${vm.errValidates}`);
			});
		},
	},
	mounted() {
		let vm = this;
		pdb.fetchBranchs(data=> {
			if (Array.isArray(data)) {
				vm.$store.commit(types.BRANCH_LIST, [...data])
				vm.chooseBranch(data.length? data[0]:[])

				// update chart data
				const colorLen = vm.colorList.length;
				const chartLabels = new Array(data.length);
				const chartDataset = {
					data: new Array(data.length),
					backgroundColor: new Array(data.length),
				};

				data.forEach((val, idx) => {
					chartLabels[idx] = val.storeName;
					chartDataset.data[idx] = val.buildsCount? val.buildsCount : 0;
					chartDataset.backgroundColor[idx] = vm.colorList[idx % colorLen];
				});
				vm.chartLabels = chartLabels;
				vm.chartDataset = [chartDataset];

				//console.log('chartLabels:', vm.chartLabels);
				//console.log('chartDataset:', vm.chartDataset);
			}
			console.log(`fetch ${data.length} branchs.`)
		})
	},
	updated() {
		//console.log(`current branch ${this.showCurrRow.storeName}`)
		this.$refs.branchTable.setCurrentRow(this.showCurrRow)
	}

};
</script>

