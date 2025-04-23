package services

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"humoBooking/internal/database"
	"humoBooking/internal/models"
	"humoBooking/pkg"
	"log"
	"strconv"
	"strings"
	"time"
)

type AuthService struct {
	userRepository    database.UserRepository
	roleService       RoleServiceInterface
	routeService      RouteServiceInterface
	scopeService      ScopeServiceInterface
	permissionService PermissionServiceInterface

	// to perform IsOwner check
	bookingService BookingServiceInterface
	roomService    RoomServiceInterface
}

const (
	salt            = "eyJhbGciOiJIUzI1NiIsInR5ed217a32a94ba416f88e16122278cCI6IkpXVCJ9"
	accessTokenKey  = "062839dc6e3f934d4ed217a32a94ba416f88e161222785ad95803fe4923dd06b"
	refreshTokenKey = "0bf4586851bb6b6b15376e7b6bff4ac4d5cee836321f349462f60e3dbb07d7a4"
	issuer          = "humo_booking"
	JWT             = "JWT"
	HS256           = "HS256"
	accessTokenTTL  = 1 * time.Hour
	refreshTokenTTL = 3 * time.Hour
)

func NewAuthService(repository database.UserRepository, roleService RoleServiceInterface, routeService RouteServiceInterface, scopeService ScopeServiceInterface, permissionService PermissionServiceInterface, bookingService BookingServiceInterface, roomService RoomServiceInterface) *AuthService {
	return &AuthService{
		userRepository:    repository,
		roleService:       roleService,
		routeService:      routeService,
		scopeService:      scopeService,
		permissionService: permissionService,
		bookingService:    bookingService,
		roomService:       roomService,
	}
}

func (a *AuthService) CheckIfUserExistsAndPasswordIsCorrect(username string, password string) (models.User, error) {
	// Identification
	foundUser, err := a.userRepository.GetUserByUsername(username)
	if err != nil {
		log.Println("AuthService.CheckIfUserExistsAndPasswordIsCorrect(): error occured during User search. Passed data: ", username)
		return models.User{}, fmt.Errorf(`error occured during User search. Passed data: '%s'`, username)
	}

	emptyUser := models.User{}
	if foundUser == emptyUser {
		log.Println("AuthService.CheckIfUserExistsAndPasswordIsCorrect(): user not found! Passed data: ", username)
		return models.User{}, fmt.Errorf(`user not found. Passed data: '%s'`, username)
	}

	// Authentication
	userPasswordHash := foundUser.Password
	passwordHash := a.GeneratePasswordHash(password)
	if !strings.EqualFold(userPasswordHash, passwordHash) {
		log.Println("AuthService.CheckIfUserExistsAndPasswordIsCorrect(): wrong password!")
		return models.User{}, fmt.Errorf("wrong password")
	}

	return foundUser, nil
}

// GenerateTokens TODO Add error return value
func (a *AuthService) GenerateTokens(user models.User, identity pkg.IPAddressIdentity) (accessToken pkg.JWTToken, refreshToken pkg.JWTToken) {
	joseHeader := pkg.JOSEHeader{
		Algorithm: HS256,
		Type:      JWT,
	}

	now := time.Now() // extra measure for checking

	accessTokenClaims := pkg.AccessTokenClaims{
		Issuer:              issuer,
		IssuedAt:            int(now.Unix()),
		ExpirationTime:      int(now.Add(accessTokenTTL).Unix()),
		Subject:             strconv.FormatInt(int64(user.UserId), 10),
		Role:                strconv.FormatInt(int64(user.RoleId), 10),
		OriginatingIdentity: identity,
	}

	refreshTokenClaims := pkg.RefreshTokenClaims{
		Issuer:              issuer,
		IssuedAt:            int(now.Unix()),
		ExpirationTime:      int(now.Add(refreshTokenTTL).Unix()),
		Subject:             strconv.FormatInt(int64(user.UserId), 10),
		OriginatingIdentity: identity,
	}

	accessToken, accessTokenGenerationError := pkg.GenerateJWTAccessToken(joseHeader, accessTokenClaims, accessTokenKey)
	if accessTokenGenerationError != nil {
		log.Println("AuthService.GenerateTokens(): error occured during access token generation.")
		log.Println(accessTokenGenerationError)
		return pkg.JWTToken(""), pkg.JWTToken("")
	}

	refreshToken, refreshTokenGenerationError := pkg.GenerateJWTRefreshToken(joseHeader, refreshTokenClaims, refreshTokenKey)
	if refreshTokenGenerationError != nil {
		log.Println("AuthService.GenerateTokens(): error occured during refresh token generation.")
		log.Println(refreshTokenGenerationError)
		return pkg.JWTToken(""), pkg.JWTToken("")
	}

	return accessToken, refreshToken
}

