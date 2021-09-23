all: frames
 # -iの順番を画像、音声の順で指定すること。順番を逆にすると壊れた動画ができる（原因不明）
	ffmpeg -r 1 -i frames/out_%d.png -i audio.wav -vcodec libx264 -pix_fmt yuv420p out.mp4

frames: conte images
	../bin/images2frames conte frames images

images: conte
	cat conte | gxargs -P6 -L1 -I{} ../bin/conte2image images {}

conte:
	nkf --overwrite --oc=UTF-8 scenario.txt
	../bin/scenario2conte scenario.txt > conte

clean:
	rm -rf images frames conte out.mp4

backup:
	mkdir backup
	mv -f images frames conte out.mp4 backup
	cp scenario.txt backup/
