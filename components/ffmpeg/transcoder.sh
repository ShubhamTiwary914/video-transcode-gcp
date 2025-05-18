# transcoder for a local file (video container- mkv, mp4) to HLS segments - ffmpeg
#!/bin/bash


check_audio_streams(){
    if [[ -n $(ffprobe -v error -select_streams a -show_entries stream=index -of csv=p=0 "$input") ]]; then
        HAS_AUDIO=true
    else
        HAS_AUDIO=false
    fi
}


input_validations(){
    if [[ -z "$1" ]]; then
        echo "args required: <video file>"
        exit 1
    fi

    if [[ ! -f "$1" ]]; then
        echo "File not found, exiting."
        exit 1
    fi
}



# start =========================
input="$1"
input_validations "$input"

mkdir -p out
cp "$input" ./out
cd out || exit 1

HAS_AUDIO=true
check_audio_streams


echo -e "HAS AUDIO:  $HAS_AUDIO\n\n\n"


# AUDIO ==================
if $HAS_AUDIO; then
    AUDIO_FLAGS="
        -map a:0 -c:a aac -b:a:0 192k -ac 2 \
        -map a:0 -c:a aac -b:a:1 128k -ac 2 \
        -map a:0 -c:a aac -b:a:2 96k -ac 2"
    VAR_STREAM_MAP="v:0,a:0 v:1,a:1 v:2,a:2"
else
    AUDIO_FLAGS=""
    VAR_STREAM_MAP="v:0 v:1 v:2"
fi

# VIDEO ==========
VIDEO_CHAINS="[0:v]split=3[v1][v2][v3]; \
        [v1]scale=w=1920:h=1080[v1out]; \
        [v2]scale=w=1280:h=720[v2out]; \
        [v3]scale=w=854:h=480[v3out]"


ffmpeg_transcode(){
    ffmpeg -i "$input" \
        -filter_complex "$VIDEO_CHAINS" \
            -map "[v1out]" -c:v:0 libx264 -b:v:0 5000k -maxrate:v:0 5350k -bufsize:v:0 7500k \
            -map "[v2out]" -c:v:1 libx264 -b:v:1 2800k -maxrate:v:1 2996k -bufsize:v:1 4200k \
            -map "[v3out]" -c:v:2 libx264 -b:v:2 1400k -maxrate:v:2 1498k -bufsize:v:2 2100k \
        $AUDIO_FLAGS \
        -f hls \
        -hls_time 5 \
        -hls_playlist_type vod \
        -hls_flags independent_segments \
        -hls_segment_type mpegts \
        -hls_segment_filename stream_%v/data%03d.ts \
        -master_pl_name master.m3u8 \
        -var_stream_map "$VAR_STREAM_MAP" \
        stream_%v/playlist.m3u8
}

ffmpeg_transcode