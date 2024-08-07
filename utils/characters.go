package utils

const UPPERCASE_LETTERS string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const LOWERCASE_LETTERS string = "abcdefghijklmnopqrstuvwxyz"
const DIGITS string = "0123456789"
const SPECIAL_CHARS string = "`~!@#$%^&*()-=_+,./<>?;':\"[]\\{}|"

const OTP_ALPHABET = UPPERCASE_LETTERS + LOWERCASE_LETTERS + DIGITS
const SLUG_ALPHABET = OTP_ALPHABET + "_-"
const PASSWORD_ALPHABET = OTP_ALPHABET + SPECIAL_CHARS
