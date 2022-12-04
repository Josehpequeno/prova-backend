package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"unicode"
)

type Body struct {
	Password string `json:"password" binding:"required"`
	Rules    []Rule `json:"rules" binding:"required"`
}

type Rule struct {
	Content string `json:"rule" binding:"required"`
	Value   int    `json:"value" binding:"required"`
}

type Response struct {
	Verify  bool     `json:"verify"`
	NoMatch []string `json:"noMatch"`
}

func minUppercase(s string, x int) bool {
	count := 0
	for _, r := range s {
		if unicode.IsUpper(r) {
			count++
		}
	}
	if count >= x {
		return true
	} else {
		return false
	}
}
func minLowercase(s string, x int) bool {
	count := 0
	for _, r := range s {
		if unicode.IsLower(r) {
			count++
		}
	}
	if count >= x {
		return true
	} else {
		return false
	}
}

func minDigit(s string, x int) bool {
	count := 0
	for _, r := range s {
		if unicode.IsDigit(r) {
			count++
		}
	}
	if count >= x {
		return true
	} else {
		return false
	}
}

func minSpecialChars(s string, x int) bool {
	chars := "!@#$%^&*()-+/{}[]\\"
	count := 0
	for _, r := range s {
		for _, c := range chars {
			if r == c {
				count++
				break
			}
		}
	}
	if count >= x {
		return true
	} else {
		return false
	}
}
func noRepeted(s string, x int) bool {
	for i := 1; i < len(s)-1; i++ {
		if s[i] == s[i-1] || s[i] == s[i+1] {
			return false
		}
	}
	return true
}

func strongPassword(body Body) Response {
	noMatch := make([]string, 0)
	lenght := len(body.Rules)
	flag := true
	for i := 0; i < lenght; i++ {
		rule := body.Rules[i].Content
		x := body.Rules[i].Value
		switch rule {
		case "minSize":
			if len(body.Password) < x {
				flag = false
				noMatch = append(noMatch, rule)
			}
		case "minUppercase":
			flag = minUppercase(body.Password, x)
			if !flag {
				noMatch = append(noMatch, rule)
			}
		case "minLowercase":
			flag = minLowercase(body.Password, x)
			if !flag {
				noMatch = append(noMatch, rule)
			}
		case "minDigit":
			flag = minDigit(body.Password, x)
			if !flag {
				noMatch = append(noMatch, rule)
			}
		case "minSpecialChars":
			flag = minSpecialChars(body.Password, x)
			if !flag {
				noMatch = append(noMatch, rule)
			}
		case "noRepeted":
			flag = noRepeted(body.Password, x)
			if !flag {
				noMatch = append(noMatch, rule)
			}
		}
	}
	return Response{flag, noMatch}
}

func main() {
	port := "8001"
	r := gin.Default()
	r.POST("/verify", func(c *gin.Context) {
		var body Body
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response := strongPassword(body)
		c.JSON(http.StatusOK, response)
	})
	r.Run(":" + port) // listen and serve on 0.0.0.0:port (for windows "localhost:port")
}
