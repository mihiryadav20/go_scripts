#!/bin/bash

# Script to run the various Go utilities in the appropriate directories

function show_usage() {
  echo "Usage: ./run.sh <command> [arguments]"
  echo ""
  echo "Available commands:"
  echo "  application <slug> <application_link>  - Update application link"
  echo "  hindi <slug> <youtube_link>            - Update Hindi YouTube link"
  echo "  regional <language> <slug> <youtube_link> - Update regional language YouTube link"
  echo ""
  echo "Examples:"
  echo "  ./run.sh application fpktnk.json https://new-link.com"
  echo "  ./run.sh hindi fpktnk.json https://www.youtube.com/watch?v=example"
  echo "  ./run.sh regional ta fpktnk.json https://www.youtube.com/watch?v=example"
  echo ""
  echo "Supported language codes for regional command:"
  echo "  te: Telugu       as: Assamese     kok: Konkani"
  echo "  gu: Gujarati     ml: Malayalam    mr: Marathi"
  echo "  mni: Manipuri    lus: Mizo        or: Odia"
  echo "  pa: Punjabi      ta: Tamil        bn: Bengali"
  echo "  ks: Kashmiri     kn: Kannada"
  exit 1
}

# Check if command is provided
if [ $# -lt 1 ]; then
  show_usage
fi

COMMAND=$1
shift

case "$COMMAND" in
  "application")
    if [ $# -lt 2 ]; then
      echo "Error: Missing arguments for application command"
      show_usage
    fi
    cd update_application_link && go run main.go "$@"
    ;;
    
  "hindi")
    if [ $# -lt 2 ]; then
      echo "Error: Missing arguments for hindi command"
      show_usage
    fi
    cd update_hindi_youtube && go run main.go "$@"
    ;;
    
  "regional")
    if [ $# -lt 3 ]; then
      echo "Error: Missing arguments for regional command"
      show_usage
    fi
    cd update_regional_youtube && go run main.go "$@"
    ;;
    
  *)
    echo "Error: Unknown command '$COMMAND'"
    show_usage
    ;;
esac
