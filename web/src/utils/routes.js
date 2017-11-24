import Branchs from "../components/branchs.vue"
import Symbols from "../components/symbols.vue"
import Authorize from "../view/authorize.vue"

export const routes = [
	{
		name: 'branchs',
		path: '/',
		component: Branchs,
	},
	{
		name: 'symbols',
		path: '/symbols',
		component: Symbols,
		//props: true, 
	},
	{
		name: 'authredirect',
		path: '/login/authorize',
		component: Authorize,
	},
	{
		name: 'edit',
		path: '/branch',
		component: Branchs,
		children: [
			{
				path: '/edit',
				component: Branchs,
				meta: {
					requireAuth: true
				}
			},
			{
				path: '/update',
				component: Branchs,
				meta: {
					requireAuth: true
				}
			}
		]
	}

];

export default routes;