package ssh

import (
	"reflect"
	"testing"
)

func Test_ConvertPuttySSH(t *testing.T) {
	tests := []struct {
		name      string
		inputKey  string
		outputKey string
		wantErr   string // empty for no errors
	}{
		{"proper rsa 2k putty key",
			`---- BEGIN SSH2 PUBLIC KEY ----
Comment: "rsa-key-20210826"
AAAAB3NzaC1yc2EAAAABJQAAAQEA7h//wvJpdAScBrlu6Vop3LdFYSrgTfumOkEm
MKTUeG4yxwVk+HiBwwVt0BWDvXbYVcoL+lourxD1CBAPk5TesHqCXOdhEpWx8KMA
jeuBzqDCp/OBIXURg60LwjCzGDZ/TFz+FXx/M0OXjPfulSlaTkX3Mcc4+nTDpGCr
73eQWt2VMGQBNZENzQYrmQaP2QRAdA7CU4PJNC1CGrx+J3NnkFnhmyDDsu4m800P
4HkyKqUayxkVtXl59av4LsYAT0VV4U7WNRLUOFFLeUh3pJLfiGky/IJpvUV1yIT7
Dvz0bSsRN7/hHMn8wc29KLl3XF10XJJSnuOGKAo0yY56DLziew==
---- END SSH2 PUBLIC KEY ----`,
			"ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEA7h//wvJpdAScBrlu6Vop3LdFYSrgTfumOkEmMKTUeG4yxwVk+HiBwwVt0BWDvXbYVcoL+lourxD1CBAPk5TesHqCXOdhEpWx8KMAjeuBzqDCp/OBIXURg60LwjCzGDZ/TFz+FXx/M0OXjPfulSlaTkX3Mcc4+nTDpGCr73eQWt2VMGQBNZENzQYrmQaP2QRAdA7CU4PJNC1CGrx+J3NnkFnhmyDDsu4m800P4HkyKqUayxkVtXl59av4LsYAT0VV4U7WNRLUOFFLeUh3pJLfiGky/IJpvUV1yIT7Dvz0bSsRN7/hHMn8wc29KLl3XF10XJJSnuOGKAo0yY56DLziew== rsa-key-20210826",
			""},
		{"proper rsa 2k putty key with line break",
			`---- BEGIN SSH2 PUBLIC KEY ----
Comment: "rsa-key-20210826"
AAAAB3NzaC1yc2EAAAABJQAAAQEA7h//wvJpdAScBrlu6Vop3LdFYSrgTfumOkEm
MKTUeG4yxwVk+HiBwwVt0BWDvXbYVcoL+lourxD1CBAPk5TesHqCXOdhEpWx8KMA
jeuBzqDCp/OBIXURg60LwjCzGDZ/TFz+FXx/M0OXjPfulSlaTkX3Mcc4+nTDpGCr
73eQWt2VMGQBNZENzQYrmQaP2QRAdA7CU4PJNC1CGrx+J3NnkFnhmyDDsu4m800P
4HkyKqUayxkVtXl59av4LsYAT0VV4U7WNRLUOFFLeUh3pJLfiGky/IJpvUV1yIT7
Dvz0bSsRN7/hHMn8wc29KLl3XF10XJJSnuOGKAo0yY56DLziew==
---- END SSH2 PUBLIC KEY ----
`,
			"ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEA7h//wvJpdAScBrlu6Vop3LdFYSrgTfumOkEmMKTUeG4yxwVk+HiBwwVt0BWDvXbYVcoL+lourxD1CBAPk5TesHqCXOdhEpWx8KMAjeuBzqDCp/OBIXURg60LwjCzGDZ/TFz+FXx/M0OXjPfulSlaTkX3Mcc4+nTDpGCr73eQWt2VMGQBNZENzQYrmQaP2QRAdA7CU4PJNC1CGrx+J3NnkFnhmyDDsu4m800P4HkyKqUayxkVtXl59av4LsYAT0VV4U7WNRLUOFFLeUh3pJLfiGky/IJpvUV1yIT7Dvz0bSsRN7/hHMn8wc29KLl3XF10XJJSnuOGKAo0yY56DLziew== rsa-key-20210826",
			""},
		{"proper rsa 2k putty key with special comment",
			`---- BEGIN SSH2 PUBLIC KEY ----
Comment: rsa-key-20210826
AAAAB3NzaC1yc2EAAAABJQAAAQEA7h//wvJpdAScBrlu6Vop3LdFYSrgTfumOkEm
MKTUeG4yxwVk+HiBwwVt0BWDvXbYVcoL+lourxD1CBAPk5TesHqCXOdhEpWx8KMA
jeuBzqDCp/OBIXURg60LwjCzGDZ/TFz+FXx/M0OXjPfulSlaTkX3Mcc4+nTDpGCr
73eQWt2VMGQBNZENzQYrmQaP2QRAdA7CU4PJNC1CGrx+J3NnkFnhmyDDsu4m800P
4HkyKqUayxkVtXl59av4LsYAT0VV4U7WNRLUOFFLeUh3pJLfiGky/IJpvUV1yIT7
Dvz0bSsRN7/hHMn8wc29KLl3XF10XJJSnuOGKAo0yY56DLziew==
---- END SSH2 PUBLIC KEY ----`,
			"ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEA7h//wvJpdAScBrlu6Vop3LdFYSrgTfumOkEmMKTUeG4yxwVk+HiBwwVt0BWDvXbYVcoL+lourxD1CBAPk5TesHqCXOdhEpWx8KMAjeuBzqDCp/OBIXURg60LwjCzGDZ/TFz+FXx/M0OXjPfulSlaTkX3Mcc4+nTDpGCr73eQWt2VMGQBNZENzQYrmQaP2QRAdA7CU4PJNC1CGrx+J3NnkFnhmyDDsu4m800P4HkyKqUayxkVtXl59av4LsYAT0VV4U7WNRLUOFFLeUh3pJLfiGky/IJpvUV1yIT7Dvz0bSsRN7/hHMn8wc29KLl3XF10XJJSnuOGKAo0yY56DLziew== rsa-key-20210826",
			""},
		{"proper rsa 2k putty key without comment",
			`---- BEGIN SSH2 PUBLIC KEY ----
AAAAB3NzaC1yc2EAAAABJQAAAQEA7h//wvJpdAScBrlu6Vop3LdFYSrgTfumOkEm
MKTUeG4yxwVk+HiBwwVt0BWDvXbYVcoL+lourxD1CBAPk5TesHqCXOdhEpWx8KMA
jeuBzqDCp/OBIXURg60LwjCzGDZ/TFz+FXx/M0OXjPfulSlaTkX3Mcc4+nTDpGCr
73eQWt2VMGQBNZENzQYrmQaP2QRAdA7CU4PJNC1CGrx+J3NnkFnhmyDDsu4m800P
4HkyKqUayxkVtXl59av4LsYAT0VV4U7WNRLUOFFLeUh3pJLfiGky/IJpvUV1yIT7
Dvz0bSsRN7/hHMn8wc29KLl3XF10XJJSnuOGKAo0yY56DLziew==
---- END SSH2 PUBLIC KEY ----`,
			"ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEA7h//wvJpdAScBrlu6Vop3LdFYSrgTfumOkEmMKTUeG4yxwVk+HiBwwVt0BWDvXbYVcoL+lourxD1CBAPk5TesHqCXOdhEpWx8KMAjeuBzqDCp/OBIXURg60LwjCzGDZ/TFz+FXx/M0OXjPfulSlaTkX3Mcc4+nTDpGCr73eQWt2VMGQBNZENzQYrmQaP2QRAdA7CU4PJNC1CGrx+J3NnkFnhmyDDsu4m800P4HkyKqUayxkVtXl59av4LsYAT0VV4U7WNRLUOFFLeUh3pJLfiGky/IJpvUV1yIT7Dvz0bSsRN7/hHMn8wc29KLl3XF10XJJSnuOGKAo0yY56DLziew==",
			""},
		{"broken putty key",
			`ABC`,
			"",
			"Invalid putty public key"},
		{"broken putty key #2",
			`---- BEGIN SSH2 PUBLIC KEY ----
Comment: "rsa-key-20210826"
AAAAB3NzaC1yc2EAAAABJQAAAQEA7h//wvJpdAScBrlu6Vop3LdFYSrgTfumOkEm
MKTUeG4yxwVk+HiBwwVt0BWDvXbYVcoL+lourxD1CBAPk5TesHqCXOdhEpWx8KMA
jeuBzqDCp/OBIXURg60LwjCzGDZ/TFz+FXx/M0OXjPfulSlaTkX3Mcc4+nTDpGCr
73eQWt2VMGQBNZENzQYrmQaP2QRAdA7CU4PJNC1CGrx+J3NnkFnhmyDDsu4m800P
4HkyKqUayxkVtXl59av4LsYAT0VV4U7WNRLUOFFLeUh3pJLfiGky/IJpvUV1yIT7
Dvz0bSsRN7/hHMn8wc29KLl3XF10XJJSnuOGKAo0yY56DLziew==
---- END SSH2 PUBLIC KEY --`,
			"",
			"Invalid putty public key"},
		{"proper ed25519 putty key",
			`---- BEGIN SSH2 PUBLIC KEY ----
Comment: "ed25519-key-20210826"
AAAAC3NzaC1lZDI1NTE5AAAAIALoF39Bi+IqrjGnRdXSZRA8ih/FcB3NXWamTLLu
o4uJ
---- END SSH2 PUBLIC KEY ----`,
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIALoF39Bi+IqrjGnRdXSZRA8ih/FcB3NXWamTLLuo4uJ ed25519-key-20210826",
			""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputKey, err := ConvertPuttySSH(tt.inputKey)

			if (err != nil) != (tt.wantErr != "") {
				t.Errorf("ConvertPuttySSH() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr != "" && (err != nil) {
				if !reflect.DeepEqual(err.Error(), tt.wantErr) {
					t.Errorf("ConvertPuttySSH() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if !reflect.DeepEqual(outputKey, tt.outputKey) {
					t.Errorf("ConvertPuttySSH = \n%#v, want \n%#v", outputKey, tt.outputKey)
				}
			}
		})
	}
}
