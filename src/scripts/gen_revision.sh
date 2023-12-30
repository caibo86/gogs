#!/bin/bash

# commit id list写入到config.git-hash文件
git rev-list HEAD | sort > config.git-hash
# branch=分支名字
branch=$(git symbolic-ref HEAD 2>/dev/null | cut -d"/" -f 3)
LOCAL_REV=$(wc -l config.git-hash | awk '{print $1}')
if [ "$LOCAL_REV" -gt 1 ] ; then
    VER=$(git rev-list "$branch" | sort | join config.git-hash - | wc -l | awk '{print $1}')
        if [ "$VER" != "$LOCAL_REV" ] ; then
            VER="$VER+$("$LOCAL_REV"-"$VER")"
        fi
        if git status | grep -q "modified:" ; then
            VER="${VER}M"
        fi
        VER="$VER.$(git rev-list HEAD -n 1 | cut -c 1-16)"
        GIT_VERSION=r$VER
else
    GIT_VERSION=
    VER="x"
fi
rm -f config.git-hash
REVISION="$GIT_VERSION.$branch"
echo "$REVISION"