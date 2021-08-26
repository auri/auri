package ssh

import (
	"reflect"
	"testing"
)

func Test_DetermineType(t *testing.T) {
	tests := []struct {
		name     string
		inputKey string
		keyType  KeyType
		length   int
		wantErr  string // empty for no errors
	}{
		{"valid RSA key",
			"ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEA7h//wvJpdAScBrlu6Vop3LdFYSrgTfumOkEmMKTUeG4yxwVk+HiBwwVt0BWDvXbYVcoL+lourxD1CBAPk5TesHqCXOdhEpWx8KMAjeuBzqDCp/OBIXURg60LwjCzGDZ/TFz+FXx/M0OXjPfulSlaTkX3Mcc4+nTDpGCr73eQWt2VMGQBNZENzQYrmQaP2QRAdA7CU4PJNC1CGrx+J3NnkFnhmyDDsu4m800P4HkyKqUayxkVtXl59av4LsYAT0VV4U7WNRLUOFFLeUh3pJLfiGky/IJpvUV1yIT7Dvz0bSsRN7/hHMn8wc29KLl3XF10XJJSnuOGKAo0yY56DLziew== rsa-key-20210826",
			RSAKey,
			2048,
			""},
		{"invalid RSA key",
			"ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEA7h//wvJpdAScBrlu6Vop3LdFYSrgTfumOkEmMKTUeG4yxwVk+HiBwwVt0BWDvXbYVcoL+lourxD1CBAPk5TesHqCXOdhEpWx8KMAjeuBzqDx+J3NnkFnhmyDDsu4m800P4HkyKqUayxkVtXl59av4LsYAT0VV4U7WNRLUOFFLeUh3pJLfiGky/IJpvUV1yIT7Dvz0bSsRN7/hHMn8wc29KLl3XF10XJJSnuOGKAo0yY56DLziew== rsa-key-20210826",
			0,
			0,
			"Invalid key format"},
		{"valid ed25519 key",
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIALoF39Bi+IqrjGnRdXSZRA8ih/FcB3NXWamTLLuo4uJ ed25519-key-20210826",
			Ed25519Key,
			0,
			""},
		{"invalid ed25519 key",
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5LLuo4uJ ed25519-key-20210826",
			0,
			0,
			"Invalid key format"},
		{"valid DSA key",
			"ssh-dss AAAAB3NzaC1kc3MAAACBAKx4LXK33lfy61FAUJWzRTPc/GVIyTL/ACeCvtVkuP59SYZdnZ9jO79zfCnoP7XV8OwpVAuHcxV2u9jUQ+vXyXTO0w8cSNfJSw/8vihwMWP+1oU3eO6pQn7g2icVCORwAN9wgBIILJ8ethsXKJX+Kl7Gh1LD/XoBnKYKYA9yAg1TAAAAFQDd5bk14tIQEkhWCdc9YglZxtGTkwAAAIBvYN1pXgnVU74tyYTI2O3kuOM6EzWso8yhyV1TEQNLyEBhj7dKn1rlfaZOxcAXajfO5W6++Cgv5ph59QQdLhRCTvSwQUSbO6Aol42QzeXQDWQb/wfMRl7hevqxxUej35Zyfr1W9PjIN3zjqs3/K9hqcSVqI1w2G9kyHs5nt9X5OAAAAIBQUK9IDascEbs8W941T2qb7yMxDWU+oTV3+8Qa+SSehWu4T1aOAJaT1XpWcVSLpDAgevk7P05JfokC5W+MhoKHKpujPXDpfQD6e5R3htaEUsxPtVFS3sataJjwnhlMkyjWIQaUfmB2UI3XJ4sm0rjvMoLwVQ/eRuIQmEJ1NTWdZQ==",
			0,
			0,
			"Invalid key format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyType, length, err := DetermineType(tt.inputKey)

			if (err != nil) != (tt.wantErr != "") {
				t.Errorf("DetermineType() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != "" && (err != nil) {
				if !reflect.DeepEqual(err.Error(), tt.wantErr) {
					t.Errorf("DetermineType() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if !reflect.DeepEqual(keyType, tt.keyType) {
					t.Errorf("DetermineType = \n%#v, want \n%#v", keyType, tt.keyType)
				}

				if !reflect.DeepEqual(length, tt.length) {
					t.Errorf("DetermineType = \n%#v, want \n%#v", length, tt.length)
				}
			}
		})
	}
}
