package db

import (
	"ClipsArchiver/internal/config"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log/slog"
	"slices"
	"time"
)

var db *sql.DB
var logger *slog.Logger

type User struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	ApexUsername string `json:"apexUsername"`
	ApexUid      string `json:"apexUid"`
}

type Clip struct {
	Id                int            `json:"id"`
	OwnerId           int            `json:"ownerId"`
	Filename          string         `json:"filename"`
	IsProcessed       bool           `json:"isProcessed"`
	CreatedAt         sql.NullTime   `json:"createdOn"`
	Duration          int            `json:"duration"`
	Map               sql.NullInt32  `json:"map"`
	GameMode          sql.NullString `json:"gameMode"`
	Legend            sql.NullInt32  `json:"legend"`
	MatchHistoryFound bool           `json:"matchHistoryFound"`
	Tags              []string       `json:"tags"`
	ThumbnailUri      string         `json:"thumbnailUri"`
	VideoUri          string         `json:"videoUri"`
}

type QueueEntry struct {
	Id               int            `json:"id"`
	ClipId           int            `json:"clipId"`
	Status           string         `json:"status"`
	StartedAt        sql.NullTime   `json:"startedAt"`
	FinishedAt       sql.NullTime   `json:"finishedAt"`
	Operation        string         `json:"operation"`
	ErrorMessage     sql.NullString `json:"errorMessage"`
	CombineWithId    sql.NullInt32  `json:"combineWithId"`
	DesiredStartTime sql.NullInt32  `json:"desiredStartTime"`
	DesiredEndTime   sql.NullInt32  `json:"desiredEndTime"`
	Filename         string
}

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ClipTag struct {
	ClipId int
	TagId  int
}

type MatchHistory struct {
	Id        int
	UserId    int
	GameStart sql.NullTime
	GameEnd   sql.NullTime
	Map       sql.NullInt32
	Legend    sql.NullInt32
	GameMode  string
}

type Map struct {
	Id        int
	Name      string
	CardImage string
	AlsName   string
}

type Legend struct {
	Id        int
	Name      string
	CardImage string
}

func SetupDb(l *slog.Logger) error {
	logger = l
	dbConfig := config.GetDatabaseInfo()
	cfg := mysql.Config{
		User:      dbConfig.Username,
		Passwd:    dbConfig.Password,
		Net:       "tcp",
		Addr:      dbConfig.Address,
		DBName:    dbConfig.Name,
		ParseTime: true,
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())

	l.Debug(fmt.Sprintf("Opened database connection to %s", dbConfig.Address))
	return err
}

func GetAllUsers() ([]User, error) {
	logger.Debug("Fetching all users")
	var users []User

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		logger.Error(fmt.Sprintf("Error fetching all users: %s", err.Error()))
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		if err = rows.Scan(&user.Id, &user.Name, &user.ApexUsername, &user.ApexUid); err != nil {
			logger.Error(fmt.Sprintf("Error fetching all users: %s", err.Error()))
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		logger.Error(fmt.Sprintf("Error fetching all users: %s", err.Error()))
		return nil, err
	}
	return users, nil
}

