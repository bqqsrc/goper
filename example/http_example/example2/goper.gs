http: {
	server: {
		listen: 10000
		action: 1
		location: {
			url: /route1/route2
			methods: post
		}
		location: {
			url: /route3/route4
			methods: get post
		}
	}
}