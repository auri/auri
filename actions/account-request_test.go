package actions

import (
	"auri/mail"
	"auri/models"
)

func (as *ActionSuite) Test_AccountRequestPageOk() {
	res := as.HTML("/").Get()
	as.Equal(200, res.Code, res.Body.String())
}

func (as *ActionSuite) Test_AccountRequestPageWithComment() {
	res := as.HTML("/?comment=testcomment").Get()
	body := res.Body.String()
	as.Equal(200, res.Code, body)
	as.Contains(body, ">testcomment</textarea>")
}

func (as *ActionSuite) Test_AccountRequestSubmitValidData() {
	defer mail.ClearMock()
	request := models.Request{Email: "user@example.com", Name: "test", LastName: "test"}
	req := as.JSON("/accountrequest")
	reg := req.Post(&request)
	as.Equal(201, reg.Code, reg.Body.String())
}

func (as *ActionSuite) Test_AccountRequestSubmitEmptyEmail() {
	request := models.Request{Email: ""}
	req := as.JSON("/accountrequest")
	reg := req.Post(request)
	as.Equal(400, reg.Code, reg.Body.String())
}

func (as *ActionSuite) Test_AccountRequestSubmitEmptyName() {
	request := models.Request{Email: "example@example.com", Name: "", LastName: "test"}
	req := as.JSON("/accountrequest")
	reg := req.Post(request)
	as.Equal(400, reg.Code, reg.Body.String())
}

func (as *ActionSuite) Test_AccountRequestSubmitEmptyLastName() {
	request := models.Request{Email: "example@example.com", Name: "test", LastName: ""}
	req := as.JSON("/accountrequest")
	reg := req.Post(request)
	as.Equal(400, reg.Code, reg.Body.String())
}

func (as *ActionSuite) Test_AccountRequestSubmitInvalidEmailDomain() {
	request := models.Request{Email: "example@invalid.example.com", Name: "test", LastName: "test"}
	req := as.JSON("/accountrequest")
	reg := req.Post(request)
	as.Equal(400, reg.Code, reg.Body.String())
}
