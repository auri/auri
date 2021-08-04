package models

func (ms *ModelSuite) Test_RequestValidateErr() {
	u := Request{}
	ms.Error(u.Validate(DB))
	u.Email = "Foo"
	ms.Error(u.Validate(DB))
}

func (ms *ModelSuite) Test_RequestValidateErrNoName() {
	u := Request{}
	ms.Error(u.Validate(DB))
	u.Name = ""
	ms.Error(u.Validate(DB))
}

func (ms *ModelSuite) Test_RequestValidateErrNoLastName() {
	u := Request{}
	ms.Error(u.Validate(DB))
	u.LastName = ""
	ms.Error(u.Validate(DB))
}

func (ms *ModelSuite) Test_RequestValidateNameOk() {
	u := Request{}
	ms.Error(u.Validate(DB))
	u.Name = "test"
	ms.Error(u.Validate(DB))
}

func (ms *ModelSuite) Test_RequestValidateLastNameOk() {
	u := Request{}
	ms.Error(u.Validate(DB))
	u.Name = "test"
	ms.Error(u.Validate(DB))
}

func (ms *ModelSuite) Test_PublicKeyNoKey() {
	u := Request{}
	ms.Error(u.Validate(DB))
	u.PublicKey = ""
	ms.Error(u.Validate(DB))
}

func (ms *ModelSuite) Test_RequestValidateOK() {
	u := Request{}
	u.Email = "user@example.com"
	ms.Error(u.Validate(DB))
	u.Name = "test"
	ms.Error(u.Validate(DB))
	u.LastName = "test"
	ms.Error(u.Validate(DB))
	u.PublicKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIK0wmN/Cr3JXqmLW7u+g9pTh+wyqDHpSQEIQczXkVx9q user@example.com"
	ms.Error(u.Validate(DB))
	checkError, err := u.Validate(DB)
	ms.NoError(err)
	ms.False(checkError.HasAny())
}

func (ms *ModelSuite) Test_RequestCreateOK() {
	u := Request{}
	u.Email = "user@example.com"
	ms.Error(u.Validate(DB))
	u.Name = "test"
	ms.Error(u.Validate(DB))
	u.LastName = "test"
	ms.Error(u.Validate(DB))
	u.PublicKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIK0wmN/Cr3JXqmLW7u+g9pTh+wyqDHpSQEIQczXkVx9q user@example.com"
	ms.Error(u.Validate(DB))
	checkError, err := u.Create(DB)
	ms.NoError(err)
	ms.False(checkError.HasAny())
}

func (ms *ModelSuite) Test_RequestExists() {
	u := Request{
		Email:     "user@example.com",
		Name:      "test",
		LastName:  "test",
		PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIK0wmN/Cr3JXqmLW7u+g9pTh+wyqDHpSQEIQczXkVx9q user@example.com",
	}
	verrs, err := u.Create(DB)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	u = Request{
		Email:     "user@example.com",
		Name:      "test",
		LastName:  "test",
		PublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIK0wmN/Cr3JXqmLW7u+g9pTh+wyqDHpSQEIQczXkVx9q user@example.com",
	}
	verrs, err = u.Create(DB)
	ms.NoError(err)
	ms.True(verrs.HasAny())
}

func (ms *ModelSuite) Test_RequestExistsInIPA() {
	u := Request{
		Email: "user-present-in-IPA@example.com",
	}
	verrs, err := u.Create(DB)
	ms.NoError(err)
	ms.True(verrs.HasAny())
}

func (ms *ModelSuite) Test_RequestLogin() {
	u := Request{
		Name:     "Abc",
		LastName: "Xyz",
	}

	ms.Equal("axyz", u.Login(0))
	ms.Equal("axyz1", u.Login(1))

	umlauts := Request{
		Name:     "ÄäÜüÖößabC",
		LastName: "XyzÄäÜüÖößxyZ",
	}

	ms.Equal("aexyzaeaeueueoeoessxyz", umlauts.Login(0))
	ms.Equal("aexyzaeaeueueoeoessxyz1", umlauts.Login(1))
}
