<template>
	<el-row style="margin-top: 40px;">
        <h3> {{ curBranch }}</h3>
        <el-card class="box-card">
			<div slot="header" class="clearfix">
				<el-row :gutter="20">
					<el-col :span="12">
						<el-input placeholder="Query Build ..."
							icon="close"
							v-model.lazy="queryWord"
							:on-icon-click="handleIconClean" >
							<el-select v-model="queryType" slot="prepend" class="queryBox">
								<el-option label="BUILD ID" value="id"></el-option>
                                <el-option label="VERSION" value="version"></el-option>
							</el-select>
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
                <el-table 
                    ref="buildTable"
                    :data="tableData" 
                    v-loading.body="loading"
                    highlight-current-row
                    @row-click="chooseBuild"
                    border>
                    <el-table-column label="BUILD ID" prop="id">
                        <template slot-scope="props">
                            <el-button type="text" @click.stop="showSymbols(props.row)">
                                {{ props.row.id }}
                            </el-button>
                        </template>
                    </el-table-column>
                    <el-table-column label="VERSION" prop="version">
                    </el-table-column>
                    <el-table-column label="BRANCH" prop="branch">
                    </el-table-column>
                    <el-table-column label="DATETIME" prop="date">
                    </el-table-column>
                    <el-table-column label="COMMENT" prop="comment" >
                    </el-table-column>
                </el-table>
		    </div>
        </el-card>
	</el-row>
</template>


<script>
import pdb from "../api/pdb"
import {mapGetters} from "vuex"
import * as types from '../utils/types'

export default {
    data() {
        return {
            loading: false,
            queryType: 'id',
            queryWord: '',
            showCurRow: {},
            updateData: 0,
            tableData: [],
            pageTotal: 0,
            pageSizeLst: [10, 20, 50, 100],
            pageSize: 10,
            pageCur: 1,
        }
    },
    computed: {
        ...mapGetters([
            'curBranch',
            'buildsList',
        ]),
    },
    watch : {
        showCurRow : function (val, old) {
            //console.log(`build current row change ${old.version} => ${val.version}`)
            this.$store.commit(types.CHANGE_BUILD, val.id || "")
        },
        updateData : function (val){
            let idx = this.pageCur == 0? 0 : this.pageCur - 1;
            let start = idx * this.pageSize;
            this.tableData = this.buildsList.slice(start, start + this.pageSize);
        },
        curBranch : function(val) {
            let vm = this;
            //this.loading = true;
            pdb.fetchBuilds(val, data => {
                if (Array.isArray(data)) {
                    vm.pageCur = 1
                    vm.pageTotal = data.length
                    vm.showCurRow = data.length > 0? data[0] : []
                    vm.$store.commit(types.BUILD_LIST, data)
                }
                vm.updateData++;
                //vm.loading = false;
                console.log(`fetch ${data.length} builds for branch ${val}.`);
            });
        },
        queryWord: function(newVal, oldVal) {
			if (!this.queryType) {
				return
			}
			if (!this.queryWord) {
				this.updateData++;
				return
            }
			this.tableData = this.buildsList.filter(item=>{
				if (this.queryType === 'id') {
					return Number(item.id) === Number(newVal);
                }
                if (this.queryType === 'version') {
                    return item.version.indexOf(newVal) !== -1;
                }
				return false;
			})
		},
    },
    methods: {
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
        chooseBuild(row) {
            this.showCurRow = row
        },
        showSymbols(row) {
            console.log(`symbols for ${this.curBranch}, ${row.version}`);
            this.showCurRow = row
            this.$router.push({name: 'symbols'})
        },
    },
    updated() {
        //console.log(`build updated set current row ${JSON.stringify(this.showCurRow)}`)
        this.$refs.buildTable.setCurrentRow(this.showCurRow)
    },
    activated() {
        //console.log(`build activated set current row ${JSON.stringify(this.showCurRow)}`)
        this.$refs.buildTable.setCurrentRow(this.showCurRow)
    }
};
</script>