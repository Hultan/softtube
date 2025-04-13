package database

const constDriverName = "mysql"
const constDateLayout = "2006-01-02 15:04:05"

//
// Videos table
//

const sqlVideosExists = "SELECT EXISTS(SELECT 1 FROM Videos WHERE id=?);"
const sqlVideosGetStatus = "SELECT status FROM Videos WHERE id=?"
const sqlVideosGet = "SELECT id, subscription_id, title, duration, published, added, status, save FROM Videos WHERE id=?"
const sqlVideosDelete = "DELETE FROM Videos WHERE id=? AND save=0"
const sqlVideosUpdateStatus = "UPDATE Videos SET status=? WHERE id=?"
const sqlVideosUpdateSave = "UPDATE Videos SET save=? WHERE id=?"
const sqlVideosUpdateDuration = "UPDATE Videos SET duration=?, seconds=? WHERE id=?"

const sqlVideosSearch = `SELECT Videos.id, Videos.subscription_id, Videos.title, Videos.duration, Videos.published, Videos.added, 
									Videos.status, Subscriptions.name , Videos.save
									FROM Videos 
									INNER JOIN Subscriptions ON Subscriptions.id = Videos.subscription_id 
									WHERE Videos.title LIKE ? OR Subscriptions.name LIKE ? 
									ORDER BY Videos.Added DESC`

const sqlVideosGetStats = `SELECT Videos.id
									FROM Videos 
									WHERE Videos.status NOT IN (0,4) OR Videos.save=1`

const sqlVideosGetLatest = `SELECT * FROM 
									(SELECT Videos.id, Videos.subscription_id, Videos.title, Videos.duration, Videos.published, Videos.added, Videos.status, Subscriptions.name, Videos.save 
									FROM Videos 
									INNER JOIN Subscriptions ON Videos.subscription_id = Subscriptions.id 
									ORDER BY $ORDER$
									LIMIT 200) as Newest

									UNION

									SELECT * FROM
										(SELECT Videos.id, Videos.subscription_id, Videos.title, Videos.duration, Videos.published, Videos.added, Videos.status, Subscriptions.name, Videos.save 
										FROM Videos 
										INNER JOIN Subscriptions ON Videos.subscription_id = Subscriptions.id 
										WHERE Videos.status NOT IN (0,4) OR Videos.save=1) as Downloaded

									ORDER BY $ORDER$`

const sqlVideosGetFailed = `SELECT Videos.id, Videos.subscription_id, Videos.title, Videos.duration, 
										Videos.published, Videos.added, Videos.status, Subscriptions.name, Videos.save 
									FROM Videos 
									INNER JOIN Subscriptions ON Videos.subscription_id = Subscriptions.id
									WHERE Videos.status = 1
									ORDER BY added desc`

const sqlVideosInsert = `INSERT IGNORE INTO Videos (id, subscription_id, title, duration, published, added, status, 
save, seconds) 
								VALUES (?, ?, ?, ?, ?, ?, 0, 0, ?);`

//
// Downloads table
//

// TODO : Make max downloads a setting
const sqlDownloadsGetAll = "SELECT video_id FROM Download LIMIT 5"
const sqlDownloadsInsert = "INSERT INTO Download (video_id) VALUES (?)"
const sqlDownloadsDelete = "DELETE FROM Download WHERE video_id=?"

//
// Log table
//

const sqlLogInsert = `INSERT INTO Log (type, message, time) VALUES (?, ?, NOW());`
const sqlLogGetLatest = `SELECT id, type, message FROM Log                 
ORDER BY id desc
LIMIT 50`

//
// Subscriptions table
//

const sqlSubscriptionsGetAll = "SELECT id, name, frequency, last_checked, next_update FROM Subscriptions"
const sqlSubscriptionsGet = "SELECT id, name, frequency, last_checked, next_update FROM Subscriptions WHERE id=?"
const sqlSubscriptionsUpdateLastChecked = "UPDATE Subscriptions SET last_checked=?, next_update=? WHERE id=?"