func GetAllTags() ([]Tag, error) {
	logger.Debug("Fetching all tags")
	var tags []Tag

	rows, err := db.Query("SELECT * FROM tags")
	if err != nil {
		logger.Error(fmt.Sprintf("Error fetching all tags: %s", err.Error()))
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var tag Tag
		if err = rows.Scan(&tag.Id, &tag.Name); err != nil {
			logger.Error(fmt.Sprintf("Error fetching all tags: %s", err.Error()))
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		logger.Error(fmt.Sprintf("Error fetching all tags: %s", err.Error()))
		return nil, err
	}
	return tags, nil
}

func GetAllLegends() ([]Legend, error) {
	logger.Debug("Fetching all legends")
	var legends []Legend

	rows, err := db.Query("SELECT * FROM legends")
	if err != nil {
		logger.Error(fmt.Sprintf("Error fetching all legends: %s", err.Error()))
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var legend Legend
		if err = rows.Scan(&legend.Id, &legend.Name, &legend.CardImage); err != nil {
			logger.Error(fmt.Sprintf("Error fetching all legends: %s", err.Error()))
			return nil, err
		}
		legends = append(legends, legend)
	}

	if err = rows.Err(); err != nil {
		logger.Error(fmt.Sprintf("Error fetching all legends: %s", err.Error()))
		return nil, err
	}
	return legends, nil
}

func GetAllMaps() ([]Map, error) {
	logger.Debug("Fetching all maps")
	var maps []Map

	rows, err := db.Query("SELECT id,name,card_image FROM maps")
	if err != nil {
		logger.Error(fmt.Sprintf("Error fetching all maps: %s", err.Error()))
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var gameMap Map
		if err = rows.Scan(&gameMap.Id, &gameMap.Name, &gameMap.CardImage); err != nil {
			logger.Error(fmt.Sprintf("Error fetching all maps: %s", err.Error()))
			return nil, err
		}
		maps = append(maps, gameMap)
	}

	if err = rows.Err(); err != nil {
		logger.Error(fmt.Sprintf("Error fetching all maps: %s", err.Error()))
		return nil, err
	}
	return maps, nil
}

func GetAllQueueEntries() ([]QueueEntry, error) {
	logger.Debug("Fetching all queue entries")
	var queueEntries []QueueEntry

	rows, err := db.Query("SELECT * FROM clips_queue")
	if err != nil {
		logger.Error(fmt.Sprintf("Error fetching all queue entries: %s", err.Error()))
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var queueEntry QueueEntry
		if err = rows.Scan(&queueEntry.Id, &queueEntry.ClipId, &queueEntry.Status, &queueEntry.StartedAt, &queueEntry.FinishedAt, &queueEntry.Operation); err != nil {
			logger.Error(fmt.Sprintf("Error fetching all queue entries: %s", err.Error()))
			return nil, err
		}

		queueEntries = append(queueEntries, queueEntry)
	}
	if err = rows.Err(); err != nil {
		logger.Error(fmt.Sprintf("Error fetching all queue entries: %s", err.Error()))
		return nil, err
	}
	return queueEntries, nil
}

func GetQueueEntryByClipId(id int) (QueueEntry, error) {
	logger.Debug(fmt.Sprintf("Fetching queue entries for clip id: %d", id))
	var queueEntry QueueEntry
	row := db.QueryRow("SELECT * FROM clips_queue WHERE clips_queue.clip_id = ?", id)

	err := row.Scan(&queueEntry.Id, &queueEntry.ClipId, &queueEntry.Status, &queueEntry.StartedAt, &queueEntry.FinishedAt, &queueEntry.Operation)
	if err != nil {
		logger.Error(fmt.Sprintf("Error fetching queue entries for clip id: %d. %s", id, err.Error()))
	}
	return queueEntry, err
}

func GetClipsForDate(dateOf time.Time) ([]Clip, error) {
	logger.Debug(fmt.Sprintf("Fetching clips for date: %s", dateOf.String()))
	var clips []Clip

	dateAfter := dateOf.AddDate(0, 0, 1)

	rows, err := db.Query("SELECT * FROM clips WHERE clips.is_processed = 1 AND clips.created_at >= ? AND clips.created_at < ?", dateOf, dateAfter)
	if err != nil {
		logger.Error(fmt.Sprintf("Error fetching clips for date: %s. %s", dateOf.String(), err.Error()))
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var clip Clip
		if err = rows.Scan(&clip.Id, &clip.OwnerId, &clip.Filename, &clip.IsProcessed, &clip.CreatedAt, &clip.Duration, &clip.Map, &clip.GameMode, &clip.Legend, &clip.MatchHistoryFound); err != nil {
			logger.Error(fmt.Sprintf("Error fetching clips for date: %s. %s", dateOf.String(), err.Error()))
			return nil, err
		}

		tags, err := GetTagsForClip(clip.Id)
		if err == nil {
			clip.Tags = tags
		}
		clip.VideoUri = fmt.Sprintf("http://10.0.0.10:8080/clips/archive/%s", clip.Filename)
		clip.ThumbnailUri = fmt.Sprintf("http://10.0.0.10:8080/clips/archive/Thumbnails/%s", clip.Filename+".png")
		clips = append(clips, clip)
	}

	if err = rows.Err(); err != nil {
		logger.Error(fmt.Sprintf("Error fetching clips for date: %s. %s", dateOf.String(), err.Error()))
		return nil, err
	}
	return clips, nil
}

func AddClip(ownerId int, filename string, createdAt time.Time) (Clip, error) {
	logger.Debug(fmt.Sprintf("Adding clip with owner ID: %d, filename: %s, createdAt: %s", ownerId, filename, createdAt.String()))
	var clip Clip
	clipResult, err := db.Exec("INSERT INTO clips (owner_id, filename, is_processed, created_at) VALUES (?, ?, ?, ?)", ownerId, filename, 0, createdAt)
	if err != nil {
		logger.Error(fmt.Sprintf("Error adding clip: %s", err.Error()))
		return clip, err
	}

	id, err := clipResult.LastInsertId()
	if err != nil {
		logger.Error(fmt.Sprintf("Error adding clip: %s", err.Error()))
		return clip, err
	}

	clip, err = GetClipById(int(id))
	if err != nil {
		logger.Error(fmt.Sprintf("Error adding clip: %s", err.Error()))
		return clip, err
	}

	err = AddClipToQueue(int(id))
	if err != nil {
		logger.Error(fmt.Sprintf("Error adding clip: %s", err.Error()))
	}

	return clip, err
}

func AddClipToQueue(clipId int) error {
	_, err := db.Exec("INSERT INTO clips_queue (clip_id, status) VALUES (?, ?)", clipId, "pending")
	if err != nil {
		logger.Error(fmt.Sprintf("Error adding clip with id %d to queue: %s", clipId, err.Error()))
	}

	return err
}

func GetTagsForClip(clipId int) ([]string, error) {
	logger.Debug(fmt.Sprintf("Fetching tags for clip id: %d", clipId))
	var tags []Tag

	rows, err := db.Query("SELECT tag_id,name FROM clips_tags INNER JOIN tags ON clips_tags.tag_id = tags.id WHERE clips_tags.clip_id = ?", clipId)
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting tags for clip with id %d: %s", clipId, err.Error()))
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var tag Tag
		if err = rows.Scan(&tag.Id, &tag.Name); err != nil {
			logger.Error(fmt.Sprintf("Error getting tags for clip with id %d: %s", clipId, err.Error()))
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		logger.Error(fmt.Sprintf("Error getting tags for clip with id %d: %s", clipId, err.Error()))
		return nil, err
	}

	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}

	return tagNames, nil
}

