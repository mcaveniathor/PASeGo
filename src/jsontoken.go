package past

type JsonToken struct {
	claims map[string]string
	iss    string
	sub    string
	aud    string
	exp    string
	nbf    string
	iat    string
	jti    string
	cached string
	footer string
	//key							SymmetricKey
}

func (t *JsonToken) get(claim string) (string, error) {
	if t.claims[claim] == nil {
		return nil, Errors.New("claim does not exist")
	}
	return claims[claim], nil
}

func (t *JsonToken) getClaims() []string {

}
