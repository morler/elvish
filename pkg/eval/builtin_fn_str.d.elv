#doc:html-id str-lt
# Outputs whether `$string`s in the given order are strictly increasing. Outputs
# `$true` when given fewer than two strings.
fn '<s' {|@string| }

#doc:html-id str-le
# Outputs whether `$string`s in the given order are strictly non-decreasing.
# Outputs `$true` when given fewer than two strings.
fn '<=s' {|@string| }

#doc:html-id str-eq
# Outputs whether `$string`s are all the same string. Outputs `$true` when given
# fewer than two strings.
fn '==s' {|@string| }

#doc:html-id str-ne
# Outputs whether `$a` and `$b` are not the same string. Equivalent to `not (==s
# $a $b)`.
fn '!=s' {|a b| }

#doc:html-id str-gt
# Outputs whether `$string`s in the given order are strictly decreasing. Outputs
# `$true` when given fewer than two strings.
fn '>s' {|@string| }

#doc:html-id str-ge
# Outputs whether `$string`s in the given order are strictly non-increasing.
# Outputs `$true` when given fewer than two strings.
fn '>=s' {|@string| }

# Output the width of `$string` when displayed on the terminal. Examples:
#
# ```elvish-transcript
# ~> wcswidth a
# ▶ (num 1)
# ~> wcswidth lorem
# ▶ (num 5)
# ~> wcswidth 你好，世界
# ▶ (num 10)
# ```
fn wcswidth {|string| }

# Override the column width of a Unicode rune to a specific non-negative value.
# If `$width` is negative, removes the override for that rune.
#
# This function is useful for handling terminal compatibility issues with certain
# Unicode characters that may not display with their expected width on specific
# terminals or environments.
#
# ```elvish-transcript
# ~> # First, see the normal width of an emoji
# ~> wcswidth 🌟
# ▶ (num 2)
# ~> # Override the width to be 1 column instead of 2
# ~> -override-wcwidth 🌟 1
# ~> wcswidth 🌟
# ▶ (num 1)
# ~> # Remove the override (use negative width)
# ~> -override-wcwidth 🌟 -1
# ~> wcswidth 🌟
# ▶ (num 2)
# ```
#
# **Note**: This function takes a single Unicode rune, not a multi-rune string.
# Use `(all)` to apply overrides to multiple characters.
#
# See also [`wcswidth`](#wcswidth).
fn '-override-wcwidth' {|rune width| }

# Convert arguments to string values.
#
# ```elvish-transcript
# ~> to-string foo [a] [&k=v]
# ▶ foo
# ▶ '[a]'
# ▶ '[&k=v]'
# ```
fn to-string {|@value| }

# Outputs a string for each `$number` written in `$base`. The `$base` must be
# between 2 and 36, inclusive. Examples:
#
# ```elvish-transcript
# ~> base 2 1 3 4 16 255
# ▶ 1
# ▶ 11
# ▶ 100
# ▶ 10000
# ▶ 11111111
# ~> base 16 1 3 4 16 255
# ▶ 1
# ▶ 3
# ▶ 4
# ▶ 10
# ▶ ff
# ```
fn base {|base @number| }