func UpdateClipTags(old Clip, new Clip) error {
	logger.Debug(fmt.Sprintf("Updating tags for clip with id %d", old.Id))
	var tagsToRemove []string
	tagsToRemove = make([]string, 0)
	for _, existingTag := range old.Tags {
		if !slices.Contains(new.Tags, existingTag) {
			tagsToRemove = append(tagsToRemove, existingTag)
		}
	}

	for _, tag := range new.Tags {
		var existingTag Tag
		row := db.QueryRow("SELECT * FROM tags WHERE tags.name = ?", tag)

		err := row.Scan(&existingTag.Id, &existingTag.Name)
		if err != nil {
			_, err = db.Exec("INSERT INTO tags (name) VALUES (?)", tag)
			if err != nil {
				logger.Error(fmt.Sprintf("Error adding tag %s to clip %d: %s", tag, old.Id, err.Error()))
				return err
			}
			row = db.QueryRow("SELECT * FROM tags WHERE tags.name = ?", tag)
			err = row.Scan(&existingTag.Id, &existingTag.Name)
			if err != nil {
				logger.Error(fmt.Sprintf("Error adding tag %s to clip %d: %s", tag, old.Id, err.Error()))
				return err
			}
		}
		var existingClipTag ClipTag
		row = db.QueryRow("SELECT * FROM clips_tags WHERE clip_id = ? AND tag_id = ?", old.Id, existingTag.Id)
		err = row.Scan(&existingClipTag.ClipId, &existingClipTag.TagId)
		if err == nil {
			continue
		}
		_, err = db.Exec("INSERT INTO clips_tags (clip_id, tag_id) VALUES (?, ?)", old.Id, existingTag.Id)
		if err != nil {
			logger.Error(fmt.Sprintf("Error adding tag %s to clip %d: %s", tag, old.Id, err.Error()))
			return err
		}
	}

	for _, tag := range tagsToRemove {
		var existingTag Tag
		row := db.QueryRow("SELECT * FROM tags WHERE tags.name = ?", tag)

		err := row.Scan(&existingTag.Id, &existingTag.Name)
		if err != nil {
			logger.Error(fmt.Sprintf("Error removing tag %s from clip %d: %s", tag, old.Id, err.Error()))
			return err
		}

		_, err = db.Exec("DELETE FROM clips_tags WHERE clip_id = ? AND tag_id = ?", new.Id, existingTag.Id)
		if err != nil {
			logger.Error(fmt.Sprintf("Error removing tag %s from clip %d: %s", tag, old.Id, err.Error()))
			return err
		}
	}
	return nil
}

