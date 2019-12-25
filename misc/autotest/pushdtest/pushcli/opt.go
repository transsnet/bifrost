package pushcli

type config struct {
	clientid string
	// = "clientid"
	service string
	// = "service"
	cleansession bool
	// = true
	payload []byte
	// = []byte("hello world")
}

type opt func(c *config)

func SetClientID(id string) opt {
	return func(c *config) {
		c.clientid = id
	}
}

func SetService(s string) opt {
	return func(c *config) {
		c.service = s
	}
}

func SetCleanSession() opt {
	return func(c *config) {
		c.cleansession = true
	}
}

func SetPayLoad(p []byte) opt {
	return func(c *config) {
		c.payload = p
	}
}
