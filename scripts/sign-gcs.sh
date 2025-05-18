
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
  curl -sL -X GET "$signer_API?bucket=$bucket&filename=$filename" | jq -r '.'
}


#main execution ========
handle_options "$@"

#check inputs
if [ -z "$filepath" ]; then
  echo "Error: -f or --file (file-path) are required."
  usage
  exit 1
fi


filename=$(basename "$filepath")
sign_url
