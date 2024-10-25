package sqs

type Option func(*Driver)

func WithClient(c sqsClient) Option {
	return func(d *Driver) {
		d.sqsClient = c
	}
}

func AutoTestConnection() Option {
	return func(d *Driver) {
		d.testConnectionOnStartup = true
	}
}

func WithUrl(name string) Option {
	return func(d *Driver) {
		d.url = name
	}
}

func WithRegion(region string) Option {
	return func(d *Driver) {
		d.region = region
	}
}

func WithSharedCredentials(filename, profile string) Option {
	return func(d *Driver) {
		d.SetSharedCredentials(filename, profile)
	}
}
