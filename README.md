# arb

A tool that reads [App Resource Bundle](https://github.com/google/app-resource-bundle) (.arb) files used in Flutter Internationalization.

## how it works

The tool reads 2 or more files.
The first being the reference or source file.
Each additional file is sync-ed with the source file:

- if a key is missing in the target file then add it
- if a key is no longer in the source file then remove it

All files are written back when entries sorted by key.

## install

    go install github.com/emicklei/arb@latest

## usage

    arb <source.arb> <target_nl.arb> <target_de.arb> ...


(c)2023. ernestmicklei.com