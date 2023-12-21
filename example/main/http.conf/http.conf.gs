server: {
	listen: 10003
	protocol: {
		name: http 
	}
	location: {
		methods: GET
		url: `/demo`
		huxing: dynamic
	}
	location: {
		methods: GET
		url: `/demo/:ac`
		huxing: dynamic
	}
}