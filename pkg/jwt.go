package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

type Signature string
type JWTToken string
type JOSEHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type AccessTokenClaims struct {
	Issuer              string            `json:"iss"`   // who/what issued token - `go_booking`
	IssuedAt            int               `json:"iat"`   // when token was issued
	ExpirationTime      int               `json:"exp"`   // when token expires
	Subject             string            `json:"sub"`   // who gets token UserID - SHOULD BE string ref:https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.2
	Role                string            `json:"roles"` // string of RoleId - compatible with OAuth2 ref:https://www.rfc-editor.org/rfc/rfc9068.html#name-roles
	OriginatingIdentity IPAddressIdentity `json:"orig"`  // custom claim - this is crucial for preventing malicious users from using tokens on another machine or another program on the same machine
}

type IPAddressIdentity struct {
	IP string `json:"ip"` // has information about Remote Address. Example: `127.0.0.1`
}

type RefreshTokenClaims struct {
	Issuer              string            `json:"iss"`  // who issued token
	IssuedAt            int               `json:"iat"`  // when token was issued
	ExpirationTime      int               `json:"exp"`  // when token expires
	Subject             string            `json:"sub"`  // who gets token UserID - SHOULD BE string ref:https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.2
	OriginatingIdentity IPAddressIdentity `json:"orig"` // custom claim - this is crucial for preventing malicious users from using tokens on another machine or another program on the same machine
}

// GenerateJWTAccessToken ref: https://datatracker.ietf.org/doc/html/rfc7519#section-3.1
func GenerateJWTAccessToken(header JOSEHeader, claims AccessTokenClaims, key string) (JWTToken, error) {
	// 1. Convert JOSEHeader struct to byte array, where the octets represent the UTF-8 representation
	headerJSON, JSONMarshallerError := json.Marshal(header)
	if JSONMarshallerError != nil {
		log.Println("pkg.JWT.GenerateJWTAccessToken(): error occured during converting JOSEHeader to JSON byte array. Details: ", JSONMarshallerError)
		return JWTToken(""), JSONMarshallerError
	}

	// 2. Convert AccessTokenClaims struct to byte array, where the octets represent the UTF-8 representation
	claimsJSON, JSONMarshallerError := json.Marshal(claims)
	if JSONMarshallerError != nil {
		log.Println("pkg.JWT.GenerateJWTAccessToken(): error occured during converting AccessTokenClaims to JSON byte array. Details: ", JSONMarshallerError)
		return JWTToken(""), JSONMarshallerError
	}

	// 3. Encode Base64URL header JSON
	encodedJOSEHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	// 4. Encode Base64URL claims JSON
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// 5. Sign header and payload(claims) using HS256 hash function and secret key
	signature := SignHeaderAndPayload(encodedJOSEHeader, encodedClaims, key)

	// 6. return concatenated with "." between encoded header, claims and hashed and encoded signature
	return JWTToken(encodedJOSEHeader + "." + encodedClaims + "." + signature), nil
}

func GenerateJWTRefreshToken(header JOSEHeader, claims RefreshTokenClaims, key string) (JWTToken, error) {
	// 1. Convert JOSEHeader struct to byte array, where the octets represent the UTF-8 representation
	headerJSON, JSONMarshallerError := json.Marshal(header)
	if JSONMarshallerError != nil {
		log.Println("pkg.JWT.GenerateJWTRefreshToken(): error occured during converting JOSEHeader to JSON byte array. Details: ", JSONMarshallerError)
		return JWTToken(""), JSONMarshallerError
	}

	// 2. Convert RefreshTokenClaims struct to byte array, where the octets represent the UTF-8 representation
	claimsJSON, JSONMarshallerError := json.Marshal(claims)
	if JSONMarshallerError != nil {
		log.Println("pkg.JWT.GenerateJWTRefreshToken(): error occured during converting RefreshTokenClaims to JSON byte array. Details: ", JSONMarshallerError)
		return JWTToken(""), JSONMarshallerError
	}

	// 3. Encode Base64URL header JSON
	encodedJOSEHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	// 4. Encode Base64URL claims JSON
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// 5. Sign header and payload(claims) using HS256 hash function and secret key
	signature := SignHeaderAndPayload(encodedJOSEHeader, encodedClaims, key)

	// 6. return concatenated with "." between encoded header, claims and hashed and encoded signature
	return JWTToken(encodedJOSEHeader + "." + encodedClaims + "." + signature), nil
}

