package postgres

import "chat/store/users"

func (s *Store) GetRecentChats(userID int64, chat users.RecentChat) error {
	query := `
        SELECT 
            sub.partner_id,
            u.username AS partner_username,
            sub.id,
            sub.content,
            sub.sender_id,
            sub.created_at,
            sub.is_read,
            sub.is_delivered
        FROM (
            SELECT DISTINCT ON (CASE WHEN sender_id = $1 THEN receiver_id ELSE sender_id END)
                id,
                content,
                sender_id,
                created_at,
                is_read,
                is_delivered,
                CASE WHEN sender_id = $1 THEN receiver_id ELSE sender_id END AS partner_id
            FROM messages
            WHERE sender_id = $1 OR receiver_id = $1
            ORDER BY (CASE WHEN sender_id = $1 THEN receiver_id ELSE sender_id END), created_at DESC
        ) sub
        JOIN users u ON u.id = sub.partner_id
        ORDER BY sub.created_at DESC;
    `

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		msg := &users.Message{}

		var partnerID int64
		var partnerUsername string

		if err := rows.Scan(
			&partnerID,
			&partnerUsername,
			&msg.ID,
			&msg.Content,
			&msg.Sender_id,
			&msg.Timestamp,
			&msg.Is_read,
			&msg.Is_delivered,
		); err != nil {
			return err
		}

		if err := chat(partnerID, partnerUsername, msg); err != nil {
			return err
		}
	}

	return nil
}
