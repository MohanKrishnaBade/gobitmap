#!/bin/sh
# This is a comment!

# remove old build and update with lastest one
rm -rf /Users/mohankrishnareddybade/go/src/github.com/gobitmap/build
cp -r /Applications/coding/react/blk-design-system-react-master/build /Users/mohankrishnareddybade/go/src/github.com/gobitmap
sed -i.bak 's+/blk-design-system-react+.+g' /Users/mohankrishnareddybade/go/src/github.com/gobitmap/build/index.html
sed -i.bak 's+/blk-design-system-react/static+..+g' /Users/mohankrishnareddybade/go/src/github.com/gobitmap/build/static/css/main.e8034e21.chunk.css
