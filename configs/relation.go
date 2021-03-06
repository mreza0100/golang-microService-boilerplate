package configs

var RelationConfigs *serviceConfigs

func init() {
	const (
		port   = 9092
		dbPort = 5435
	)
	RelationConfigs = &serviceConfigs{
		Addr:    "localhost:" + str(port),
		Port:    port,
		DBPort:  dbPort,
		Timeout: stdConnectionTimeout,
	}
}
