#!/bin/bash

#desc: 把文档从ppt 转化为 jpg
#dep: 系统需要装libreoffice 和 convert 工具

filename=$1
size=${2:-600}

if [ -e ${filename} -o ! -f ${filename} ]; then
	 echo "${filename} not exist"
	 exit 0
fi

libreoffice --headless --convert-to temp.pdf ${filename}

echo "convert to pdf ok , going to convert jpg"

convert -density ${size} temp.pdf  output/out%d.jpg

echo "convert to jpg done, "

rm -f temp.pdf
