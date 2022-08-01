package util

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"log"
	"miniproject/ds"
	"miniproject/model"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
)

type AccessTokenCustomClaim struct {
	UserID uint64
	jwt.StandardClaims
}

type refreshTokenClaim struct {
	UserID uint64
	jwt.StandardClaims
}
type RefreshTokenCustomClaim struct {
	UID         uint64
	TokenID     uuid.UUID
	TokenString string
	Expire      time.Duration
}

func GetAccessToken(userID uint64, key *rsa.PrivateKey) (string, error) {
	//
	unixTime := time.Now()
	expTime := unixTime.Add(60 * 15) // 15min

	claims := &AccessTokenCustomClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime.Unix(),
			ExpiresAt: expTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenData, err := token.SignedString(key)
	if err != nil {
		log.Printf(err.Error())
	}
	return tokenData, err
}

func ValidateAccessToken(token string, key *rsa.PublicKey) (*AccessTokenCustomClaim, error) {
	//
	claim := AccessTokenCustomClaim{}
	tokenData, err := jwt.ParseWithClaims(token, claim, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		log.Printf(err.Error())
	}
	if !tokenData.Valid {
		log.Printf(err.Error())
	}
	data, ok := tokenData.Claims.(*AccessTokenCustomClaim)
	if !ok {
		log.Printf(err.Error())
	}
	return data, nil
}

func GetRefreshToken(userId uint64, key *rsa.PrivateKey) (*RefreshTokenCustomClaim, error) {
	//
	unixTime := time.Now()
	expTime := unixTime.Add(60 * 60 * 24 * 7) // 1week
	tokenID, err := uuid.NewRandom()
	if err != nil {
		log.Printf(err.Error())
	}

	claim := &refreshTokenClaim{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime.Unix(),
			ExpiresAt: expTime.Unix(),
			Id:        tokenID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	tokenData, err := token.SignedString(key)
	if err != nil {
		log.Printf(err.Error())
	}

	return &RefreshTokenCustomClaim{
		UID:         userId,
		TokenID:     tokenID,
		TokenString: tokenData,
		Expire:      expTime.Sub(unixTime),
	}, nil
}

func ValidRefreshToken(token string, key *rsa.PublicKey) (*RefreshTokenCustomClaim, error) {
	//
	claim := &refreshTokenClaim{}
	tokenData, err := jwt.ParseWithClaims(token, claim, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		log.Printf(err.Error())
	}
	if !tokenData.Valid {
		log.Printf(err.Error())
	}
	data, ok := tokenData.Claims.(*refreshTokenClaim)
	if !ok {
		log.Printf(err.Error())
	}
	uuidData, err := uuid.Parse(data.Id)
	return &RefreshTokenCustomClaim{
		UID:         data.UserID,
		TokenID:     uuidData,
		TokenString: token,
		Expire:      time.Duration(data.ExpiresAt),
	}, nil
}

//@@@hash password
func HashPassword(pass string) (string, error) {
	//
	salt := make([]byte, 64)
	_, err := rand.Read(salt)
	if err != nil {
		return "", nil
	}

	unHash, err := scrypt.Key([]byte(pass), salt, 32768, 8, 1, 32)
	if err != nil {
		log.Fatalf("error on crypting password %v", err.Error())
	}
	hashed := fmt.Sprintf("%v.%v", hex.EncodeToString(unHash), hex.EncodeToString(salt))
	return hashed, nil

}

//validate hash
func ValidateHashedPassword(hash, plainPassword string) (bool, error) {
	//
	splitHash := strings.Split(hash, ".")
	salt, err := hex.DecodeString(splitHash[1])
	if err != nil {
		log.Fatalf("error on : %v", err.Error())
	}
	plainHashed, err := scrypt.Key([]byte(plainPassword), salt, 32768, 8, 1, 32)
	if err != nil {
		log.Fatalf("error on %v", err.Error())
	}

	return (hex.EncodeToString(plainHashed)) == (splitHash[0]), nil
}

const (
	Yangon    = "0001"
	Mandalay  = "0002"
	Naypyitaw = "0003"
	Taunggyi  = "0004"
)

func GetBankCardNumber(city string) (string, error) {
	//
	var allCity int64
	var c, cStr, bankStr string
	info := &model.UserCityTotalBankCard{}

	tx := ds.DB.Begin()
	db := tx.Model(&model.UserCityTotalBankCard{}).Where("id = ?", 1)
	err := db.First(&info).Error
	if err == gorm.ErrRecordNotFound {
		totalCard := &model.UserCityTotalBankCard{}
		switch city {
		case "yangon":
			allCity += 1
			c = Yangon
			totalCard = &model.UserCityTotalBankCard{
				YangonCard: allCity,
			}
		case "mandalay":
			allCity += 1
			c = Mandalay
			totalCard = &model.UserCityTotalBankCard{
				MandalayCard: allCity,
			}
		case "naypyitaw":
			allCity += 1
			c = Naypyitaw
			totalCard = &model.UserCityTotalBankCard{
				NaypyitawCard: allCity,
			}
		case "taunggyi":
			allCity += 1
			c = Taunggyi
			totalCard = &model.UserCityTotalBankCard{
				TaunggyiCard: allCity,
			}
		default:
			allCity = 0
			c = "Unknown"
		}

		err = tx.Model(&model.UserCityTotalBankCard{}).Create(&totalCard).Error
		if err != nil {
			log.Printf(err.Error())
			tx.Rollback()
		}
		if len(city) > 0 {
			cStr = strconv.Itoa(int(allCity))
			length := len(cStr)
			for i := length; i < 8; i++ {
				bankStr += "0"
			}
		}

	} else {
		var aStr int64
		updateBankCard := &model.UserCityTotalBankCard{}
		if city == "yangon" {
			aStr = info.YangonCard + 1
			c = Yangon
			updateBankCard = &model.UserCityTotalBankCard{
				YangonCard: aStr,
			}
		} else if city == "mandalay" {
			aStr = info.MandalayCard + 1
			c = Mandalay
			updateBankCard = &model.UserCityTotalBankCard{
				MandalayCard: aStr,
			}
		} else if city == "naypyitaw" {
			aStr = info.NaypyitawCard + 1
			c = Naypyitaw
			updateBankCard = &model.UserCityTotalBankCard{
				NaypyitawCard: aStr,
			}
		} else if city == "taunggyi" {
			aStr = info.TaunggyiCard + 1
			c = Taunggyi
			updateBankCard = &model.UserCityTotalBankCard{
				TaunggyiCard: aStr,
			}
		}

		db := tx.Model(&model.UserCityTotalBankCard{}).Where("id = ?", 1)
		err = db.Updates(&updateBankCard).Error
		if err != nil {
			log.Printf(err.Error())
			tx.Rollback()
		}

		if len(city) > 0 {
			cStr = strconv.Itoa(int(aStr))
			length := len(cStr)
			for i := length; i < 8; i++ {
				bankStr += "0"
			}
		}

	}
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf(err.Error())
	}
	err = tx.Commit().Error
	if err != nil {
		log.Printf(err.Error())
		tx.Rollback()
	}

	return c + bankStr + cStr, nil
}

func SaveUserBankCard(city string, userId uint64, bankcard string) error {
	//

	card := &model.UserBankCard{
		Uid:           userId,
		BanCardNumber: bankcard,
	}
	err := ds.DB.Model(&model.UserBankCard{}).Create(&card).Error
	if err != nil {
		log.Printf(err.Error())
	}

	return nil
}
