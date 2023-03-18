package core

// ------------------------- Params -------------------------

type SendVerificationCodeMessage struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Code  string `json:"code"`
}

type SendResetPasswordTokenMessage struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type SendNewSignInSessionMessage struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	ClientIP  string `json:"client-ip"`
	UserAgent string `json:"user-agent"`
}
