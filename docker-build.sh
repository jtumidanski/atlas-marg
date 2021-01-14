if [[ "$1" = "NO-CACHE" ]]
then
   docker build --no-cache -f Dockerfile.dev --tag atlas-marg:latest .
else
   docker build -f Dockerfile.dev --tag atlas-marg:latest .
fi