func (a *AuthService) Create(user models.User) (models.User, error) {
	passwordHash := a.GeneratePasswordHash(user.Password)
	user.Password = passwordHash

	return a.userRepository.Create(user)
}

func (a *AuthService) UpdatePassword(userId int, password string) (models.User, error) {
	passwordHash := a.GeneratePasswordHash(password)

	return a.userRepository.UpdatePassword(models.User{UserId: userId, Password: passwordHash})
}

func (a *AuthService) UpdateUsername(userId int, username string) (models.User, error) {
	return a.userRepository.UpdateUsername(models.User{UserId: userId, UserName: username})
}

func (a *AuthService) UpdateRole(userId int, roleId int) (models.User, error) {
	return a.userRepository.UpdateUserRole(models.User{UserId: userId, RoleId: roleId})
}

func (a *AuthService) GeneratePasswordHash(password string) string {
	sha256Hasher := sha256.New()

	// hash password
	sha256Hasher.Write([]byte(password))
	// add salt to hashed password
	hashedAndSaltedPassword := sha256Hasher.Sum([]byte(salt))

	return fmt.Sprintf("%x", hashedAndSaltedPassword)
}

func (a *AuthService) ValidateAccessToken(encodedToken string, ipAddress string) *JWTTokenValidator {
	sentFrom := pkg.IPAddressIdentity{IP: ipAddress}

	validator := NewJWTTokenValidator(encodedToken, accessTokenKey, sentFrom)

	if validator.IsEverythingValid != true {
		log.Println("AuthService.ValidateTokensForRefresh(): error occured during access token validation.")
		return validator
	}

	return validator
}

func (a *AuthService) ValidateRefreshToken(encodedToken string, ipAddress string) *JWTTokenValidator {
	sentFrom := pkg.IPAddressIdentity{IP: ipAddress}

	validator := NewJWTTokenValidator(encodedToken, refreshTokenKey, sentFrom)

	if validator.IsEverythingValid != true {
		log.Println("AuthService.ValidateTokensForRefresh(): error occured during refresh token validation.")
		return validator
	}

	return validator
}

func (a *AuthService) SignHeaderAndPayload(encodedJOSEHeader string, encodedClaims string) string {
	return pkg.SignHeaderAndPayload(encodedJOSEHeader, encodedClaims, accessTokenKey)
}

const (
	AllScopeId   = 1
	OwnerScopeId = 2
)

func (a *AuthService) CheckPermissions(destination string, recordType string, recordString string, subject string, roleString string) (bool, error) {
	roleId, conversionError := strconv.Atoi(roleString)
	if conversionError != nil {
		log.Println("AuthService.CheckPermissions(): error occured during conversion from `role` string to `RoleId` integer")
		return false, conversionError
	}

	// get routeId by url
	route, err := a.routeService.GetRouteByURL(destination)
	if err != nil {
		log.Println("AuthService.CheckPermissions(): error occured during getting route by URL. Passed data: ", destination)
		return false, err
	}

	// check if no route has been found
	emptyRoute := models.Route{}
	if route == emptyRoute {
		log.Println("AuthService.CheckPermissions(): no route has been found. Passed data: ", destination)
		return false, fmt.Errorf("no route has been found. Passed data: %s", destination)
	}

	// find permissions by roleId and routeId
	permissions, err := a.permissionService.GetPermissionsByRoleIdAndRouteId(roleId, route.RouteId)
	if err != nil {
		log.Printf("AuthService.CheckPermissions(): error occured during getting permissions by RoleId and RouteId. Passed data: RoleId=%d RouteId=%d", roleId, route.RouteId)
		return false, err
	}

	// check if no permissions were found
	if len(permissions) == 0 {
		log.Printf("AuthService.CheckPermissions(): no permissions has been found by RoleId and RouteId. Passed data: RoleId=%d RouteId=%d", roleId, route.RouteId)
		return false, errors.New("no permissions has been found. Access denied")
	}

	userId, conversionError := strconv.Atoi(subject)
	if conversionError != nil {
		log.Println("AuthService.CheckPermissions(): error occured during conversion from `subject` string to `UserId` integer")
		return false, conversionError
	}

	// check if user has right to perform action over chosen record
	for _, permission := range permissions {
		// check if role has rights over of ALL the records
		if permission.ScopeId == AllScopeId {
			return true, nil // then he has rights over all records - everything is ok
		}
		if permission.ScopeId == OwnerScopeId {
			recordId, conversionError := strconv.Atoi(recordString)
			if conversionError != nil {
				log.Println("AuthService.CheckPermissions(): error occured during conversion from `recordString` string to `recordId` integer")
				return false, conversionError
			}
			isOwner, isOwnerError := a.CheckIfUserIsOwner(userId, recordType, recordId)
			log.Printf("AuthService.CheckPermissions(): if permission.ScopeId == OwnerScopeId | IsUserOwner=%t userId=%d recordType=%s recordId=%d", isOwner, userId, recordType, recordId)
			if isOwnerError != nil {
				return false, isOwnerError
			}
			if isOwner == true {
				return true, nil
			}
		}
	}

	return false, nil
}

