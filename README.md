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

- Primary task is managed by containers with FFMPEG, handled via Cloud Run jobs -> runs transcoding batch jobs. These containers just need an input stream & some output stream, and transcodes the video.
  
- So there's two GCS buckets: input & output, should be pretty obvious by those names.
  
- There's 3 Cloud Run services:
  - auth:  handled getting the signed URLs for the GCS buckets (no public access)
  - input trigger:  triggered by GCS upload event.
  - consumer:  consume from queue, a trigger batch job   

- Cloud Task Queue helps rate limit many concurrent jobs at a time period, set to work for <=100 jobs


---

### Provisioning Steps (on GCP with terraform):

Prerequisities:
- [Taskfile](https://taskfile.dev/docs/installation)
- [Terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)
- [Gcloud CLI](https://cloud.google.com/sdk/docs/install)  -> then do [authentication via gcloud CLI](https://cloud.google.com/docs/authentication/gcloud)

<br />

##### 1. Getting the env variables for terraform:
```
TF_VAR_project_id=
TF_VAR_project_number=
TF_VAR_region=
TF_VAR_region_queue=

TF_VAR_region_bucket=
TF_VAR_hls_bucketname=

TF_VAR_job_name=

TF_VAR_artifact_reponame=
TF_VAR_artifact_packagename_signbucket=
TF_VAR_artifact_packagename_processjoTF_VAR_project_id=
TF_VAR_project_number=
TF_VAR_region=
TF_VAR_region_queue=

TF_VAR_region_bucket=
TF_VAR_hls_bucketname=

TF_VAR_job_name=

TF_VAR_artifact_reponame=
TF_VAR_artifact_packagename_signbucket=
TF_VAR_artifact_packagename_processjob=
TF_VAR_artifact_packagename_jobrunner=

TF_VAR_default_SA=
TF_VAR_user_SA=b=
TF_VAR_artifact_packagename_jobrunner=

TF_VAR_default_SA=
TF_VAR_user_SA=
```
> these are to be set in `.env` on the project root dir. (shown in .env.example)


Getting the project number and id:
```bash
gcloud projects describe $(gcloud config get-value project) --format="value(projectId,projectNumber)"
```

Get list of regions (pick any one closer to you):
```bash
gcloud compute regions list
```

Get list of regions that support task queue & bucket (since not all may be suppported from the compute regions list) - also pick one closer to you:
```bash
gcloud tasks locations list
gcloud storage location list
```

Service accounts:
```bash
#Get default Service Account:
gcloud iam service-accounts list --filter="displayName:Compute Engine default service account" --format="value(email)"


#Create User Service account:
gcloud iam service-accounts create my-sa --display-name="My Service Account"

gcloud iam service-accounts keys create my-sa-key.json \
  --iam-account=my-sa@$(gcloud config get-value project).iam.gserviceaccount.com
```

> For the remaining: names, your choice.


<br />

Then provision with Terraform:
```bash
task tf:init
task tf:plan
task tf:apply  #enter "yes" on prompt
```


2. Run the Interface:
```bash
task src:run-interface
```

<br />

The received signed URL at the end, can be used to stream out video from sources like:
- [Live Push](https://livepush.io/hlsplayer/index.html)
- [VLC Media Player](https://www.videolan.org/vlc/)








