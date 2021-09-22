all: frames
	ffmpeg -r 1 -i frames/out_%d.png -vcodec libx264 -pix_fmt yuv420p out.mp4

frames: conte images
	../bin/images2frames conte frames images

images: conte
	cat conte | xargs -P6 -L1 -I{} ../bin/conte2image images {}

conte:
	nkf --overwrite --oc=UTF-8 scenario.txt
	../bin/scenario2conte scenario.txt > conte

clean:
	rm -rf images frames conte out.mp4

backup:
	mkdir backup
	mv -f images frames conte out.mp4 backup
	cp scenario.txt backup/