func UpdateClip(clip Clip) error {
	logger.Debug(fmt.Sprintf("Updating clip %d", clip.Id))
	_, err := db.Exec("UPDATE clips SET clips.map = ?, clips.game_mode = ?, clips.legend = ?, clips.match_history_found = ? WHERE clips.id = ?", clip.Map, clip.GameMode, clip.Legend, clip.MatchHistoryFound, clip.Id)
	logger.Error(fmt.Sprintf("Error updating clip %d: %s", clip.Id, err.Error()))
	return err
}

func GetClipById(clipId int) (Clip, error) {
	logger.Debug(fmt.Sprintf("Getting clip with id %d", clipId))
	var clip Clip
	row := db.QueryRow("SELECT * FROM clips WHERE clips.id = ?", clipId)

	err := row.Scan(&clip.Id, &clip.OwnerId, &clip.Filename, &clip.IsProcessed, &clip.CreatedAt, &clip.Duration, &clip.Map, &clip.GameMode, &clip.Legend, &clip.MatchHistoryFound)
	tags, err := GetTagsForClip(clip.Id)
	if err == nil {
		clip.Tags = tags
	}
	clip.VideoUri = fmt.Sprintf("http://10.0.0.10:8080/clips/archive/%s", clip.Filename)
	clip.ThumbnailUri = fmt.Sprintf("http://10.0.0.10:8080/clips/archive/Thumbnails/%s", clip.Filename+".png")
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get clip with id %d: %s", clipId, err.Error()))
	}
	return clip, err
}

func GetClipByFilename(filename string) (Clip, error) {
	var clip Clip
	row := db.QueryRow("SELECT * FROM clips WHERE clips.filename = ?", filename)

	err := row.Scan(&clip.Id, &clip.OwnerId, &clip.Filename, &clip.IsProcessed, &clip.CreatedAt, &clip.Duration, &clip.Map, &clip.GameMode, &clip.Legend, &clip.MatchHistoryFound)

	if err != nil {
		return clip, err
	}

	tags, err := GetTagsForClip(clip.Id)
	if err == nil {
		clip.Tags = tags
	}
	clip.VideoUri = fmt.Sprintf("http://10.0.0.10:8080/clips/archive/%s", clip.Filename)
	clip.ThumbnailUri = fmt.Sprintf("http://10.0.0.10:8080/clips/archive/Thumbnails/%s", clip.Filename+".png")
	return clip, err
}

func DeleteClipById(clipId int) error {
	_, err := db.Exec("DELETE FROM clips_queue WHERE clip_id = ?", clipId)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM clips_tags WHERE clip_id = ?", clipId)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM clips WHERE id = ?", clipId)
	return err
}

func GetMatchHistoryByUserIdAndTimeStampRange(userId int, startRange time.Time, endRange time.Time) (MatchHistory, error) {
	var matchHistory MatchHistory
	row := db.QueryRow("SELECT * FROM match_history WHERE match_history.user_id = ? AND match_history.game_start >= ? AND match_history.game_end <= ?", userId, startRange, endRange)

	err := row.Scan(&matchHistory.Id, &matchHistory.UserId, &matchHistory.GameStart, &matchHistory.GameEnd, &matchHistory.Map, &matchHistory.Legend, &matchHistory.GameMode)
	return matchHistory, err
}

func GetUserByApexUid(uid string) (User, error) {
	var user User

	row := db.QueryRow("SELECT * FROM users WHERE users.apex_uid = ?", uid)
	err := row.Scan(&user.Id, &user.Name, &user.ApexUsername, &user.ApexUid)

	return user, err
}

