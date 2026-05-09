tries() {
	local dir
	dir="$(command tries "$@")"
	[[ -n "$dir" ]] && cd "$dir"
}
