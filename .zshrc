export ZSH="$HOME/.oh-my-zsh"

ZSH_THEME="robbyrussell"
DISABLE_MAGIC_FUNCTIONS=true

# Standard plugins can be found in $ZSH/plugins/
# Custom plugins may be added to $ZSH_CUSTOM/plugins/
plugins=(git)

source $ZSH/oh-my-zsh.sh

zstyle ':completion:*:*' ignored-patterns '*ORIG_HEAD'

export MOEBEL_CODE="$HOME/dev/furniture" # Change this to your code path aswell
if [ -f $MOEBEL_CODE/env/misc/zshrc ]; then
     . $MOEBEL_CODE/env/misc/zshrc
fi

export EDITOR="vim"
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$GOPATH/bin
export PATH="$PATH:$(yarn global bin)"
export NVM_DIR="$HOME/.nvm"
export PATH="$HOME/.cargo/bin:$PATH"
export PATH="$HOME/.local/bin:$PATH"
export WINIT_HIDPI_FACTOR=1
export DENO_INSTALL="/home/gempir/.deno"
export PATH="$DENO_INSTALL/bin:$PATH"

export GTK_THEME=Adwaita:dark
export QT_STYLE_OVERRIDE=adwaita

alias loadnvm=". /usr/local/opt/nvm/nvm.sh"

if [ "$(uname)" = "Darwin" ]; then
  . ~/.zshrc_mac
fi

if [ -f $HOME/.profile ]; then
     . $HOME/.profile
fi

sl () { streamlink twitch.tv/"$@" audio_only --hls-live-edge 1 --twitch-disable-hosting  }
alias ls="ls -l"
alias initsubmodules="git submodule update --init --recursive"
alias restartcompton="killall -USR1 compton"
alias ciscovpn="/opt/cisco/anyconnect/bin/vpnui"
alias tm="tmux attach || tmux"
alias tmn="tmux new"
alias ktm="killall -9 tmux"
alias dev="cd ~/dev"q
alias wowdev="cd ~/Games/World\ of\ Warcraft/_retail_/Interface/AddOns/gempUI"

# ssh autocompletion via config file
h=()
if [[ -r ~/.ssh/config ]]; then
  h=($h ${${${(@M)${(f)"$(cat ~/.ssh/config)"}:#Host *}#Host }:#*[*?]*})
fi
if [[ -r ~/.ssh/known_hosts ]]; then
  h=($h ${${${(f)"$(cat ~/.ssh/known_hosts{,2} || true)"}%%\ *}%%,*}) 2>/dev/null
fi
if [[ $#h -gt 0 ]]; then
  zstyle ':completion:*:ssh:*' hosts $h
  zstyle ':completion:*:slogin:*' hosts $h
fi

alias ydl-hq="youtube-dl -f bestvideo+bestaudio"

alias yw="yarn watch"

if ! type "pbcopy" > /dev/null; then
    alias pbpaste='xsel --clipboard --output'
    alias pbcopy='xsel --clipboard --input'
fi

alias cb="git symbolic-ref --short HEAD | pbcopy" # copy current branch name into clipboard

