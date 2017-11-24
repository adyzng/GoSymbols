<script>
export default {
	name: 'authredirect',
	mounted() {
		const hash = window.location.search.slice(1)
		console.log(`authredirect hash: ${hash}`)
        window.opener.location.href = window.location.origin + '/api/auth/authorize#' + hash
        window.close()
	}
};
</script>
