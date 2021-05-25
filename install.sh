#!/usr/bin/env bash

set -e

RED='\033[0;31m'
CYAN='\033[0;36m'
ORANGE='\033[0;33m'
GREEN='\033[0;32m'
NC='\033[0m'

# Usage: make_home_symlink path_to_file name_of_file_in_home_to_symlink_to
make_home_symlink() {
    if [ -z "$1" ]; then
        >&2 echo -e "${RED}make_home_symlink: missing first argument: path to file${NC}"
        exit 1
    fi

    THIS_DOTFILE_PATH="$PWD/$1"
    if [ -z "$2" ]; then
        HOME_DOTFILE_PATH="$HOME/$1"
    else
        HOME_DOTFILE_PATH="$HOME/$2"
    fi

    printf "%s -> %s" "$1" "$HOME_DOTFILE_PATH"

    dir=`dirname "${HOME_DOTFILE_PATH}"`
    mkdir -p $dir

    if [ -L "$HOME_DOTFILE_PATH" ]; then
        echo -e " ${ORANGE}skipping, already a symlink${NC}"
        return
    fi

    if [ -f "$HOME_DOTFILE_PATH" ] && [ ! -L "$HOME_DOTFILE_PATH" ]; then
        printf " You already have a regular file at %s. Do you want to remove it? (y/n) " "$HOME_DOTFILE_PATH"
        read -r response
        if [ "$response" = "y" ] || [ "$response" = "Y" ]; then
            rm "$HOME_DOTFILE_PATH"
        else
            return
        fi
    fi

    if [[ "$(uname)" == "Darwin"* ]]; then
        # macOS permissions won't allow us to use a symlink
        ln "$THIS_DOTFILE_PATH" "$HOME_DOTFILE_PATH" 2>/dev/null
    elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
        ln -s "$THIS_DOTFILE_PATH" "$HOME_DOTFILE_PATH" 2>/dev/null
    fi
    echo -e " ${GREEN}done!${NC}"
}

print_big_notice() {
    echo -e "${CYAN}==================${NC}"
    echo -e "${CYAN}${@}${NC}"
    echo -e "${CYAN}==================${NC}"
}


make_home_symlink ".gitconfig"
make_home_symlink ".vimrc"
make_home_symlink ".vim/colors/onedark.vim"
make_home_symlink ".vim/autoload/onedark.vim"
make_home_symlink ".zshrc"
make_home_symlink ".config/alacritty/alacritty.yml"
make_home_symlink ".config/streamlink"
make_home_symlink ".config/chromium-flags.conf"
make_home_symlink ".tmux.conf"

if [[ "$(uname)" == "Darwin"* ]]; then
    print_big_notice "Detected macOS"
    make_home_symlink ".zshrc_mac"
    make_home_symlink ".hushlogin"
    make_home_symlink ".config/alacritty/macos.yml" ".config/alacritty/os.yml"
    make_home_symlink ".local/share/chatterino/Settings/commands.json" "Library/Application\ Support/chatterino/Settings/commands.json"
    make_home_symlink ".local/share/chatterino/Settings/window-layout.json" "Library/Application Support/chatterino/Settings/window-layout.json"
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    print_big_notice "Detected Linux"

    make_home_symlink ".config/alacritty/linux.yml" ".config/alacritty/os.yml"

    make_home_symlink ".config/gtk-3.0/settings.ini" ".config/gtk-3.0/settings.ini"

    make_home_symlink ".config/i3"
    make_home_symlink ".config/i3blocks"
    make_home_symlink ".config/compton.conf"
    make_home_symlink ".config/wallpaper.jpg"

    make_home_symlink ".config/Code/User/settings.json"

    make_home_symlink ".local/share/chatterino/Settings/commands.json"
    make_home_symlink ".local/share/chatterino/Settings/window-layout.json"
    
    make_home_symlink ".local/share/fonts"
    make_home_symlink ".imwheelrc"
    make_home_symlink ".sharenix.json"

    fc-cache -f
fi

