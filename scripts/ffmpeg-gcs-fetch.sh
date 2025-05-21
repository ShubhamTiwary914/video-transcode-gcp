
filename="$1"
out="$2"

sign=$(./sign-gcs.sh -f $filename)

ffmpeg -headers "Content-Type: video/mp4" -i "$sign" $out
