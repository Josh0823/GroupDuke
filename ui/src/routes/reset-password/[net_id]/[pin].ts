/** @type {import('@sveltejs/kit').RequestHandler} */
export async function get({ params }) {
	if (!params.net_id || !params.pin) {
		return {
			status: 400
		};
	}

	return {
		body: {
			username: params.net_id,
			pin: params.pin
		}
	};
}
