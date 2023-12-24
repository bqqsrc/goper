http: {
	demo: {
		action: 1
		name: default
	}
	server: {
		listen: 10000
		location: {
			url: /route1/route2
		}
		location: {
			url: /route3/route4
			method: {
				demo: {
					action: 2
					name: location1
				}
				methods: get 
			}
		}
	}
	server: {
		listen: 10001
		location: {
			url: /route1/route2
			method: {
				demo: {
					action: 3
					name: location2
				}
				methods: post 
			}
			method: {
				demo: {
					action: 4
					name: location3
				}
				methods: get 
			}
		}
	}
}