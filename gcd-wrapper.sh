DIR=`dirname "$0"`
GCD_USER=`id | sed 's/[^(]*(//;s/).*//'`
GCD_PWD_FILE="${TMPDIR-/tmp}/gcd-$GCD_USER-gcd.pwd.$$"
$DIR/gcd "$GCD_PWD_FILE" "$@"

if test -r "$GCD_PWD_FILE"; then
	GCD_PWD="`cat "$GCD_PWD_FILE"`"
	if test -n "$GCD_PWD" && test -d "$GCD_PWD"; then
		cd "$PWD/$GCD_PWD"
	fi
	unset GCD_PWD
fi

rm -f "$GCD_PWD_FILE"
unset GCD_PWD_FILE
