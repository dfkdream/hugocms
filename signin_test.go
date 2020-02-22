package main

import "testing"

func TestHashPassword(t *testing.T) {
	t.Log(hashPassword("HelloWorld"))
}

func TestValidatePassword(t *testing.T) {
	t.Log(validatePassword("HelloWorld", "8fdfc4d95e3d82413534e4ef4e442b119db1331329e388e4cee857890db905c3", "5f7f0b8e9d7555499a1b84a825f4719c204759ab72a7cdf433527e72929049a8"))
}
