#!/usr/bin/env sh

sesh () {
    local seshName=""
    while IFS= read -r line; do
        case "$line" in
            "[DIR]:  "*)
                echo "THERE"
                dir="${line#"[DIR]:  "}"  # Strip "[DIR]:  "
                seshName="sesh-$(basename "$dir")"
                cd "$dir"
                ;;
            "[EXIT]"*) return ;;
            *) echo "$line" ;;
        esac
    done < <(seshCore "$@")
    if [[ -n "$seshName" ]]; then
        tmux attach-session -t "$seshName" || echo "Creating new session $seshName ..." && tmux new-session -s "$seshName"
    fi
}
