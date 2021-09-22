all: frames
	ffmpeg -r 1 -i frames/out%d.png -vcodec libx264 -pix_fmt yuv420p out.mp4

frames: timesheet
	../bin/images2frames timesheet frames

timesheet: conte
	mkdir images
	cat conte | xargs -I{} ../bin/conte2timesheet {} > timesheet

conte:
	cat scenario.txt | xargs -I{} ../bin/scenario2conte {} > conte