func (a *AuthService) CheckIfUserIsOwner(userId int, recordType string, idValue int) (bool, error) {
	if recordType == "room" {
		foundRoom, err := a.roomService.GetRoomById(idValue)
		if err != nil {
			return false, err
		}
		if foundRoom.CreatedBy == userId {
			return true, nil
		}
		return false, nil
	}
	if recordType == "user" {
		foundUser, err := a.userRepository.GetUserById(idValue)
		if err != nil {
			return false, err
		}
		if foundUser.UserId == userId {
			return true, nil
		}
		return false, nil
	}
	if recordType == "booking" {
		foundBooking, err := a.bookingService.GetBookingById(idValue)
		if err != nil {
			return false, err
		}
		if foundBooking.CreatedBy == userId {
			return true, nil
		}
		return false, nil
	}
	return false, nil
}

const (
	LessThanOnePeriodError              = "JWT token contains less than one period ('.') character"
	WrongJWTTypeError                   = "JWT token should contain 3 parts: 1. JOSEHeader, 2. AccessTokenClaims, 3. Signature"
	EmptyJWTPartsError                  = "JOSEHeader, AccessTokenClaims, Signature should not be empty strings"
	IntegrityNotIntactError             = "JWT-Token was changed along the way"
	SentNotFromOriginatingIdentityError = "tokens are sent from another system&program! Possible fraudulent activity"
	TokenIsExpiredError                 = "token is expired"
)

type JWTTokenValidator struct {
	JWTTokenString string `json:"passed_jwt_params"`

	JOSEHeader         pkg.JOSEHeader
	AccessTokenClaims  pkg.AccessTokenClaims
	RefreshTokenClaims pkg.RefreshTokenClaims
	Signature          pkg.Signature
	SentFromIdentity   pkg.IPAddressIdentity
	tokenKey           string

	ValidationError   error `json:"validation_errors"`
	IsExpired         bool  `json:"is_expired"`
	IsEverythingValid bool  `json:"is_everything_valid"`
}

func NewJWTTokenValidator(jwtTokenString string, tokenKey string, sentFromIdentity pkg.IPAddressIdentity) *JWTTokenValidator {
	validator := &JWTTokenValidator{JWTTokenString: jwtTokenString, tokenKey: tokenKey, SentFromIdentity: sentFromIdentity, IsEverythingValid: false}
	validator.IsTokenValid()

	return validator
}

func (j *JWTTokenValidator) IsTokenValid() {
	j.ValidateToken()

	if j.ValidationError == nil {
		j.IsEverythingValid = true
	}
}

