package utils

const UPPERCASE_LETTERS string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const LOWERCASE_LETTERS string = "abcdefghijklmnopqrstuvwxyz"
const DIGITS string = "0123456789"
const SPECIAL_CHARS string = "`~!@#$%^&*()-=_+,./<>?;':\"[]\\{}|"

const SLUG_ALPHABET = UPPERCASE_LETTERS + LOWERCASE_LETTERS + DIGITS + "_-"
const PASSWORD_ALPHABET = UPPERCASE_LETTERS + LOWERCASE_LETTERS + DIGITS + SPECIAL_CHARS