func ExtractAccessTokenClaims(input string) (AccessTokenClaims, error) {
	// 9.   Otherwise, base64url decode the Message
	decodedData, decodingError := base64.RawURLEncoding.DecodeString(input)
	if decodingError != nil {
		log.Println("pkg.JWT.ExtractAccessTokenClaims(): error during decoding Base64RawURL second part of JWT token (AccessTokenClaims). Passed data: ", input)
		return AccessTokenClaims{}, decodingError
	}

	// 9. following the restriction that no line breaks, whitespace, or other additional characters have been used.
	decodedDataStr := string(decodedData)
	decodedDataAfterTrimming := strings.TrimSpace(decodedDataStr)
	if decodedDataStr != decodedDataAfterTrimming {
		log.Println("pkg.JWT.ExtractAccessTokenClaims(): second part of JWT token has line break, whitespace or tab character(s). Passed data: ", input)
		return AccessTokenClaims{}, errors.New("second part of JWT token has line break, whitespace or tab character(s). Passed data: " + input)
	}

	var claims AccessTokenClaims
	// 10.  Verify that the resulting octet sequence is a UTF-8-encoded
	//        representation of a completely valid JSON object conforming to
	//        RFC 7159 [RFC7159]; let the JWT AccessTokenClaims Set be this JSON object.
	unmarshallingError := json.Unmarshal(decodedData, &claims)
	if unmarshallingError != nil {
		log.Println("pkg.JWT.ExtractAccessTokenClaims(): error during unmarshalling Base64RawURL-decoded second part of JWT Token. Reasons: not AccessTokenClaims or has invalid data types. Passed data: ", decodedData)
		return AccessTokenClaims{}, unmarshallingError
	}

	return claims, nil
}

func ExtractRefreshTokenClaims(input string) (RefreshTokenClaims, error) {
	// 9.   Otherwise, base64url decode the Message
	decodedData, decodingError := base64.RawURLEncoding.DecodeString(input)
	if decodingError != nil {
		log.Println("pkg.JWT.ExtractRefreshTokenClaims(): error during decoding Base64RawURL second part of JWT token (AccessTokenClaims). Passed data: ", input)
		return RefreshTokenClaims{}, decodingError
	}

	// 9. following the restriction that no line breaks, whitespace, or other additional characters have been used.
	decodedDataStr := string(decodedData)
	decodedDataAfterTrimming := strings.TrimSpace(decodedDataStr)
	if decodedDataStr != decodedDataAfterTrimming {
		log.Println("pkg.JWT.ExtractRefreshTokenClaims(): second part of JWT token has line break, whitespace or tab character(s). Passed data: ", input)
		return RefreshTokenClaims{}, errors.New("second part of JWT token has line break, whitespace or tab character(s). Passed data: " + input)
	}

	var claims RefreshTokenClaims
	// 10.  Verify that the resulting octet sequence is a UTF-8-encoded
	//        representation of a completely valid JSON object conforming to
	//        RFC 7159 [RFC7159]; let the JWT AccessTokenClaims Set be this JSON object.
	unmarshallingError := json.Unmarshal(decodedData, &claims)
	if unmarshallingError != nil {
		log.Println("pkg.JWT.ExtractRefreshTokenClaims(): error during unmarshalling Base64RawURL-decoded second part of JWT Token. Reasons: not AccessTokenClaims or has invalid data types. Passed data: ", decodedData)
		return RefreshTokenClaims{}, unmarshallingError
	}

	return claims, nil
}

