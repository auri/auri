package models

import (
	"auri/config"

	"github.com/gobuffalo/envy"
)

func (ms *ModelSuite) Test_PublicKeyRSA2048Err() {
	request := Request{PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC2vjjBwydqP7mHutYi957wU7ezEu4WR/fEsmKQUy+nnhF4zl7Y+wiLWKrlPk08ziLQPHCJV/9vzE/fM86BwAVqUmtCzHUYtiNxGHAJ6mVYZkkM9/gkQa+XY3TA1w2ToMZv7OsMjeBLbdLTHdpOL2BTmRLQB+QmcA9q2JMK5wjgQ4rhf6rbYyMuIu3aOyueuoPcgNAqwOFsyoR2/ucjCv6+7PVi/3uafhYb+oKnwpPkIzfuows5NbhmbYXpn0IVRBO68tNlk0tONaoPtjvvR2knOX1T7+V/CYRYrxPPabdMOMEqlXA47sQ2Tq8FKYhaCzvs7FuKIiptWzjkdTBv2PK1"}
	validate := request.ValidateSSHOrPassword()
	ms.True(validate.HasAny())
	ms.Equal("RSA key is too short, it should be at least 4096 bit", validate.Error())
}

func (ms *ModelSuite) Test_PublicKeyRSANotAllowed() {
	envy.Temp(func() {
		request := Request{PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCdhXwVZAe7NwGTWjFsQmL8deCS5ZjtvVVxnrnu1UpYBKZph50ElXwLQ839v8CtSVR1lLS1KpIYg1i2vl42NbY6+q2G/pMS1rUVVrGqbycP04F765EcG78jP0byPoWpEoMwv7oP8w5gGR4m3+i81St1nLBczbMfomB3bpiBAH7jZZu+VXRPT88AJzJ4UiSBxTloa9i8mG6bMeYT+wxXEpv2H+YhEfxXoxQFXvvr1Fd3RX4y3HzGk/Oum+DXbF3cQqQ6lYndmBRjYzlY9jLfFhxVXaH/TYbJVQY4k4zcilHMMXnQp/p/vXEpyNZw7rTMxV3IcYl+ihE8FKZ/PiyJp872uCNpCoo7UZhYsbMggT4KwNTrZKnenB1aMA9fKE43zpaIsiJ60WS6JJj/cYsQzlOcXFGzpD6hTzY7D7GLMxoKck9nnZ2LO39x3mRRyjamhOAtVMkT6IJ5y3F+dj0xiX73KcxAZ0YpYKK1b9V8AVQcidGpjS3ISLvo3hDW1H1mt3OPzpl0rcU9lV0oGqTYJsrC3mcU6kv0Abxr94voB/4ur2nYTv6KseOoH2sBsKwR/L2ClBg/wvto6DtRnG9hQQTRPwdGaow/FipcAi5hm2ErdH1fg4DEhyB0CE7C5oVq0lIf/o2/vr5OdRvCjqDWFmDTkI2zmlq7aW4/eVJlE+HYww=="}
		envy.Set("ALLOW_RSA_KEYS", "false")
		ms.NoError(config.LoadConfig())
		validate := request.ValidateSSHOrPassword()
		ms.True(validate.HasAny())
		ms.Equal("RSA key is not allowed", validate.Error())
	})

}

func (ms *ModelSuite) Test_PublicKeyED25519NotAllowed() {
	envy.Temp(func() {
		request := Request{PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIK0wmN/Cr3JXqmLW7u+g9pTh+wyqDHpSQEIQczXkVx9q test@test.de"}
		envy.Set("ALLOW_ED25519_KEYS", "false")
		ms.NoError(config.LoadConfig())
		validate := request.ValidateSSHOrPassword()
		ms.True(validate.HasAny())
		ms.Equal("ED25519 key is not allowed", validate.Error())
	})
}

//Validate ssh_key with new line
func (ms *ModelSuite) Test_PublicKeyOkWithNewLine() {
	envy.Temp(func() {
		u := Request{PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIK0wmN/Cr3JXqmLW7u+g9pTh+wyqDHpSQEIQczXkVx9q test@test.de" + "\n"}
		ms.NoError(config.LoadConfig())
		checkError := u.ValidateSSHOrPassword()
		ms.False(checkError.HasAny())
	})
}

//Not allowed symbols are (! and $)
func (ms *ModelSuite) Test_PublicKeyNotAllowedSymbols() {
	u := Request{}
	ms.Error(u.Validate(DB))
	u.PublicKey = "ssh-ed25519 !AAAAC3NzaC1lZDI1NTE5AAAAEjkPQpaaFHDTIdUJVx8jA4esihdfioaalmWXUlY0QK$ test@tes.DE"
	ms.Error(u.ValidateSSHOrPassword())
}
func (ms *ModelSuite) Test_PublicKeyWrongFormat() {
	u := Request{}
	ms.Error(u.Validate(DB))
	u.PublicKey = "ssh-public AAAAC3NzaC1lZDI1NTE5AAAAEjkPQpaaFHDTIdUJVx8jA4esihdfioaalmWXUlY0QK test@tes.DE"
	ms.Error(u.ValidateSSHOrPassword())
}

func (ms *ModelSuite) Test_CreatePwShortLength() {
	request := Request{Password: "Short1!", PasswordConfirmation: "Short1!"}
	checkError := request.ValidateSSHOrPassword()
	ms.True(checkError.HasAny())
	ms.Equal("Password should be at least 12 characters", checkError.Error())
}
func (ms *ModelSuite) Test_CreatePwNoMatch() {
	request := Request{Password: "Test123!testtest", PasswordConfirmation: "Pass4321*"}
	checkError := request.ValidateSSHOrPassword()
	ms.True(checkError.HasAny())
	ms.Equal("Password and it's confirmation are not equal", checkError.Error())
}

func (ms *ModelSuite) Test_CreatePwPwdcEmpty() {
	request := Request{Password: "Test123!testtest", PasswordConfirmation: ""}
	checkError := request.ValidateSSHOrPassword()
	ms.True(checkError.HasAny())
	ms.Equal("Password and it's confirmation are not equal", checkError.Error())
}

func (ms *ModelSuite) Test_CreatePwNOK2ConditionsWereMetNoNumbersNoSymbole() {
	request := Request{Password: "Testtestfootesttest", PasswordConfirmation: "Testtestfootesttest"}
	checkError := request.ValidateSSHOrPassword()
	ms.True(checkError.HasAny())
	ms.Equal("Password do not meet 3 of 4 categories", checkError.Error())

}
func (ms *ModelSuite) Test_CreatePwNOK2ConditionsWereMetNoLowerCaseNoUpercase() {
	request := Request{Password: "123456789123456789!!", PasswordConfirmation: "123456789123456789!!"}
	checkError := request.ValidateSSHOrPassword()
	ms.True(checkError.HasAny())
	ms.Equal("Password do not meet 3 of 4 categories", checkError.Error())
}

func (ms *ModelSuite) Test_CreatePwOK3ConditionsWereMetNoUpperCaseLetter() {
	request := Request{Password: "test123!testtest", PasswordConfirmation: "test123!testtest"}
	checkError := request.ValidateSSHOrPassword()
	ms.False(checkError.HasAny())
}
func (ms *ModelSuite) Test_CreatePwOK3ConditionsWereMetNoLowerCaseLetter() {
	request := Request{Password: "TEST123!TESTTEST", PasswordConfirmation: "TEST123!TESTTEST"}
	checkError := request.ValidateSSHOrPassword()
	ms.False(checkError.HasAny())
}
func (ms *ModelSuite) Test_CreatePwOK3ConditionsWereMetNoNumber() {
	request := Request{Password: "TESTtest!testtest", PasswordConfirmation: "TESTtest!testtest"}
	checkError := request.ValidateSSHOrPassword()
	ms.False(checkError.HasAny())
}
func (ms *ModelSuite) Test_CreatePwOK3ConditionsWereMetNoSymbol() {
	request := Request{Password: "TESTtest123testtest", PasswordConfirmation: "TESTtest123testtest"}
	checkError := request.ValidateSSHOrPassword()
	ms.False(checkError.HasAny())
}
