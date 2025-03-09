#! /bin/bash

if pacman -Qt | grep neovim; then
  echo "you have neovim"
else
  sudo pacman -S neovim
fi