// ValidateToken ref: https://datatracker.ietf.org/doc/html/rfc7519#section-7.2
func (j *JWTTokenValidator) ValidateToken() {
	// 1. Verify that the JWT contains at least one period ('.') character
	tokenParts := strings.Split(j.JWTTokenString, ".")
	dotsCount := len(tokenParts) - 1
	if dotsCount < 1 {
		log.Println("JWTTokenValidator: ", LessThanOnePeriodError, "Passed data: ", j.JWTTokenString)
		j.ValidationError = errors.New(LessThanOnePeriodError)
		return
	}
	if len(tokenParts) != 3 {
		log.Println("JWTTokenValidator: ", WrongJWTTypeError, "Passed data: ", j.JWTTokenString)
		j.ValidationError = errors.New(WrongJWTTypeError)
		return
	}

	encodedJOSEHeader, encodedClaims, encodedSignature := tokenParts[0], tokenParts[1], tokenParts[2]
	// check if they are not nil strings
	if encodedJOSEHeader == "" || encodedClaims == "" || encodedSignature == "" {
		log.Println("JWTTokenValidator: ", EmptyJWTPartsError, " Passed data: ", j.JWTTokenString)
		j.ValidationError = errors.New(EmptyJWTPartsError)
		return
	}

	// verify the JOSE header and Message wasn't changed along the way
	expectedSignature := pkg.SignHeaderAndPayload(encodedJOSEHeader, encodedClaims, accessTokenKey)
	if expectedSignature != encodedSignature {
		log.Println("JWTTokenValidator: ", IntegrityNotIntactError, " Passed data: ", j.JWTTokenString)
		j.ValidationError = errors.New(IntegrityNotIntactError)
		return
	}

	// 2. Let the Encoded JOSE JOSEHeader be the portion of the JWT before the first period ('.') character
	joseHeader, joseHeaderExtractionError := pkg.ExtractJOSEHeader(encodedJOSEHeader)
	if joseHeaderExtractionError != nil {
		log.Println("JWTTokenValidator: error during extraction of JOSE JOSEHeader. Details: ", joseHeaderExtractionError)
		j.ValidationError = joseHeaderExtractionError
		return
	}
	j.JOSEHeader = joseHeader

	// 9.   Otherwise, base64url decode the Message
	if j.tokenKey == accessTokenKey {
		accessTokenClaims, claimsExtractionError := pkg.ExtractAccessTokenClaims(encodedClaims)
		if claimsExtractionError != nil {
			log.Println("JWTTokenValidator: ", claimsExtractionError, " Passed data: ", encodedClaims)
			j.ValidationError = claimsExtractionError
			return
		}
		j.AccessTokenClaims = accessTokenClaims
	} else {
		refreshTokenClaims, claimsExtractionError := pkg.ExtractRefreshTokenClaims(encodedClaims)
		if claimsExtractionError != nil {
			log.Println("JWTTokenValidator: ", claimsExtractionError, " Passed data: ", encodedClaims)
			j.ValidationError = claimsExtractionError
			return
		}
		j.RefreshTokenClaims = refreshTokenClaims
	}
	// no need for signature - we've already made sure that our token is valid with given signature
	// plus it doesn't contain anything
	j.Signature = pkg.Signature(encodedSignature)

	now := time.Now()

	// security check
	emptyAccessTokenClaims := pkg.AccessTokenClaims{}
	if j.AccessTokenClaims != emptyAccessTokenClaims {
		// check if not expired
		accessTokenExpirationTime := time.Unix(int64(j.AccessTokenClaims.ExpirationTime), 0)
		if accessTokenExpirationTime.Before(now) {
			log.Println("JWTTokenValidator: ", TokenIsExpiredError, " Passed data: ", j.AccessTokenClaims)
			j.IsExpired = true
			j.ValidationError = errors.New(TokenIsExpiredError)
			return
		}
		// check if token came from original identity (same IP)
		if j.AccessTokenClaims.OriginatingIdentity.IP != j.SentFromIdentity.IP {
			log.Println("JWTTokenValidator: ", SentNotFromOriginatingIdentityError, " Passed data: ", j.SentFromIdentity.IP)
			j.ValidationError = errors.New(SentNotFromOriginatingIdentityError)
			return
		}
	} else {
		// check if not expired
		refreshTokenExpirationTime := time.Unix(int64(j.RefreshTokenClaims.ExpirationTime), 0)
		if refreshTokenExpirationTime.Before(now) {
			log.Println("JWTTokenValidator: ", TokenIsExpiredError, " Passed data: ", j.RefreshTokenClaims)
			j.IsExpired = true
			j.ValidationError = errors.New(TokenIsExpiredError)
			return
		}
		// check if token came from original identity (same IP)
		if j.RefreshTokenClaims.OriginatingIdentity.IP != j.SentFromIdentity.IP {
			log.Println("JWTTokenValidator: ", SentNotFromOriginatingIdentityError, " Passed data: ", j.SentFromIdentity.IP)
			j.ValidationError = errors.New(SentNotFromOriginatingIdentityError)
			return
		}
	}
}
