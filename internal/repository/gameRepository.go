package data

import (
	"wordGame/internal/models"
	db "wordGame/pkg"
)

func LeaderBoardFill(user_id, category_id, score int) error { //Saving all data from the game to the leaderBoard
	_, err := db.GetDB().Exec("INSERT INTO game_results (user_id, category_id, score) VALUES ($1, $2, $3);", user_id, category_id, score)
	return err
}

func ShowLeaderBoardByCategory(catId int) ([]models.LeaderBoard, error) {
	rows, err := db.GetDB().Query(`
	SELECT u.id, u.username, SUM(r.score) AS total_score
	FROM game_results r
	JOIN users u ON u.id = r.user_id
	WHERE r.category_id = $1
	GROUP BY u.id, u.username
	ORDER BY total_score DESC
	LIMIT 10
`, catId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	leaders := make([]models.LeaderBoard, 0)

	for rows.Next() {
		var lb models.LeaderBoard

		err := rows.Scan(
			&lb.UserId,
			&lb.UserName,
			&lb.Score,
		)
		if err != nil {
			return nil, err
		}

		leaders = append(leaders, lb)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return leaders, nil
}
