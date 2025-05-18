#!/bin/bash


source "$(pwd)/.env"


filepath=""
signer_API="$SIGNER_API"
bucket="$BUCKET_NAME"
signed=""
filename=""


usage() {
  echo "Usage: $0 [OPTIONS]"
  echo "Options:"
  echo " -h, --help         Display this help message"
  echo " -s, --sign <url>   GCS signed URL"
  echo " -f, --file <path>  File path to upload"
}

handle_options() {
  while [ $# -gt 0 ]; do
    case "$1" in
      -h|--help)
        usage
        exit 0
        ;; 
      -f|--file)
        shift
        filepath="$1"
        ;;
      *)
        echo "Invalid option: $1"
        usage
        exit 1
        ;;
    esac
    shift
  done
}


sign_url() { 
  curl -X PUT -sLv -H "Content-Length: $(stat -c%s $filepath)" "$signer_API?bucket=$bucket&filename=$filename" 
}


#main execution ========
handle_options "$@"

#check inputs
if [ -z "$filepath" ]; then
  echo "Error: -f or --file (file-path) are required."
  usage
  exit 1
fi

if [[ ! -f "$filepath" ]]; then
  echo "File not found, exiting."
  usage
  exit 1
fi


filename=$(basename "$filepath")


sig=$(sign_url)
echo $sig


# {
#   echo "Trying to sign URL..."
#   signed=$(sign_url)
#   echo "Signed successfully: $signed"

#   echo -e "\nUploading..."
#   curl -f -X PUT -H "Content-Type: video/mp4" --upload-file "$filepath" "$signed"
#   echo "Upload successful, file: gs://$bucket/$filename"
# } || {
#   echo "Error occurred during signing or upload." >&2
#   exit 1
# }
