tries() {
	local dir
	dir="$(go run /Users/stefan/src/tries/2026-05-05-go/try-go)"
	[[ -n "$dir" ]] && cd "$dir"
}
