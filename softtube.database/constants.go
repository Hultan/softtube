package database

const driverName = "sqlite3"
const dateLayout = "2006-01-02T15:04:05-0700"

func getConnectionString(databasePath string) string {
	return "file:" + databasePath + "?parseTime=true&_timeout=5000"
}
