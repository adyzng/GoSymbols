import {Pie, mixins} from 'vue-chartjs'
const {reactiveProp} = mixins;

function merge(...objects) {
	let result = {};
	function assign(val, key) {
		if (typeof result[key] === 'object' && typeof val === 'object') {
			result[key] = merge(result[key], val);
		} else {
			result[key] = val;
		}
	}
	objects.forEach(val => {
		for (let key in val) {
			assign(val[key], key)
		}
	})
	return result;
}

const PieChart = {
	name: 'pie-chart',
	extends: Pie,
	//mixins : [reactiveProp],
	props: {
		chartLabel: {
			type: Array,
			required: true,
		},
		chartData: {
			type: Array,
			required: true,
		},
		chartOption: {
			type: Object,
			default: ()=> {},
		},
	},
	watch: {
		'chartData': function(newVal, oldVal) {
			const options = merge(this.options, this.chartOption);
			//console.log(`renderChart with label: ${this.chartLabel}`);
			this.renderChart({
				labels: this.chartLabel,
				datasets: this.chartData,
			}, options );
		},
	},
	data() {
        return {
			options : {
				responsive: true, 
				maintainAspectRatio: false
			},
		}
	},
};

export default {
	PieChart,
}

export {
	PieChart,
};