func GetMapByAlsName(alsName string) (Map, error) {
	var gameMap Map
	row := db.QueryRow("SELECT * FROM maps WHERE maps.als_name = ?", alsName)
	err := row.Scan(&gameMap.Id, &gameMap.Name, &gameMap.CardImage, &gameMap.AlsName)

	return gameMap, err
}

func GetLegendByName(name string) (Legend, error) {
	var legend Legend
	row := db.QueryRow("SELECT * FROM legends WHERE legends.name = ?", name)
	err := row.Scan(&legend.Id, &legend.Name, &legend.CardImage)

	return legend, err
}

func AddNewMatchHistory(matchHistory MatchHistory) error {
	_, err := db.Exec("INSERT INTO match_history (user_id, game_start, game_end, map, legend, game_mode) VALUES (?, ?, ?, ?, ?, ?)", matchHistory.UserId, matchHistory.GameStart, matchHistory.GameEnd, matchHistory.Map, matchHistory.Legend, matchHistory.GameMode)
	return err
}

func GetMatchHistoriesForClip(clip Clip) ([]MatchHistory, error) {
	var matchHistories []MatchHistory
	clipStart := clip.CreatedAt.Time.Add(time.Duration(clip.Duration*-1) * time.Second)
	rows, err := db.Query("SELECT * FROM match_history WHERE user_id = ? AND ((game_start <= ? AND game_end >= ?) OR (game_start <= ? AND game_end >= ?) OR (? <= game_start AND ? >= game_end) OR (? <= game_start AND ? >= game_end))", clip.OwnerId, clipStart, clipStart, clip.CreatedAt, clip.CreatedAt, clipStart, clipStart, clip.CreatedAt, clip.CreatedAt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var matchHistory MatchHistory
		if err := rows.Scan(&matchHistory.Id, &matchHistory.UserId, &matchHistory.GameStart, &matchHistory.GameEnd, &matchHistory.Map, &matchHistory.Legend, &matchHistory.GameMode); err != nil {
			return nil, err
		}
		matchHistories = append(matchHistories, matchHistory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return matchHistories, nil
}

func UpdateQueueEntryStatusToQueued(clipId int) error {
	_, err := db.Exec("UPDATE clips_queue SET clips_queue.status = 'queued' WHERE clips_queue.clip_id = ?", clipId)
	return err
}

func UpdateQueueEntryStatusToTranscoding(clipId int) error {
	_, err := db.Exec("UPDATE clips_queue SET clips_queue.status = 'transcoding', clips_queue.started_at = ? WHERE clips_queue.clip_id = ?", time.Now(), clipId)
	return err
}

func UpdateQueueEntryStatusToFinished(clipId int) error {
	_, err := db.Exec("UPDATE clips_queue SET clips_queue.status = 'finished', clips_queue.started_at = ? WHERE clips_queue.clip_id = ?", time.Now(), clipId)
	return err
}

func UpdateQueueEntryStatusToError(clipId int, errorMessage string) error {
	_, err := db.Exec("UPDATE clips_queue SET clips_queue.status = 'error', clips_queue.finished_at = ?, clips_queue.error_message = ? WHERE clips_queue.clip_id = ?", time.Now(), errorMessage, clipId)
	return err
}

func UpdateClipOnTranscodeFinish(clipId int, durationSeconds float64) error {
	_, err := db.Exec("UPDATE clips SET clips.is_processed = 1, clips.duration = ? WHERE clips.id = ?", durationSeconds, clipId)
	return err
}

func GetAllPendingQueueEntries() ([]QueueEntry, error) {
	var queueEntries []QueueEntry

	rows, err := db.Query("SELECT clips_queue.*,clips.filename FROM clips_queue INNER JOIN clips on clips_queue.clip_id = clips.id where status = 'pending'")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var queueEntry QueueEntry
		if err := rows.Scan(&queueEntry.Id, &queueEntry.ClipId, &queueEntry.Status, &queueEntry.StartedAt, &queueEntry.FinishedAt, &queueEntry.Operation); err != nil {
			return nil, err
		}
		queueEntries = append(queueEntries, queueEntry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return queueEntries, nil
}