func SignHeaderAndPayload(encodedJOSEHeader string, encodedClaims string, key string) string {
	// create HS256 hasher
	sha256Hasher := hmac.New(sha256.New, []byte(key))

	// 1. Concatenate encoded header and claims
	concatenatedEncodedHeaderAndClaims := encodedJOSEHeader + "." + encodedClaims
	// 2. Apply HashFunc with key to hash encoded header and claims separated by .
	sha256Hasher.Write(
		[]byte(concatenatedEncodedHeaderAndClaims),
	)
	// 3. Get signature hash in byte array
	hashedHeaderAndClaims := sha256Hasher.Sum(nil)
	// 4. Get signature -> Encode Base64URL hashed hashedHeaderAndClaims
	signature := base64.RawURLEncoding.EncodeToString(hashedHeaderAndClaims)

	return signature
}

func ExtractJOSEHeader(input string) (JOSEHeader, error) {
	// 3. Base64url decode the Encoded JOSE JOSEHeader
	decodedData, decodingError := base64.RawURLEncoding.DecodeString(input)
	if decodingError != nil {
		log.Println("pkg.JWT.ExtractJOSEHeader(): error during decoding Base64RawURL first part of JWT token (JOSEHeader). Passed data: ", input)
		return JOSEHeader{}, decodingError
	}

	// 3. following the restriction that no line breaks, whitespace, or other additional characters have been used
	decodedDataStr := string(decodedData)
	if strings.Contains(decodedDataStr, "\n") == true || strings.Contains(decodedDataStr, " ") == true ||
		strings.Contains(decodedDataStr, "\t") == true {
		log.Println("pkg.JWT.ExtractJOSEHeader(): first part of JWT token has line break, whitespace or tab character(s). Passed data: ", input)
		return JOSEHeader{}, errors.New("first part of JWT token has line break, whitespace or tab character(s). Passed data: " + input)
	}

	var joseHeader JOSEHeader
	//  4.   Verify that the resulting octet sequence is a UTF-8-encoded
	//        representation of a completely valid JSON object conforming to
	//        RFC 7159 [RFC7159]; let the JOSE JOSEHeader be this JSON object.
	// 6.   Determine whether the JWT is a JWS or a JWE using any of the
	//        methods described in Section 9 of [JWE].
	//  8.   If the JOSE JOSEHeader contains a "cty" (content type) value of
	//        "JWT", then the Message is a JWT that was the subject of nested
	//        signing or encryption operations.  In this case, return to Step
	//        1, using the Message as the JWT.
	//  Answer to 6: NO `enc` field is determined in my implementation of JWT-token
	// Answer to 8: NO `cty` field is determined in my implementation of JWT-token
	unmarshallingError := json.Unmarshal(decodedData, &joseHeader)
	if unmarshallingError != nil {
		log.Println("pkg.JWT.ExtractJOSEHeader(): error during unmarshalling Base64RawURL-decoded first part of JWT Token. Reasons: not JOSE JOSEHeader or has invalid data types. Passed data: ", decodedData)
		return JOSEHeader{}, unmarshallingError
	}

	//  5.   Verify that the resulting JOSE JOSEHeader includes only parameters
	//        and values whose syntax and semantics are both understood and
	//        supported or that are specified as being ignored when not
	//        understood.
	// TODO add more algorithms
	if joseHeader.Algorithm != "HS256" || joseHeader.Type != "JWT" {
		log.Println("pkg.JWT.ExtractJOSEHeader(): error in JOSE JOSEHeader - has unsupported values in parameters. JOSE JOSEHeader: ", joseHeader)
		return JOSEHeader{}, fmt.Errorf("error in JOSE JOSEHeader - has unsupported values in parameters. JOSE JOSEHeader: %#v", joseHeader)
	}

	return joseHeader, nil
}

func (h JOSEHeader) String() string {
	return fmt.Sprintf(`{"alg":"%s","typ":"%s"}`, h.Algorithm, h.Type)
}

func (i IPAddressIdentity) String() string {
	return fmt.Sprintf(`{"ip":"%s"}`, i.IP)
}
