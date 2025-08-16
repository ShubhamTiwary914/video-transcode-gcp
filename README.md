# video-transcode-gcp


## Context 

Transcodes video (with [FFMPEG](https://ffmpeg.org/)) -> into three variants of HLS segment playlists:
- 1080p60 HLS 
- 720p60 HLS
- 480p30 HLS

with video codecs: [H.264/AVC](https://en.wikipedia.org/wiki/Advanced_Video_Coding)

<br />

References (Overview for media/video transcoding):
- [live transcoding for streaming in twitch](https://blog.twitch.tv/en/2017/10/10/live-video-transmuxing-transcoding-f-fmpeg-vs-twitch-transcoder-part-i-489c1c125f28/)
- [transcoding in ffmpeg overview](https://ffmpeg.org/ffmpeg.html#Transcoding)


Why do this? It allows for [Adaptive Bitrate streaming](https://www.cloudflare.com/learning/video/what-is-adaptive-bitrate-streaming/) over HLS (an application layer streaming protocol), that adjusts video bitrate according to network throughput.

Services like Youtube, Twitch have had similar services (now they probably have in house mechanisms)

<br />


---

## High Level Overview

<img width="1210" height="637" alt="screenshot_2025-08-17-044405" src="https://github.com/user-attachments/assets/f1dcff47-657c-4a35-a312-4163683b35db" />


> All the services here are on GCP's premises & provisioned with Terraform. 
<br/>
A rundown on what each does:

- Primary task is managed by containers with FFMPEG, handled via Cloud Run jobs -> runs transcoding batch jobs    (`Just think of it as a black box for now`) These containers just need an input stream & some output stream, and transcodes the video.
  
- So there's two GCS buckets: input & output, should be pretty obvious by those names.
  
- There's 3 Cloud Run services:
  - auth:  handled getting the signed URLs for the GCS buckets (no public access)
  - input trigger:  triggered by GCS upload event.
  - consumer:  consume from queue, a trigger batch job   (`probably a different name would be better, will just stick with this for now`)

- Cloud Task Queue helps rate limit many concurrent jobs at a time period, set to work for <=100 jobs


---

### Provisioning Steps (on GCP with terraform):
`will fill next parts later (if i don't forget)`...





