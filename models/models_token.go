package models

import (
	"database/sql"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func TokenInTokenBlockList(token string) bool {
	tokenData, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return false
	}
	result, err := GetTokenBlockListFromSignature(tokenData.Signature)
	fmt.Println("token :", result.Signature, err, result.Id)
	if err != nil || result.Id == 0 {
		return false
	} else {
		return true
	}
}

func GetTokenBlockListFromSignature(signature string) (TokenBlockList, error) {
	row := DB.QueryRow(`SELECT * from tokenblocklist WHERE signature = ?`, signature)
	token := TokenBlockList{}
	var err error
	if err = row.Scan(&token.Id, &token.Signature, &token.Token, &token.UserId); err == sql.ErrNoRows {
		//fmt.Println("not found")
		return TokenBlockList{}, err
	}
	return token, nil
}

func InsertTokenBlockList(token string) (bool, error) {
	tokenData, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return false, err
	}

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO tokenblocklist ( userid, token, signature ) VALUES ( ?, ?, ? )")
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(tokenData.Claims.(jwt.MapClaims)["sub"], token, tokenData.Signature)
	if err != nil {
		return false, err
	}

	tx.Commit()
	return true, nil
}
