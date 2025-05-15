package storage

import (
	"context"
	"database/sql"
)

type DeviceCodes struct {
	db *sql.DB
}

// SaveDeviceCode stores user code and corresponding device code and login, sets expiration time
func (d *DeviceCodes) SaveDeviceCode(ctx context.Context, deviceCode, userCode, login string, expiresIn int64) error {
	_, err := d.db.ExecContext(ctx, `
		INSERT INTO devices (device_code, user_code, login, expires_at, activated)
		VALUES ($1, $2, $3, NOW() + ($4 || ' seconds')::interval, false)
	`, deviceCode, userCode, login, expiresIn)
	if err != nil {
		return err
	}
	return nil
}

// ActivateDevice sets user code's 'activated' column to true
func (d *DeviceCodes) ActivateDevice(ctx context.Context, userCode string) error {
	result, err := d.db.ExecContext(ctx, `
	UPDATE devices
	SET activated = true
	WHERE user_code = $1 AND expires_at > NOW()
	`, userCode)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetDeviceLoginIfActivated returns login and checks if the user has already authenticated
func (d *DeviceCodes) GetDeviceLoginIfActivated(ctx context.Context, userCode string) (string, bool, error) {
	var login string
	var activated bool

	err := d.db.QueryRowContext(ctx, `
	SELECT login, activated
	FROM devices
	WHERE user_code = $1 AND expires_at > NOW()
	`, login).Scan(&login, &activated)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}

	return login, activated, nil
}
