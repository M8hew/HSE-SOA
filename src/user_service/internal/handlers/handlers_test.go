package handlers

import (
	"testing"

	"user_service/api"

	"github.com/stretchr/testify/require"
)

type mockDB struct {
	counter       int
	loginID       map[string]int
	loginPassword map[string][16]byte
}

func newMockDB() mockDB {
	return mockDB{
		counter:       0,
		loginID:       make(map[string]int),
		loginPassword: make(map[string][16]byte),
	}
}

func (w mockDB) getUserLogin(userID int) (login string, err error) {
	return "", nil
}

func (w mockDB) addNewUser(userLogin string, hashPassword [16]byte) (id int, err error) {
	w.loginPassword[userLogin] = hashPassword
	w.counter++
	w.loginID[userLogin] = w.counter
	return w.loginID[userLogin], nil
}

func (w mockDB) getUserPasswordId(userLogin string) ([]byte, int, error) {
	password := w.loginPassword[userLogin]
	return password[:], w.loginID[userLogin], nil
}

func (w mockDB) updateUser(userId int, userData api.PutUpdateJSONRequestBody) error {
	return nil
}

func TestHashPassword(t *testing.T) {
	for _, tc := range []struct {
		name     string
		password string
		salt     string
		result   [16]byte
	}{
		{
			name:     "Easy password no salt",
			password: "password123",
			result:   [16]byte{0x48, 0x2C, 0x81, 0x1D, 0xA5, 0xD5, 0xB4, 0xBC, 0x6D, 0x49, 0x7F, 0xFA, 0x98, 0x49, 0x1E, 0x38},
		},
		{
			name:     "Easy password with salt",
			password: "password123",
			salt:     "salt",
			result:   [16]byte{0xCF, 0xFB, 0xA2, 0x6E, 0xD2, 0x54, 0x8E, 0xD5, 0xD0, 0x9E, 0x29, 0x3B, 0x0B, 0x3C, 0x51, 0x7D},
		},
		{
			name:     "Strong password no salt",
			password: "yTl+Io@dhD*lV7^_",
			result:   [16]byte{0xFB, 0x66, 0xD1, 0x0F, 0x48, 0x17, 0xB1, 0x82, 0xE8, 0x28, 0x0C, 0x57, 0x40, 0xDB, 0xF6, 0x99},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			hash := hashPassword(tc.password, tc.salt)
			require.Equal(t, tc.result, hash)
		})
	}
}

func TestParseConfig(t *testing.T) {
	configPath := "../../build/config.yaml"

	config, err := parseYAMLConfig(configPath)

	require.Nil(t, err)
	require.Contains(t, config, "user_sessions_lifetime")
	require.Contains(t, config, "jwt_public_path")
	require.Contains(t, config, "jwt_private_path")

	require.Equal(t, config["user_sessions_lifetime"], 3600)
}
